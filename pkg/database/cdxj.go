package database

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/google/uuid"

	"github.com/crossedbot/hermes-archiver/pkg/indexer/models"
	"github.com/crossedbot/simplecdxj"
)

const (
	CdxjKeyPrefix = "cdxj:"
)

var (
	ErrNotFound = errors.New("cdxj record not found")
)

type CdxjRecords interface {
	Init() error
	Set(rec models.Record) (string, error)
	Get(recordId string) (models.Record, error)
	Find(surt string, types []string, before, after int64, limit int) (models.Records, error)
}

type cdxjRecords struct {
	*redisearch.Client
	ctx context.Context
}

func New(ctx context.Context, addr string) CdxjRecords {
	return &cdxjRecords{
		Client: redisearch.NewClient(addr, "cdxjs"),
		ctx:    ctx,
	}
}

func (d *cdxjRecords) Init() error {
	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewSortableTextField("surt", 100)).
		AddField(redisearch.NewSortableNumericField("timestamp")).
		AddField(redisearch.NewTextField("type")).
		AddField(redisearch.NewTextField("content"))
	d.Drop()
	return d.CreateIndex(schema)
}

func (d *cdxjRecords) Set(rec models.Record) (string, error) {
	id := uuid.New().String()
	idx := fmt.Sprintf("%s%s", CdxjKeyPrefix, id)
	content, err := json.Marshal(rec.Content)
	if err != nil {
		return "", err
	}
	doc := redisearch.NewDocument(idx, 1.0)
	doc.Set("surt", rec.Surt).
		Set("timestamp", rec.Timestamp.Unix()).
		Set("type", rec.Type.String()).
		Set("content", base64.URLEncoding.EncodeToString(content))
	return id, d.Index(doc)
}

func (d *cdxjRecords) Get(recordId string) (models.Record, error) {
	idx := fmt.Sprintf("%s%s", CdxjKeyPrefix, recordId)
	doc, err := d.Client.Get(idx)
	if err != nil {
		return models.Record{}, err
	}
	if doc == nil {
		return models.Record{}, ErrNotFound
	}
	return parseCdxjRecordDoc(*doc)
}

func (d *cdxjRecords) Find(surt string, types []string, before, after int64, limit int) (models.Records, error) {
	raw := []string{}
	if surt != "" {
		raw = append(raw, fmt.Sprintf("@surt:%%%s%%", surt))
	}
	if len(types) > 0 {
		raw = append(raw, fmt.Sprintf("@type:(%s)", strings.Join(types, "|")))
	}
	if after > 0 && before > 0 {
		// filter within date range
		raw = append(raw, fmt.Sprintf("@timestamp:[%d %d]", after, before))
	} else if after > 0 {
		// filter by lower range
		raw = append(raw, fmt.Sprintf("@timestamp:[(%d +inf]", after))
	} else if before > 0 {
		// filter by upper range
		raw = append(raw, fmt.Sprintf("@timestamp:[-inf (%d]", before))
	}
	qs := "*"
	if len(raw) > 0 {
		qs = strings.Join(raw, " ")
	}
	q := redisearch.NewQuery(qs)
	if limit > 0 {
		q.Limit(0, limit)
	}
	docs, _, err := d.Client.Search(q)
	if err != nil {
		return models.Records{}, err
	}
	recs := []models.Record{}
	for _, doc := range docs {
		rec, err := parseCdxjRecordDoc(doc)
		if err != nil {
			return models.Records{}, err
		}
		recs = append(recs, rec)
	}
	return models.Records{
		Count:   len(recs),
		Results: recs,
	}, nil
}

func parseCdxjRecordDoc(doc redisearch.Document) (models.Record, error) {
	rec := models.Record{Id: strings.TrimPrefix(doc.Id, CdxjKeyPrefix)}
	if s, ok := doc.Properties["surt"]; ok {
		rec.Surt = s.(string)
	}
	if t, ok := doc.Properties["timestamp"]; ok {
		ts := t.(string)
		i64, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return models.Record{}, err
		}
		rec.Timestamp = time.Unix(i64, 0)
	}
	if t, ok := doc.Properties["type"]; ok {
		ty, err := simplecdxj.ParseRecordType(t.(string))
		if err != nil {
			return models.Record{}, err
		}
		rec.Type = ty
	}
	if c, ok := doc.Properties["content"]; ok {
		b, err := base64.URLEncoding.DecodeString(c.(string))
		if err != nil {
			return models.Record{}, err
		}
		if err := json.Unmarshal(b, &rec.Content); err != nil {
			return models.Record{}, err
		}
	}
	return rec, nil
}
