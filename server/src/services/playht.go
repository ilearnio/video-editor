package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"videoeditor/src/config"
)

type PlayHtTextToSpeechRequest struct {
	Quality      string  `json:"quality,omitempty"`
	OutputFormat string  `json:"output_format,omitempty"`
	SampleRate   int64   `json:"sample_rate,omitempty"`
	Voice        string  `json:"voice"`
	Text         string  `json:"text"`
	Seed         float64 `json:"seed"`
	Speed        float64 `json:"speed"`
}

type PlayHtTextToSpeechResponse struct {
	ID       string  `json:"id"`
	Progress int64   `json:"progress"`
	Stage    string  `json:"stage"`
	URL      string  `json:"url"`
	Duration float64 `json:"duration"`
	Size     float64 `json:"size"`
}

func playHtParseTTSStreamResponse(response string) (*PlayHtTextToSpeechResponse, error) {
	lines := strings.Split(strings.TrimRight(response, "\r\n"), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	lastLine := lines[len(lines)-1]
	if len(lastLine) < 6 {
		return nil, fmt.Errorf("invalid response format")
	}
	result := &PlayHtTextToSpeechResponse{}
	err := json.Unmarshal([]byte(lastLine[6:]), result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	if result.Stage != "complete" {
		return nil, fmt.Errorf("TTS generation failed. Response:\n %v", response)
	}
	return result, nil
}

func PlayHtTextToSpeech(params PlayHtTextToSpeechRequest) (*PlayHtTextToSpeechResponse, error) {
	url := "https://play.ht/api/v2/tts"
	if params.Quality == "" {
		params.Quality = "premium"
	}
	if params.OutputFormat == "" {
		params.OutputFormat = "mp3"
	}
	if params.SampleRate == 0 {
		params.SampleRate = 32000
	}
	if params.Seed == 0 {
		params.Seed = 0
	}

	headers := map[string]string{
		"Accept":        "text/event-stream",
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + config.EnvConfig.PLAYHT_SECRET,
		"X-User-Id":     config.EnvConfig.PLAYHT_USER_ID,
	}

	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return playHtParseTTSStreamResponse(string(body))
}
