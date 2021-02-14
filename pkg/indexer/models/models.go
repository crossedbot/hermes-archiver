package models

import (
	"time"

	"github.com/crossedbot/go-warc-indexer"
	"github.com/crossedbot/simplecdxj"
)

type Records struct {
	Count   int      `json:"count"`
	Results []Record `json:"results"`
}

type Record struct {
	Id        string                `json:"id"`
	Surt      string                `json:"surt"`
	Timestamp time.Time             `json:"timestamp"`
	Type      simplecdxj.RecordType `json:"type"`
	Content   warcindexer.JsonBlock `json:"content"`
}
