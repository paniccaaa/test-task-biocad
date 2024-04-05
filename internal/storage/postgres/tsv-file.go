package postgres

import (
	"context"
	"fmt"
)

type TSVFile struct {
	ID           int    `json:"id"`
	FileName     string `json:"filename"`
	ErrorMessage string `json:"error_message"`
}

type ITSVFile interface {
	GetFileByName(fileName string) (int, error)
	IsFileProceed(fileName string) bool
	SaveFile(tsvFile *TSVFile) error
	UpdateFile(fileName, errMessage string) error
}

func (p *PostgresStore) GetFileByName(fileName string) (int, error) {
	// Здесь делаем запрос к базе данных, чтобы получить идентификатор файла по его имени
	// Возвращаем идентификатор файла
	var id int

	err := p.db.QueryRow("SELECT id FROM tsv_files WHERE filename = $1;", fileName).Scan(&id)
	if err != nil { //мы не нашли
		return -90, err
	}

	return id, nil
}

func (p *PostgresStore) IsFileProceed(fileName string) bool {
	// Здесь делаем запрос к базе данных, чтобы проверить, есть ли файл с таким именем и путем в таблице tsv_files
	// Если файл уже обработан, возвращаем true, иначе false

	var id int

	err := p.db.QueryRow("SELECT id FROM tsv_files WHERE filename = $1;", fileName).Scan(&id)
	if err != nil { //мы не нашли такой файл
		return false //false
	}

	return true //true
}

func (p *PostgresStore) SaveFile(tsvFile *TSVFile) error {
	const op = "postgres.SaveFile"

	query := `INSERT INTO tsv_files (filename, error_messages) VALUES ($1, $2);`

	_, err := p.db.ExecContext(context.Background(), query, tsvFile.FileName, tsvFile.ErrorMessage)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStore) UpdateFile(fileName, errMessage string) error {
	const op = "postgres.UpdateFile"

	query := `update tsv_files set error_message = $1 where filename = $2`

	_, err := p.db.ExecContext(context.Background(), query, errMessage, fileName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
