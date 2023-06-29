package requests

type TextToSpeechRequest struct {
	VideoQuoteId string  `json:"videoQuoteId"`
	Voice        string  `json:"voice"`
	Text         string  `json:"text"`
	Seed         float64 `json:"seed"`
	Speed        float64 `json:"speed"`
}
