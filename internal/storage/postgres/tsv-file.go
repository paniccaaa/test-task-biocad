package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type TSVFile struct {
	ID           int    `json:"id"`
	FileName     string `json:"filename"`
	ErrorMessage string `json:"error_message"`
}

type ITSVFile interface {
	GetFileIDByName(fileName string) (int, error)
	SaveFile(tsvFile *TSVFile) error
	//GetTSVFileByID(tsvFileID int) (*TSVFile, error)
	UpdateFile(fileName, errMessage string) error
}

func (p *PostgresStore) GetFileIDByName(fileName string) (int, error) {
	const op = "storage.postgres.GetFileIDByName"
	var id int

	err := p.db.QueryRow("SELECT id FROM tsv_files WHERE filename = $1;", fileName).Scan(&id)
	// if not found
	if err != nil {
		if err == sql.ErrNoRows {
			return -11, nil
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *PostgresStore) GetTSVFileByID(tsvFileID int) (*TSVFile, error) {
	const op = "storage.postgres.GetTSVFileByID"
	tf := &TSVFile{}

	row := p.db.QueryRow("SELECT * FROM tsv_files WHERE id = $1;", tsvFileID)
	if err := row.Scan(&tf.ID, &tf.FileName, &tf.ErrorMessage); err != nil {
		fmt.Println("Ошибка сканирования строки:", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tf, nil
}

func (p *PostgresStore) SaveFile(tsvFile *TSVFile) (int, error) {
	const op = "storage.postgres.SaveFile"

	query := `INSERT INTO tsv_files (filename) VALUES ($1) RETURNING id;`

	var id int
	err := p.db.QueryRow(query, tsvFile.FileName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *PostgresStore) UpdateFile(fileName, errMessage string) error {
	const op = "storage.postgres.UpdateFile"

	query := `update tsv_files set error_message = $1 where filename = $2`

	_, err := p.db.ExecContext(context.Background(), query, errMessage, fileName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
