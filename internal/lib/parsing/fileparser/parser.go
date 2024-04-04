package fileparser

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"

	"github.com/paniccaaa/test-task-biocad/internal/lib/parsing/dirscanner"
	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
)

type Parser struct {
	Storage *postgres.PostgresStore
	Queue   chan dirscanner.ScanTask
	Log     *slog.Logger
}

func NewParser(storage *postgres.PostgresStore, queue chan dirscanner.ScanTask, log *slog.Logger) *Parser {
	return &Parser{
		Storage: storage,
		Queue:   queue,
		Log:     log,
	}
}

func (p *Parser) ProcessNext() {
	for task := range p.Queue {
		if err := p.parseFile(task.FilePath); err != nil {
			p.Log.Error("failed to parse file", slog.String("file", task.FilePath), slog.String("err", fmt.Sprint(err)))
		}
	}
}

func (p *Parser) parseFile(filePath string) error {
	const op = "lib.parsing.fileparser"

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("%s: failed to open file: %w", op, err)
	}

	defer file.Close()

	r := csv.NewReader(file)
	r.Comma = '\t'

	_, err = r.Read()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for {
		row, err := r.Read()
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		item := &postgres.DataItem{
			N:         row[0],
			MQTT:      row[1],
			Invid:     row[2],
			UnitGUID:  row[3],
			MsgID:     row[4],
			Text:      row[5],
			Context:   row[6],
			Class:     row[7],
			Level:     row[8],
			Area:      row[9],
			Addr:      row[10],
			Block:     row[11],
			Type:      row[12],
			Bit:       row[13],
			InvertBit: row[14],
			TSVFileID: 1, // Заполните это поле в соответствии с вашей логикой
		}

		if err := p.Storage.SaveDataItem(item); err != nil {
			return fmt.Errorf("%s: failed to save data item to db: %w", op, err)
		}

	}

}
