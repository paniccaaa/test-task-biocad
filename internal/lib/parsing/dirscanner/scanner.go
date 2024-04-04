package dirscanner

import (
	"log/slog"
	"os"
	"time"
)

type ScanTask struct {
	FilePath string
}

type Scanner struct {
	Queue   chan ScanTask
	Log     *slog.Logger
	DirPath string
}

func NewScanner(queue chan ScanTask, log *slog.Logger, dirPath string) *Scanner {
	return &Scanner{
		Queue:   queue,
		Log:     log,
		DirPath: dirPath,
	}
}

func (s *Scanner) Start() {
	ticker := time.NewTicker(5 * time.Second) // Создаем таймер с интервалом 30 секунд
	defer ticker.Stop()                       // Обязательно останавливаем таймер перед выходом из функции

	for {
		select {
		case <-ticker.C: // Канал, который срабатывает каждые 30 секунд
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
		task := ScanTask{FilePath: filePath}
		s.Queue <- task
		//fmt.Println(task) //{./input_files/1.tsv}
	}
}
