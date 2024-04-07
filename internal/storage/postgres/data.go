package postgres

import (
	"context"
	"fmt"
)

type DataItem struct {
	ID        int    `json:"id"`
	N         string `json:"n"`
	MQTT      string `json:"mqtt"`
	Invid     string `json:"invid"`
	UnitGUID  string `json:"unit_guid"`
	MsgID     string `json:"msg_id"`
	Text      string `json:"text"`
	Context   string `json:"context"`
	Class     string `json:"class"`
	Level     string `json:"level"`
	Area      string `json:"area"`
	Addr      string `json:"addr"`
	Block     string `json:"block"`
	Type      string `json:"type"`
	Bit       string `json:"bit"`
	InvertBit string `json:"invert_bit"`
	TSVFileID int    `json:"tsv_file_id"`
}

type IData interface {
	SaveDataItem(item *DataItem) error
	GetUniqueUnitGUID() ([]string, error)
	GetDataByUnitGUID(unitGUID int) ([]*DataItem, error)
}

func (p *PostgresStore) SaveDataItem(item *DataItem) error {
	const op = "storage.postgres.SaveDataItem"

	query := `INSERT INTO data (n, mqtt, invid, unit_guid, msg_id, text,
															context, class, level, area, addr, block,
															type, bit, invert_bit, tsv_file_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`

	_, err := p.db.ExecContext(context.Background(), query, item.N, item.MQTT, item.Invid, item.UnitGUID, item.MsgID,
		item.Text, item.Context, item.Class, item.Level, item.Area, item.Addr, item.Block, item.Type, item.Bit,
		item.InvertBit, item.TSVFileID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStore) GetUniqueUnitGUID() ([]string, error) {
	const op = "storage.postgres.SaveDataItem"

	query := `SELECT DISTINCT unit_guid FROM data`
	sliceOfUnitGUID := []string{}

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	for rows.Next() {
		var unitGUID string

		if err := rows.Scan(&unitGUID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		sliceOfUnitGUID = append(sliceOfUnitGUID, unitGUID)
	}

	return sliceOfUnitGUID, nil
}

func (p *PostgresStore) GetDataByUnitGUID(unitGUID string) ([]*DataItem, error) {
	const op = "storage.postgres.GetFileByName"
	data := []*DataItem{}

	rows, err := p.db.Query("SELECT * FROM data WHERE unit_guid = $1", unitGUID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	for rows.Next() {
		var d DataItem
		if err := rows.Scan(&d.ID, &d.N, &d.MQTT, &d.Invid, &d.UnitGUID, &d.MsgID, &d.Text, &d.Context, &d.Class,
			&d.Level, &d.Area, &d.Addr, &d.Block, &d.Type, &d.Bit, &d.InvertBit, &d.TSVFileID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		data = append(data, &d)
	}
	return data, nil
}
