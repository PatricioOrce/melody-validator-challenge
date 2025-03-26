package application

import (
	"math"
)

var noteValues = map[string]int{
	"A": 9, "B": 11, "C": 0, "D": 2,
	"E": 4, "F": 5, "G": 7,
}
var altValues = map[string]int{
	"b": -1, "none": 0, "#": 1,
}

func GetUniquePosition(name string, alteration string, octave int) int {
	return Pos(name) + Alt(alteration) + (12 * octave)
}

func Pos(note string) int {
	return noteValues[note]
}
func Alt(alt string) int {
	return altValues[alt]
}

func CalculateFrequency(uniquePosition int) float64 {
	pow := float64(uniquePosition-57) / 12
	return 440 * math.Pow(2, pow)
}

func GetFrequency(name string, alteration string, octave int) float64 {
	return math.Round(CalculateFrequency(GetUniquePosition(name, alteration, octave))*100) / 100
}
