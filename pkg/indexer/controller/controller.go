package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/crossedbot/common/golang/config"
	cdxjdb "github.com/crossedbot/hermes-archiver/pkg/database"
	"github.com/crossedbot/hermes-archiver/pkg/indexer/models"
	"github.com/crossedbot/simplecdxj"
)

var (
	// Errors
	ErrorRecordNotFound = errors.New("record not found")
)

// Controller represents an interface of an WARC-CDXJ record indexer's
// controller
type Controller interface {
	// FindRecords searches for records that match the given values and returns
	// a list of matching records
	FindRecords(
		surt string,
		types []simplecdxj.RecordType,
		before, after string,
		limit int,
	) (models.Records, error)

	// GetRecord returns the record for the given record ID
	GetRecord(id string) (models.Record, error)
}

// controller implements the Controller interface
type controller struct {
	ctx context.Context
	db  cdxjdb.CdxjRecords
}

// Config represents the configuration of an Indexer controller
type Config struct {
	DatabaseAddr string `toml:"database_addr"`
}

// control is a Contoller singleton
var control Controller
var controllerOnce sync.Once
var V1 = func() Controller {
	controllerOnce.Do(func() {
		var cfg Config
		if err := config.Load(&cfg); err != nil {
			panic(err)
		}
		ctx := context.Background()
		db := cdxjdb.New(ctx, cfg.DatabaseAddr)
		control = New(ctx, db)
	})
	return control
}

// New returns a new Controller
func New(ctx context.Context, db cdxjdb.CdxjRecords) Controller {
	return &controller{ctx: ctx, db: db}
}

// FindRecords searches for records that match the given values and returns a
// list of matching records
func (c *controller) FindRecords(
	surt string,
	types []simplecdxj.RecordType,
	before, after string,
	limit int,
) (models.Records, error) {
	var err error
	s := []string{}
	for _, t := range types {
		s = append(s, t.String())
	}
	b := int64(0)
	if before != "" {
		b, err = strconv.ParseInt(before, 10, 64)
		if err != nil {
			return models.Records{},
				fmt.Errorf("'before' (%s) is not an integer", before)
		}
	}
	a := int64(0)
	if after != "" {
		a, err = strconv.ParseInt(after, 10, 64)
		if err != nil {
			return models.Records{},
				fmt.Errorf("'after' (%s) is not an integer", after)
		}
	}
	return c.db.Find(surt, s, b, a, limit)
}

// GetRecord returns the record for the given record ID
func (c *controller) GetRecord(id string) (models.Record, error) {
	r, err := c.db.Get(id)
	if err == cdxjdb.ErrNotFound {
		return models.Record{}, ErrorRecordNotFound
	}
	return r, err
}
