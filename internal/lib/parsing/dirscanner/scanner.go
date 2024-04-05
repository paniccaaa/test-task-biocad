package dirscanner

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

type ScanTask struct {
	FilePath string
}

type Scanner struct {
	Queue   chan ScanTask
	Storage *postgres.PostgresStore
	Log     *slog.Logger
	DirPath string
	Ctx     context.Context // добавлено поле контекста
}

func NewScanner(queue chan ScanTask, storage *postgres.PostgresStore, log *slog.Logger, dirPath string, ctx context.Context) *Scanner {
	return &Scanner{
		Queue:   queue,
		Log:     log,
		DirPath: dirPath,
		Ctx:     ctx,
	}
}

func (s *Scanner) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // Создаем таймер с интервалом 30 секунд

	defer ticker.Stop() // Обязательно останавливаем таймер перед выходом из функции

	for {
		select {
		// case <-ctx.Done(): // Канал, который срабатывает каждые 5 секунд
		// 	return
		case <-ticker.C:
			s.scanDirectory()
		}
	}
}

func (s *Scanner) scanDirectory() {
	entries, err := os.ReadDir(s.DirPath)
	if err != nil {
		s.Log.Error("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := s.DirPath + string(os.PathSeparator) + entry.Name()

		//without check from db
		task := ScanTask{FilePath: filePath}
		s.Queue <- task

		//with check from db
		// id, err := s.Storage.GetFileByName(filePath)
		// if err != nil {
		// 	s.Log.Error("failed to get id: %w", err)
		// }
		// if id == -90 {
		// 	tsvFile := &postgres.TSVFile{
		// 		FileName: filePath,
		// 	}

		// 	err := s.Storage.SaveFile(tsvFile)
		// 	if err != nil {
		// 		s.Log.Error("failed to save tsv_file to db", "err", err)
		// 	}

		// 	task := ScanTask{FilePath: filePath}
		// 	s.Queue <- task
		// } else {
		// 	s.Log.Info("file already processed, skipping", slog.String("file", filePath))
		// 	continue
		// }
	}
}
