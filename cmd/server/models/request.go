package models

type MelodyRequest struct {
	Melody string `json:"melody"`
}
type PlayMelodyRequest struct {
	Tempo Tempo  `json:"tempo"`
	Notes []Note `json:"notes"`
}
