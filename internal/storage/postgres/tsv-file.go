package postgres

import "time"

type TSVFile struct {
	ID           int       `json:"id"`
	FileName     string    `json:"filename"`
	Path         string    `json:"path"`
	ProcessedAt  time.Time `json:"processed_at"`
	ErrorMessage string    `json:"error_message"`
}
