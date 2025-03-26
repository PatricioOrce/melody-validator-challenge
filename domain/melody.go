package domain

type Tempo struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type Note struct {
	Type       string  `json:"type"`
	Name       string  `json:"name"`
	Octave     int     `json:"octave"`
	Alteration string  `json:"alteration"`
	Duration   float64 `json:"duration"`
	Frequency  float64 `json:"frequency"`
}

type Silence struct {
	Type     string  `json:"type"`
	Duration float64 `json:"duration"`
}
