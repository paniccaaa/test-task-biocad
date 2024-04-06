package generator

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

type Generator struct {
	Log        *slog.Logger
	OutputPath string
	Storage    *postgres.PostgresStore
}

func NewGenerator(log *slog.Logger, outputPath string, storage *postgres.PostgresStore) *Generator {
	return &Generator{
		Log:        log,
		OutputPath: outputPath,
		Storage:    storage,
	}
}

func (g *Generator) Start() {
	ticker := time.NewTicker(10 * time.Second) // Создаем таймер с интервалом 30 секунд

	defer ticker.Stop() // Обязательно останавливаем таймер перед выходом из функции

	for {
		select {

		case <-ticker.C:
			g.ScanDirectory()
		}
	}
}

func (g *Generator) CheckDB() {
	unitsGUID, err := g.Storage.GetUniqueUnitGUID()
	if err != nil {
		fmt.Errorf("failed to get unique unit guid %w", err)
	}

	fmt.Println(unitsGUID)

	// for i, v := range unitsGUID {

	// }
}

func (g *Generator) ScanDirectory() {

	g.CheckDB()

	entries, err := os.ReadDir(g.OutputPath)
	if err != nil {
		g.Log.Error("failed to read output directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fmt.Println("hello from generator", entry.Name())

		// filePath := g.OutputPath + string(os.PathSeparator) + entry.Name()

		// id, err := g.Storage.GetFileIDByName(filePath)
		// if err != nil {
		// 	g.Log.Error("failed to get id: %w", err)
		// }

		// if id == -11 {
		// 	tsvFile := &postgres.TSVFile{
		// 		FileName:     filePath,
		// 		ErrorMessage: "",
		// 	}

		// 	id, err := s.Storage.SaveFile(tsvFile)
		// 	if err != nil {
		// 		s.Log.Error("failed to save tsv_file to db", "err", err)
		// 	}

		// 	task := ScanTask{
		// 		FilePath: filePath,
		// 		FileID:   id,
		// 	}
		// 	s.Queue <- task

		// } else if id != 11 {
		// 	s.Log.Info("file already processed, skipping", slog.String("file", filePath))
		// 	continue
		// }
	}
}

func CreatePDF() {

}
