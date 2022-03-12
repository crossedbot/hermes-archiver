package models

import (
	"time"

	"github.com/crossedbot/go-warc-indexer"
	"github.com/crossedbot/simplecdxj"
)

// Records represents a list of CDXJ records
type Records struct {
	Count   int      `json:"count"`
	Results []Record `json:"results"`
}

// Record represents a single CDXJ record
type Record struct {
	Id        string                `json:"id"`
	Surt      string                `json:"surt"`
	Timestamp time.Time             `json:"timestamp"`
	Type      simplecdxj.RecordType `json:"type"`
	Content   warcindexer.JsonBlock `json:"content"`
}
