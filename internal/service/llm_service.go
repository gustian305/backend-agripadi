package service

import (
	"backend/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	// "strings"
	"time"
)

type LLMService struct {
	Client  *http.Client
	APIKey  string
	BaseURL string
}

func NewLLMService() *LLMService {

	return &LLMService{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey:  config.Cfg.GROQAPIKEY,
		BaseURL: "https://api.groq.com/openai/v1/chat/completions",
	}
}

// Request Struct
type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
}

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response Struct
type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func greetingMessage() string {
	return `
Halo! 👋

Saya adalah AI asisten pertanian padi pada aplikasi ini.

Saya dapat membantu Anda:
• Mengidentifikasi hama padi dari gambar
• Memberikan rekomendasi pestisida
• Menjelaskan cara penggunaan pestisida
• Memberikan informasi budidaya padi

Silakan kirim pertanyaan tentang padi atau foto hama tanaman padi Anda.
`
}

func helpMessage() string {
	return `
Anda dapat menggunakan aplikasi ini untuk:

1️⃣ Mengirim foto hama padi untuk dideteksi AI  
2️⃣ Mendapat rekomendasi pestisida otomatis  
3️⃣ Bertanya tentang budidaya padi  

Contoh pertanyaan:
• "Apa penyebab daun padi menguning?"
• "Bagaimana mengatasi wereng?"
`
}

func thanksMessage() string {
	return "Sama-sama! Jika ada pertanyaan tentang tanaman padi, silakan tanyakan."
}

// func isGreeting(message string) bool {

// 	m := strings.ToLower(strings.TrimSpace(message))

// 	greetings := []string{
// 		"halo",
// 		"hai",
// 		"hi",
// 		"selamat pagi",
// 		"selamat siang",
// 		"selamat sore",
// 		"selamat malam",
// 	}

// 	for _, g := range greetings {
// 		if m == g {
// 			return true
// 		}
// 	}

// 	return false
// }

// func appIntroductionMessage() string {
// 	return `
// Halo! 👋

// Saya adalah AI asisten pertanian padi pada aplikasi ini.

// Saya dapat membantu Anda untuk:
// 1. Mengidentifikasi hama padi dari gambar menggunakan AI
// 2. Memberikan rekomendasi pestisida berdasarkan sistem pakar
// 3. Menjelaskan cara penggunaan pestisida dengan aman
// 4. Memberikan pengetahuan tentang budidaya padi

// Silakan kirim:
// • Pertanyaan tentang padi
// • Atau kirim foto hama pada tanaman padi Anda

// Saya siap membantu! 🌾
// `
// }

// func (s *LLMService) IsRiceRelated(question string) bool {

// 	q := strings.ToLower(question)

// 	keywords := []string{
// 		"padi",
// 		"beras",
// 		"sawah",
// 		"tanaman padi",
// 		"hama padi",
// 		"wereng",
// 		"penggerek batang",
// 		"ulat grayak",
// 		"walang sangit",
// 		"pestisida padi",
// 		"budidaya padi",
// 		"irigasi sawah",
// 		"panen padi",
// 	}

// 	for _, k := range keywords {
// 		if strings.Contains(q, k) {
// 			return true
// 		}
// 	}

// 	return false
// }

func (s *LLMService) callLLM(messages []GroqMessage) (string, error) {
	reqBody := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: messages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		s.BaseURL,
		bytes.NewBuffer(bodyBytes),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.APIKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("groq request failed: " + string(respBytes))
	}

	var result GroqResponse

	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no response from groq")
	}

	return result.Choices[0].Message.Content, nil
}

func riceSystemPrompt() string {
	return `
Anda adalah asisten AI pertanian yang hanya membahas tanaman padi.

ATURAN:
- Hanya menjawab pertanyaan tentang tanaman padi
- Tidak menjawab pertanyaan tentang tanaman lain
- Jika pertanyaan di luar topik padi, katakan bahwa sistem hanya mendukung tanaman padi
- Gunakan bahasa sederhana yang mudah dipahami petani
- Jawaban harus singkat dan praktis tapi mengedukasi
`
}

func (s *LLMService) GenerateRiceKnowledge(
	question string,
) (string, error) {

	messages := []GroqMessage{
		{
			Role:    "system",
			Content: riceSystemPrompt(),
		},
		{
			Role:    "user",
			Content: question,
		},
	}

	return s.callLLM(messages)
}

func formatPesticideList(pestisida []string) string {

	if len(pestisida) == 0 {
		return "Tidak ada pestisida yang direkomendasikan."
	}

	var result string

	for i, p := range pestisida {
		result += fmt.Sprintf("%d. %s\n", i+1, p)
	}

	return result
}


func (s *LLMService) GeneratePesticideExplanation(
	hama string,
	luasLahan float64,
	fasePadi string,
	pestisida []string,
) (string, error) {

	pesticideList := formatPesticideList(pestisida)

	prompt := fmt.Sprintf(`
Anda adalah penyuluh pertanian yang membantu petani padi.

Gunakan bahasa sederhana dan praktis seperti penyuluh lapangan.
Jawaban harus singkat, jelas, dan mudah dipahami petani.

=============================
KONDISI SAWAH
=============================
Hama terdeteksi : %s
Luas lahan      : %.2f hektar
Fase padi       : %s

=============================
REKOMENDASI PESTISIDA
=============================
%s

=============================
FORMAT JAWABAN WAJIB
=============================

Tulis jawaban dengan format berikut:

🦗 **Bahaya Hama**
Jelaskan secara singkat mengapa hama ini berbahaya bagi tanaman padi.

🌿 **Pestisida yang Direkomendasikan**
Jelaskan fungsi pestisida yang direkomendasikan.

🧪 **Cara Mencampur Pestisida**
Jelaskan langkah mencampur pestisida dengan air.

🚜 **Perkiraan Kebutuhan Larutan**
Perkirakan kebutuhan larutan semprot untuk %.2f hektar sawah.

🌾 **Cara Penyemprotan yang Efektif**
Berikan tips penyemprotan yang efektif.

⏰ **Waktu Penyemprotan Terbaik**
Jelaskan kapan waktu terbaik melakukan penyemprotan.

⚠️ **Tips Keselamatan**
Berikan tips keselamatan saat menggunakan pestisida.

Gunakan bullet point agar mudah dibaca petani.
`, hama, luasLahan, fasePadi, pesticideList, luasLahan)

	messages := []GroqMessage{
		{
			Role:    "system",
			Content: riceSystemPrompt(),
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return s.callLLM(messages)
}