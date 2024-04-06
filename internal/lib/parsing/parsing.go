package parsing

import (
	"log/slog"

	"github.com/paniccaaa/test-task-biocad/internal/config"
	"github.com/paniccaaa/test-task-biocad/internal/lib/parsing/dirscanner"
	"github.com/paniccaaa/test-task-biocad/internal/lib/parsing/fileparser"
	"github.com/paniccaaa/test-task-biocad/internal/lib/parsing/generator"
	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

type ScanTask struct {
	FilePath string
}

func Start(cfg *config.Config, log *slog.Logger, storage *postgres.PostgresStore) error {
	// chan for queue
	fileQueue := make(chan dirscanner.ScanTask)

	// goroutine for start scan dir
	go func() {
		dirscan := dirscanner.NewScanner(fileQueue, storage, log, cfg.InputPath)
		dirscan.Start()
	}()

	// goroutine for parsing files
	go func() {
		fileProcessor := fileparser.NewParser(storage, fileQueue, log)
		fileProcessor.ProcessNext()
	}()

	// goroutine
	go func() {
		generator := generator.NewGenerator(log, cfg.OutputPath, storage)
		generator.Start()
	}()

	// scanTask := <-fileQueue

	// _ = scanTask
	// fmt.Println("parsing", scanTask.FilePath, scanTask.FileID)

	return nil
}
