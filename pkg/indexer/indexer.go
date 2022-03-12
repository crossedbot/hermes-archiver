package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crossedbot/common/golang/logger"
	"github.com/crossedbot/go-warc-indexer"
	"github.com/crossedbot/simplecdxj"
	"github.com/fsnotify/fsnotify"

	cdxjdb "github.com/crossedbot/hermes-archiver/pkg/database"
	"github.com/crossedbot/hermes-archiver/pkg/indexer/models"
)

// Indexer represents an interface to an WARC indexer
type Indexer interface {
	// Start starts the indexer's listener
	Start() error

	// SetEncryptionKey sets the encryption key and salt
	SetEncryptionKey(key, salt []byte)
}

// indexer implements the Indexer interface
type indexer struct {
	warcindexer.Indexer
	ctx     context.Context
	dir     string
	watcher *fsnotify.Watcher
	db      cdxjdb.CdxjRecords
}

// New returns a new WARC Indexer
func New(
	ctx context.Context,
	ipfsAddr string,
	watchDir string,
	db cdxjdb.CdxjRecords,
) (Indexer, error) {
	in, err := warcindexer.New(ctx, ipfsAddr)
	if err != nil {
		return nil, err
	}
	return &indexer{
		Indexer: in,
		ctx:     ctx,
		dir:     watchDir,
		db:      db,
	}, nil
}

// Start starts the indexer's listener for new WARC files written into the
// Indexer's directory
func (in *indexer) Start() error {
	var err error
	if in.watcher == nil {
		in.watcher, err = newWatcher(in.dir)
		if err != nil {
			return err
		}
	}
	go in.watch()
	return nil
}

// watch watches for new files written into the Indexer's directory
func (in *indexer) watch() {
	stop := false
	for !stop {
		select {
		case <-in.ctx.Done():
			stop = true
			in.watcher.Close()
			in.watcher = nil
		case event := <-in.watcher.Events:
			// if a new file is written, try to index it
			if event.Op == fsnotify.Write {
				ids, err := in.index(event.Name)
				if err != nil {
					logger.Error(fmt.Sprintf(
						"Failed to index record: %s",
						err,
					))
					continue
				} else if len(ids) > 0 {
					logger.Info(fmt.Sprintf(
						"Indexed records: %s",
						strings.Join(ids, ", "),
					))
				} else {
					logger.Info(fmt.Sprintf(
						"No records indexed for '%s'",
						filepath.Base(event.Name),
					))
				}
			}
		case err := <-in.watcher.Errors:
			logger.Error(err)
		}
	}
}

// index indexes the file at the given name, and returns their database IDs
func (in *indexer) index(name string) ([]string, error) {
	// ensure the file is a regular file
	isRegular, err := isRegularFile(name)
	if err != nil {
		return []string{}, err
	} else if !isRegular {
		return []string{},
			fmt.Errorf("'%s' is not a regular file", name)
	}
	// index the file and store the results
	cdxj, err := in.Index(name)
	if err != nil {
		return []string{}, err
	}
	return store(in.db, cdxj)
}

// newWatcher creates a new filesystem watcher
func newWatcher(start string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := filepath.Walk(start,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsDir() {
				return watcher.Add(path)
			}
			return nil
		}); err != nil {
		watcher.Close()
		return nil, err
	}
	return watcher, nil
}

// store stores a given CDXJ record into the given CDXJRecords datastore
func store(db cdxjdb.CdxjRecords, cdxj simplecdxj.CDXJ) ([]string, error) {
	ids := []string{}
	for _, rec := range cdxj.Records {
		content := warcindexer.JsonBlock{}
		if err := json.Unmarshal(rec.Content, &content); err != nil {
			return nil, err
		}
		id, err := db.Set(models.Record{
			Surt:      rec.SURT,
			Timestamp: rec.Timestamp,
			Type:      rec.Type,
			Content:   content,
		})
		if err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// isRegularFile returns true if the file at the given file path is a regular
// file, otherwise false is returned
func isRegularFile(name string) (bool, error) {
	stats, err := os.Lstat(name)
	if err != nil {
		return false, err
	}
	return stats.Mode().IsRegular(), nil
}
