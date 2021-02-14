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
	cdxjdb "github.com/crossedbot/hermes-archiver/pkg/database"
	"github.com/crossedbot/hermes-archiver/pkg/indexer/models"
	"github.com/crossedbot/simplecdxj"
	"github.com/fsnotify/fsnotify"
)

type Indexer interface {
	Start() error
	SetEncryptionKey(key, salt []byte)
}

type indexer struct {
	warcindexer.Indexer
	ctx     context.Context
	dir     string
	watcher *fsnotify.Watcher
	db      cdxjdb.CdxjRecords
}

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

func (in *indexer) watch() {
	stop := false
	for !stop {
		select {
		case <-in.ctx.Done():
			stop = true
			in.watcher.Close()
			in.watcher = nil
		case event := <-in.watcher.Events:
			if event.Op == fsnotify.Write {
				stats, err := os.Lstat(event.Name)
				if err == nil && stats.Mode().IsRegular() {
					cdxj, err := in.Index(event.Name)
					if err == nil {
						ids, err := store(in.db, cdxj)
						if err != nil {
							logger.Error(fmt.Sprintf(
								"Failed to index record: %s",
								err,
							))
						}
						logger.Info(fmt.Sprintf(
							"Indexed records: %s",
							strings.Join(ids, ", "),
						))
					}
				}
			}
		case err := <-in.watcher.Errors:
			logger.Error(err)
		}
	}
}

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
