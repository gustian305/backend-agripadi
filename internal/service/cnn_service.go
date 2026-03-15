package service

import (
	"backend/internal/dto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type CNNServiceInterface interface {
	DetectHama(imagePath string) (*dto.CNNPredictionResponse, error)
}

type CNNService struct {
	BaseURL string
	Client  *http.Client
}

func NewCNNService() *CNNService {
	return &CNNService{
		BaseURL: "http://192.168.1.9:8001",
		Client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (s *CNNService) DetectHama(imagePath string) (*dto.CNNPredictionResponse, error) {

	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", s.BaseURL+"/predict", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cnn service error: %s", resp.Status)
	}

	var cnnResp dto.CNNPredictionResponse

	err = json.NewDecoder(resp.Body).Decode(&cnnResp)
	if err != nil {
		return nil, err
	}

	result := dto.CNNPredictionResponse{
		Success:    cnnResp.Success,
		Prediction: cnnResp.Prediction,
		Confidence: cnnResp.Confidence,
	}

	return &result, nil
}
