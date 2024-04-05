package parsing

import (
	"context"
	"log/slog"

	"github.com/paniccaaa/test-task-biocad/internal/config"
	"github.com/paniccaaa/test-task-biocad/internal/lib/parsing/dirscanner"
	"github.com/paniccaaa/test-task-biocad/internal/lib/parsing/fileparser"
	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

type ScanTask struct {
	FilePath string
}

func Start(cfg *config.Config, log *slog.Logger, storage *postgres.PostgresStore) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация очереди для обработки файлов
	fileQueue := make(chan dirscanner.ScanTask)

	// Инициализация фоновой задачи для сканирования директории
	dirscan := dirscanner.NewScanner(fileQueue, storage, log, cfg.DirPath, ctx)
	go func() {
		dirscan.Start(ctx)
	}()

	// Создание механизма обработки файлов
	fileProcessor := fileparser.NewParser(storage, fileQueue, log)

	go fileProcessor.ProcessNext()

	return nil
}
