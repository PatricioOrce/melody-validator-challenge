package models

import (
	"melody-validator-challenge/domain"
)

type ErrorResponse struct {
	Cause string `json:"cause"`
}
type MelodyResponse struct {
	Tempo domain.Tempo  `json:"tempo"`
	Notes []interface{} `json:"notes"` // Puede contener Note o Silence
}
