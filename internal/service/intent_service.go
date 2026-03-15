package service

import (
	"backend/internal/models"
	"regexp"
	"strings"
)

type Intent string

const (
	IntentGreeting Intent = "greeting"
	IntentHelp     Intent = "help"
	IntentThanks   Intent = "thanks"
	IntentRice     Intent = "rice"
	IntentUnknown  Intent = "unknown"
)

type IntentService struct{
	pesticideKeywords []string
}

func NewIntentService(expertService *ExpertSystemService) *IntentService {

	pesticides, err := expertService.GetAllPestisida()
	if err != nil {
		pesticides = []models.Pestisida{}
	}

	var keywords []string

	for _, p := range pesticides {

		name := strings.ToLower(p.NamaPestisida)

		keywords = append(keywords, name)

		// optional: ambil kata pertama juga
		parts := strings.Split(name, " ")
		if len(parts) > 0 {
			keywords = append(keywords, parts[0])
		}
	}

	return &IntentService{
		pesticideKeywords: keywords,
	}
}
func normalizeText(text string) string {

	text = strings.ToLower(text)

	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	text = re.ReplaceAllString(text, "")

	text = strings.TrimSpace(text)

	return text
}

func containsKeyword(text string, keywords []string) bool {

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}

	return false
}

func (s *IntentService) DetectIntent(message string) Intent {

	text := normalizeText(message)

	/// DETECT PESTICIDE NAME
	if containsKeyword(text, s.pesticideKeywords) {
		return IntentRice
	}

	/// GREETING
	greetings := []string{
		"halo", "hai", "hi",
		"selamat pagi", "selamat siang",
		"selamat sore", "selamat malam",
	}

	if containsKeyword(text, greetings) {
		return IntentGreeting
	}

	/// HELP
	helpKeywords := []string{
		"bisa apa",
		"fitur apa",
		"bantuan",
		"help",
		"cara pakai",
	}

	if containsKeyword(text, helpKeywords) {
		return IntentHelp
	}

	/// THANKS
	thanks := []string{
		"terima kasih",
		"makasih",
		"thanks",
	}

	if containsKeyword(text, thanks) {
		return IntentThanks
	}

	/// RICE RELATED
	riceKeywords := []string{
		"padi",
		"sawah",
		"beras",
		"hama",
		"wereng",
		"ulat",
		"penggerek",
		"walang sangit",
		"pestisida",
	}

	if containsKeyword(text, riceKeywords) {
		return IntentRice
	}

	return IntentUnknown
}