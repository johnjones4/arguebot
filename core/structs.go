package core

type Agent struct {
	Intention string `json:"intention"`
}

type TranscriptEntry struct {
	Agent int    `json:"agent"`
	Text  string `json:"text"`
}
