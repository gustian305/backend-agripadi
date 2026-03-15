package service

import (
	"backend/internal/dto"
	"backend/internal/models"
	"gorm.io/gorm"
)

type ExpertSystemInterface interface {
	GetPestisidaByHama(hama string) ([]dto.PestisidaResponse, error)
	CalculateDosage(pestisida []dto.PestisidaResponse, luasLahan float64) ([]dto.PestisidaRecommendation, error)
	GetAllPestisida() ([]models.Pestisida, error)
}

type ExpertSystemService struct {
	DB *gorm.DB
}

func NewExpertSystemService(db *gorm.DB) *ExpertSystemService {
	return &ExpertSystemService{
		DB: db,
	}
}

func (s *ExpertSystemService) SelectBestPesticides(pestisida []dto.PestisidaResponse, limit int) ([]dto.PestisidaResponse, error) {
	if len(pestisida) <= limit {
		return pestisida, nil
	}

	return pestisida[:limit], nil
}

func (s *ExpertSystemService) GetAllPestisida() ([]models.Pestisida, error) {
	var pestisida []models.Pestisida
	err := s.DB.Find(&pestisida).Error
	if err != nil {
		return nil, err
	}
	return pestisida, nil
}

func (s *ExpertSystemService) GetPestisidaByHama(hama string) ([]dto.PestisidaResponse, error) {

	var pestisida []models.Pestisida

	err := s.DB.
		Where("LOWER(hama_sasaran) LIKE LOWER(?)", "%"+hama+"%").
		Find(&pestisida).Error

	if err != nil {
		return nil, err
	}

	var result []dto.PestisidaResponse	

	for _, p := range pestisida {
		result = append(result, dto.PestisidaResponse{
			NamaPestisida:    p.NamaPestisida,
			BahanAktif:       p.BahanAktif,
			BentukFormulasi:  p.BentukFormulasi,
			JenisPestisida:   p.JenisPestisida,
			SasaranKomoditas: p.SasaranKomoditas,
			HamaSasaran:      p.HamaSasaran,
			DosisNilai:       p.DosisNilai,
			DosisSatuan:      p.DosisSatuan,
		})
	}

	return result, nil
}

func (s *ExpertSystemService) CalculateDosage(pestisida []dto.PestisidaResponse, luasLahan float64) ([]dto.PestisidaRecommendation, error) {

	var recommendations []dto.PestisidaRecommendation
	
	for _, p := range pestisida {

		totalDosis := p.DosisNilai * luasLahan

		rec := dto.PestisidaRecommendation{
			NamaPestisida: p.NamaPestisida,
			BahanAktif:    p.BahanAktif,
			DosisPerHa:    p.DosisNilai,
			TotalDosis:    totalDosis,
			Satuan:        p.DosisSatuan,
		}

		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}