package models

import ()

type Replay struct {
	Uri       string `json:"uri"`
	Title     string `json:"title"`
	Sha       string `json:"sha"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"content"`
}
