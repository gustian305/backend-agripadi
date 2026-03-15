package utils

// import (
// 	"backend/internal/dto"
// 	"fmt"
// 	"strings"
// )

// type PromptBuilder struct{}

// func NewPromptBuilder() *PromptBuilder {
// 	return &PromptBuilder{}
// }

// /*
// BuildDiagnosePrompt digunakan ketika sistem menemukan pestisida
// */
// func (p *PromptBuilder) BuildDiagnosePrompt(
// 	hama string,
// 	luasLahan float64,
// 	result *dto.ExpertResult,
// ) string {

// 	return fmt.Sprintf(`
// Sistem mendeteksi hama pada tanaman.

// DATA DIAGNOSA
// -------------
// Hama yang terdeteksi: %s
// Luas lahan petani: %.2f hektar

// REKOMENDASI SISTEM
// -------------
// Nama Pestisida: %s
// Dosis penggunaan: %.2f %s per liter air
// Total kebutuhan pestisida: %.2f %s
// Total kebutuhan air: %.2f liter

// TUGAS ANDA
// -------------
// Jelaskan kepada petani dengan bahasa sederhana:

// 1. Fungsi pestisida tersebut
// 2. Cara mencampur pestisida dengan air
// 3. Cara penyemprotan yang benar
// 4. Waktu penyemprotan terbaik
// 5. Hal-hal yang harus diperhatikan agar aman bagi petani dan tanaman

// Gunakan bahasa yang mudah dipahami oleh petani.
// `,
// 		hama,
// 		luasLahan,
// 		result.NamaPestisida,
// 		result.DosisPerLiter,
// 		result.SatuanDosis,
// 		result.TotalDosis,
// 		result.SatuanDosis,
// 		result.KebutuhanAir,
// 	)
// }

// /*
// BuildNoPesticidePrompt digunakan jika hama tidak ditemukan dalam dataset
// */
// func (p *PromptBuilder) BuildNoPesticidePrompt(hama string) string {

// 	return fmt.Sprintf(`
// Sistem mendeteksi hama pada tanaman.

// DATA DIAGNOSA
// -------------
// Hama yang terdeteksi: %s

// Pestisida untuk hama ini tidak tersedia dalam database sistem.

// TUGAS ANDA
// -------------
// Berikan saran penanganan hama secara umum kepada petani.

// Jelaskan:
// 1. Cara mengendalikan hama secara manual
// 2. Cara pencegahan agar hama tidak menyebar
// 3. Cara menjaga kesehatan tanaman
// 4. Tips perawatan tanaman agar tidak mudah terserang hama

// Gunakan bahasa sederhana agar mudah dipahami petani.
// `, hama)
// }

// func (p *PromptBuilder) CleanInput(input string) string {
// 	return strings.TrimSpace(strings.ToLower(input))
// }