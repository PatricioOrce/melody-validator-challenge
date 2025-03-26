package application

import (
	"melody-validator-challenge/cmd/server/models"
	"melody-validator-challenge/domain"
	"regexp"
	"strconv"
	"strings"
)

var americanToNote = map[string]string{
	"A": "la", "B": "si", "C": "do", "D": "re",
	"E": "mi", "F": "fa", "G": "sol", "S": "silence",
}

func MapMelody(melody string) models.MelodyResponse {
	fields := GetMelodyFields(melody)
	var response models.MelodyResponse

	response.Tempo.Value, _ = strconv.Atoi(fields[0])
	response.Tempo.Unit = "bpm"

	for i := 1; i < len(fields); i++ {
		if fields[i][0] == 'S' {
			response.Notes = append(response.Notes, domain.Silence{
				Type:     americanToNote[string(fields[i][0])],
				Duration: GetDuration(fields[i]),
			})
		} else {
			response.Notes = append(response.Notes, domain.Note{
				Type:       "note",
				Name:       americanToNote[string(fields[i][0])],
				Duration:   GetDuration(fields[i]),
				Octave:     GetOctave(fields[i]),
				Alteration: GetAlteration(fields[i]),
				Frequency:  GetFrequency(string(fields[i][0]), GetAlteration(fields[i]), GetOctave(fields[i])),
			})
		}
	}

	return response
}
func GetAlteration(note string) string {
	if len(note) == 1 {
		return "none"
	}
	re := regexp.MustCompile(`a=([^;}]+)`)
	alterationValue := re.FindStringSubmatch(note)
	if len(alterationValue) < 2 {
		return "none"
	}
	return alterationValue[1]
}
func GetOctave(note string) int {
	if len(note) == 1 {
		return 4
	}
	re := regexp.MustCompile(`o=([^;}]+)`)
	octaveValue := re.FindStringSubmatch(note)
	if len(octaveValue) < 2 {
		return 4
	}
	result, _ := strconv.Atoi(octaveValue[1])

	return result
}
func GetDuration(note string) float64 {
	if len(note) == 1 {
		return 1
	}
	re := regexp.MustCompile(`d=([^;}]+)`)
	durationValue := re.FindStringSubmatch(note)
	if len(durationValue) < 2 {
		return 1
	}

	dur, err := parseDuration(durationValue[1])
	if !err {
		return 1
	}
	return dur
}
func parseDuration(value string) (float64, bool) {
	if strings.Contains(value, "/") {
		parts := strings.Split(value, "/")
		if len(parts) != 2 {
			return 0, false
		}
		numerator, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, false
		}
		denom, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return 0, false
		}
		if denom == 0 {
			return 0, false
		}
		return numerator / denom, true
	}
	duration, _ := strconv.ParseFloat(value, 64)
	return duration, true
}
