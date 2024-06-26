package dirscanner

import (
	"log/slog"
	"os"
	"time"

	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

type ScanTask struct {
	FilePath string
	FileID   int
}

type Scanner struct {
	Queue     chan ScanTask
	Storage   *postgres.PostgresStore
	Log       *slog.Logger
	InputPath string
}

func NewScanner(queue chan ScanTask, storage *postgres.PostgresStore, log *slog.Logger, inputPath string) *Scanner {
	return &Scanner{
		Queue:     queue,
		Storage:   storage,
		Log:       log,
		InputPath: inputPath,
	}
}

func (s *Scanner) Start() {
	ticker := time.NewTicker(20 * time.Second) //5

	defer ticker.Stop()

	for {
		select {

		case <-ticker.C:
			s.scanDirectory()
		}
	}
}

func (s *Scanner) scanDirectory() {
	entries, err := os.ReadDir(s.InputPath)
	if err != nil {
		s.Log.Error("failed to read input directory", slog.String("err", err.Error()))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := s.InputPath + string(os.PathSeparator) + entry.Name()

		id, err := s.Storage.GetFileIDByName(filePath)
		if err != nil {
			s.Log.Error("failed to get id", slog.String("err", err.Error()))
		}

		if id == -11 {
			tsvFile := &postgres.TSVFile{
				FileName:     filePath,
				ErrorMessage: "",
			}

			id, err := s.Storage.SaveFile(tsvFile)
			if err != nil {
				s.Log.Error("failed to save tsv_file to db", slog.String("err", err.Error()))
			}

			task := ScanTask{
				FilePath: filePath,
				FileID:   id,
			}
			s.Queue <- task

		} else if id != 11 {
			s.Log.Info("file already processed, skipping", slog.String("file", filePath))
			continue
		}
	}
}
