package dto

type PestisidaResponse struct {
	NamaPestisida    string  `json:"nama_pestisida"`
	BahanAktif       string  `json:"bahan_aktif"`
	BentukFormulasi  string  `json:"bentuk_formulasi"`
	JenisPestisida   string  `json:"jenis_pestisida"`
	SasaranKomoditas string  `json:"sasaran_komoditas"`
	HamaSasaran      string  `json:"hama_sasaran"`
	DosisNilai       float64 `json:"dosis_nilai"`
	DosisSatuan      string  `json:"dosis_satuan"`
}

type PestisidaRecommendation struct {
	NamaPestisida string  `json:"nama_pestisida"`
	BahanAktif    string  `json:"bahan_aktif"`
	DosisPerHa    float64 `json:"dosis_per_ha"`
	TotalDosis    float64 `json:"total_dosis"`
	Satuan        string  `json:"satuan"`
}