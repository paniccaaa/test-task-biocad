package postgres

import "time"

type TSVFile struct {
	ID           int       `json:"id"`
	FileName     string    `json:"filename"`
	Path         string    `json:"path"`
	ProcessedAt  time.Time `json:"processed_at"`
	ErrorMessage string    `json:"error_message"`
}

type ITSVFile interface {
	GetFileByName(fileName string) int //return id
	IsFileProceed(fileName string) bool
}

func (p *PostgresStore) GetFileByName(fileName string) (int, error) {
	// Здесь делаем запрос к базе данных, чтобы получить идентификатор файла по его имени
	// Возвращаем идентификатор файла
	var id int

	err := p.db.QueryRow("SELECT id FROM tsv_files WHERE filename = $1", fileName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PostgresStore) IsFileProceed(fileName string) bool {
	// Здесь делаем запрос к базе данных, чтобы проверить, есть ли файл с таким именем и путем в таблице tsv_files
	// Если файл уже обработан, возвращаем true, иначе false

	var id int

	err := p.db.QueryRow("SELECT id FROM tsv_files WHERE filename = $1", fileName).Scan(&id)
	if err != nil {
		return false
	}

	return true
}
