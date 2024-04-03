package postgres

type DataItem struct {
	ID        int    `json:"id"`
	N         int    `json:"n"`
	MQTT      string `json:"mqtt"`
	Invid     string `json:"invid"`
	UnitGUID  string `json:"unit_guid"`
	MsgID     string `json:"msg_id"`
	Text      string `json:"text"`
	Context   string `json:"context"`
	Class     string `json:"class"`
	Level     int    `json:"level"`
	Area      string `json:"area"`
	Addr      string `json:"addr"`
	Block     string `json:"block"`
	Type      string `json:"type"`
	Bit       int    `json:"bit"`
	InvertBit string `json:"invert_bit"`
	TSVFileID int    `json:"tsv_file_id"`
}
