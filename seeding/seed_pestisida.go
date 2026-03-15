package seed

import (
	"log"

	"backend/config"
	"backend/internal/models"
	"backend/internal/utils"
)

func SeedPestisida() error {

	filePath := "dataset/dataset pestisida - Full Padi.csv"

	data, err := utils.ParsePestisidaCSV(filePath)
	if err != nil {
		return err
	}

	csvCount := len(data)

	if csvCount == 0 {
		log.Println("No data found in CSV")
		return nil
	}

	var dbCount int64

	err = config.DB.Model(&models.Pestisida{}).Count(&dbCount).Error
	if err != nil {
		return err
	}

	// skip seeding jika jumlah sama
	if int(dbCount) == csvCount {
		log.Println("Seeder skipped: dataset already imported")
		return nil
	}

	log.Println("Importing pestisida dataset...")

	// bulk insert
	err = config.DB.CreateInBatches(data, 200).Error
	if err != nil {
		return err
	}

	log.Println("Pestisida dataset imported:", csvCount)

	return nil
}