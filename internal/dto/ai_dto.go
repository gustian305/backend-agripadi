package dto

type CNNPredictionResponse struct {
	Success    bool    `json:"success"`
	Prediction string  `json:"prediction"`
	Confidence float64 `json:"confidence"`
}


type AIContext struct {
	Hama       string  `json:"hama"`
	Confidence float64 `json:"confidence"`
	LuasLahan  float64 `json:"luas_lahan"`
	FaseTanaman string `json:"fase_tanaman"`
}

type AIChatResponse struct {
	Chat string `json:"chat"`
	Hama string `json:"hama"`
	Confidence float64 `json:"confidence"`
	
	PestisidaRecommendation []PestisidaRecommendation `json:"pestisida_recomendation"`
}