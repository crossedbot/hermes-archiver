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
	ErrorRecordNotFound = errors.New("record not found")
)

type Controller interface {
	FindRecords(
		surt string,
		types []simplecdxj.RecordType,
		before, after string,
		limit int,
	) (models.Records, error)
	GetRecord(id string) (models.Record, error)
}

type controller struct {
	ctx context.Context
	db  cdxjdb.CdxjRecords
}

type Config struct {
	DatabaseAddr string `toml:"database_addr"`
}

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

func New(ctx context.Context, db cdxjdb.CdxjRecords) Controller {
	return &controller{ctx: ctx, db: db}
}

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

func (c *controller) GetRecord(id string) (models.Record, error) {
	r, err := c.db.Get(id)
	if err == cdxjdb.ErrNotFound {
		return models.Record{}, ErrorRecordNotFound
	}
	return r, err
}
