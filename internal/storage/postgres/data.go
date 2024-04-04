package postgres

import (
	"context"
	"fmt"
)

type DataItem struct {
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
}

func (p *PostgresStore) SaveDataItem(item *DataItem) error {
	const op = "postgres.SaveDataItem"

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
