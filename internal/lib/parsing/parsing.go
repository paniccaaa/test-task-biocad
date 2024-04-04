package parsing

import (
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
	go func() {
		dirscan := dirscanner.NewScanner(fileQueue, log, cfg.DirPath)
		dirscan.Start()
	}()

	// Создание механизма обработки файлов
	fileProcessor := fileparser.NewParser(storage, fileQueue, log)

	go fileProcessor.ProcessNext()

	// // Создание механизма создания выходных файлов
	// outputGenerator := outputgenerator.NewGenerator(config.OutputDir, storage, log)

	// // Процесс обработки файлов
	// go func() {
	// 	for {
	// 		fileProcessor.ProcessNext()
	// 		time.Sleep(time.Second * 5) // Пауза между обработкой файлов
	// 	}
	// }()

	// // Процесс создания выходных файлов
	// go func() {
	// 	for {
	// 		outputGenerator.GenerateNext()
	// 		time.Sleep(time.Second * 5) // Пауза между созданием выходных файлов
	// 	}
	// }()

	return nil
}