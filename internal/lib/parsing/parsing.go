package parsing

import (
	"fmt"
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
	// Инициализация очереди для обработки файлов
	fileQueue := make(chan dirscanner.ScanTask)

	// Инициализация фоновой задачи для сканирования директории
	dirscan := dirscanner.NewScanner(fileQueue, storage, log, cfg.InputPath)
	go func() {
		dirscan.Start()
	}()

	// Создание механизма обработки файлов
	fileProcessor := fileparser.NewParser(storage, fileQueue, log)

	go func() {
		fileProcessor.ProcessNext()
	}()

	scanTask := <-fileQueue

	_ = scanTask
	fmt.Println("parsing", scanTask.FilePath, scanTask.FileID)

	// tsvFile, err := storage.GetTSVFileByID(scanTask.FileID)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(tsvFile)
	return nil
}
