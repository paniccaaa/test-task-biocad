package generator

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/paniccaaa/test-task-biocad/internal/storage/postgres"
	"github.com/signintech/gopdf"
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
	ticker := time.NewTicker(10 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.scanDirectory()
		}
	}
}

func (g *Generator) scanDirectory() {

	unitsGUID, err := g.Storage.GetUniqueUnitGUID()
	if err != nil {
		g.Log.Error("failed to get slice of unit guid", slog.String("err", err.Error()))
		return
	}

	for _, unitGUID := range unitsGUID {

		pdfFilePath := filepath.Join(g.OutputPath, unitGUID+".pdf")

		if _, err := os.Stat(pdfFilePath); err == nil {
			g.Log.Info("PDF file already exists, skipping", slog.String("file", pdfFilePath))
			continue
		}

		dataItems, err := g.Storage.GetDataByUnitGUID(unitGUID)
		if err != nil {
			g.Log.Error("failed to get data by unit_guid", slog.String("err", err.Error()))
			continue
		}

		// 1 request to db (only for optimization)
		tsvFile, err := g.Storage.GetTSVFileByID(dataItems[0].TSVFileID)
		if err != nil {
			g.Log.Error("failed to get tsvFile by id", slog.String("err", err.Error()))
			continue
		}

		err = g.createPDF(unitGUID, dataItems, tsvFile)
		if err != nil {
			g.Log.Error("failed to create pdf", slog.String("err", err.Error()))
		}

	}
}

func (g *Generator) createPDF(unitGUID string, dataItems []*postgres.DataItem, tf *postgres.TSVFile) error {
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	pdf.AddPage()

	err := pdf.AddTTFFont("Roboto", "./ttf/Roboto-Black.ttf")
	if err != nil {
		return err
	}

	err = pdf.SetFont("Roboto", "", 14)
	if err != nil {
		return err
	}

	for _, d := range dataItems {

		properties := []string{
			"ID: " + strconv.Itoa(d.ID),
			"N: " + d.N,
			"MQTT: " + d.MQTT,
			"Invid: " + d.Invid,
			"Unit_GUID: " + d.UnitGUID,
			"MsgID: " + d.MsgID,
			"Text: " + d.Text,
			"Context: " + d.Context,
			"Class: " + d.Class,
			"Level: " + d.Level,
			"Area: " + d.Area,
			"Addr: " + d.Addr,
			"Block: " + d.Block,
			"Type: " + d.Type,
			"Bit: " + d.Bit,
			"InvertBit: " + d.InvertBit,
		}

		for _, prop := range properties {
			pdf.Cell(nil, prop)
			pdf.Br(20)
		}
		pdf.Br(20)

		pdf.AddPage()
	}

	pdf.Br(30)
	pdf.Cell(nil, "TSV File ID: "+strconv.Itoa(tf.ID))
	pdf.Br(20)
	pdf.Cell(nil, "TSV File Name: "+tf.FileName)
	pdf.Br(20)
	pdf.Cell(nil, "TSV File Error Message: "+tf.ErrorMessage)

	pdfFilePath := filepath.Join(g.OutputPath, unitGUID+".pdf")
	err = pdf.WritePdf(pdfFilePath)
	if err != nil {
		return err
	}

	return nil
}
