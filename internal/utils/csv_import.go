package utils

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"backend/internal/models"
)

func ParsePestisidaCSV(filePath string) ([]models.Pestisida, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var pestisidas []models.Pestisida

	for i, row := range rows {

		// skip header
		if i == 0 {
			continue
		}

		nilai := strings.Replace(strings.TrimSpace(row[6]), ",", ".", 1)

		dosis, err := strconv.ParseFloat(nilai, 64)
		if err != nil {
			dosis = 0
		}

		data := models.Pestisida{
			NamaPestisida:    row[0],
			BahanAktif:       row[1],
			BentukFormulasi:  row[2],
			JenisPestisida:   row[3],
			SasaranKomoditas: row[4],
			HamaSasaran:      row[5],
			DosisNilai:       dosis,
			DosisSatuan:      row[7],
		}
		pestisidas = append(pestisidas, data)
	}

	return pestisidas, nil
}
