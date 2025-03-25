package application

import (
	"strconv"
	"strings"
)

// Consts
var validNotes = map[string]bool{
	"A": true, "B": true, "C": true, "D": true,
	"E": true, "F": true, "G": true, "S": true,
}
var validNotesAttributes = map[string]bool{
	"d": true,
	"a": true,
	"o": true,
}
var validAlterations = map[string]bool{
	"#": true,
	"b": true,
	"n": true,
}

// Methods
func ValidateBPM(bpmStr string, index *int) bool {
	bpm, err := strconv.Atoi(bpmStr)
	if err != nil || bpm <= 0 {
		*index = 0
		return false
	}
	*index += len(bpmStr) + 2 //Si tenemos error en el BPM, el indice se posiciona en 0 y se suma el largo del BPM + el espacio siguiente
	return true
}
func ValidateNote(note string, index *int) bool {
	if !validNotes[string(note[0])] {
		return false
	}
	*index++ //Avanzo posicion por Nota
	noteLen := len(note)
	if noteLen > 1 {
		if string(note[1]) != "{" {
			return false
		}
		*index++ //Avanzo posicion por llave de apertura
		rawAttributes := GetRawAttributes(note, noteLen)
		if !ValidateAttributes(rawAttributes, string(note[0]) == "S", index) {
			return false
		}
		if string(note[noteLen-1]) != "}" {
			return false
		} else {
			*index++ // Sumo por el }
		}
	}
	*index++ //Avanzo posicion por espacio
	return true
}
func GetRawAttributes(note string, noteLen int) string {
	var attributes string
	//Previamente validé opentag por lo que ahora solo valido closing tag para evitar validaciones redundantes
	if string(note[noteLen-1]) != "}" {
		attributes = note[2:]
	} else {
		attributes = note[2 : noteLen-1]
	}
	return attributes
}
func ValidateAttributes(rawAttributes string, isSilence bool, index *int) bool {
	//Valido que el final no sea un ; o } para no perder esos Edge case {o=2;} o {o=2;}}
	lenRawAttr := len(rawAttributes)
	if string(rawAttributes[lenRawAttr-1]) == ";" ||
		string(rawAttributes[lenRawAttr-1]) == "}" {
		*index += lenRawAttr
		return false
	}
	trimmed := strings.TrimSuffix(rawAttributes, ";")
	attributes := strings.Split(trimmed, ";")
	repAttrs := make(map[string]bool)
	for i, attr := range attributes {
		trimmedAttr := strings.TrimSpace(attr)
		if trimmedAttr == "" {
			return false
		}
		keyVals := strings.SplitN(trimmedAttr, "=", 2)
		if len(keyVals) != 2 {
			return false
		}
		key := strings.TrimSpace(keyVals[0])
		value := strings.TrimSpace(keyVals[1])

		// Verifico q no se repitan los key attrs
		if repAttrs[key] {
			return false
		}
		repAttrs[key] = true
		if !ValidateAttrKeys(isSilence, key, index) {
			return false

		}
		switch key {
		case "d":
			if !ValidateDuration(value, index) {
				return false
			}
		case "o":
			if !ValidateOctave(value, index) {
				return false
			}
		case "a":
			if !ValidateAlteration(value, index) {
				return false
			}
		}

		if i < len(attributes)-1 { // Aumento por ;, menos en el ultimo ya que no deberia de tener
			*index++
		}

	}
	return true
}
func ValidateDuration(val string, pos *int) bool {
	// Rechazar notación decimal
	if dotIndex := strings.Index(val, "."); dotIndex > -1 {
		*pos += dotIndex
		return false
	}
	// Caso fraccional
	if strings.Contains(val, "/") {
		parts := strings.Split(val, "/")
		if len(parts) != 2 {
			return false
		}
		num, err1 := strconv.Atoi(parts[0])
		den, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return false
		}
		if den == 0 {
			*pos += len(val) - 1 //la pos actual el largo del value del atributo -1 ej  len("d=1/0") 5 - 1 = 4 => error en d=1/* <--
			return false
		}
		dur := float64(num) / float64(den)
		if dur <= 0 || dur > 4 {
			return false
		}
		*pos += len(val)
		return true

	}
	// Caso entero
	num, err := strconv.Atoi(val)
	if err != nil {
		return false
	}
	if num <= 0 || num > 4 {
		return false
	}
	*pos += len(val) // Acá le sumo el len del value del attr
	return true
}
func ValidateOctave(val string, pos *int) bool {
	oct, err := strconv.Atoi(val)
	if err != nil {
		return false
	}
	if oct < 0 || oct > 8 {
		return false
	}
	*pos++ //sumo la pos de la octava
	return true
}
func ValidateAlteration(val string, pos *int) bool {

	if !validAlterations[val] {
		return false
	}
	*pos++
	return true
}
func ValidateAttrKeys(isSilence bool, key string, index *int) bool {
	if isSilence && key != "d" {
		return false
	}
	if !isSilence && !validNotesAttributes[key] {
		return false
	}
	*index += 2 //Esto seria el Key Attr nas el '='
	return true
}
func GetMelodyFields(melody string) []string {
	tokens := strings.Fields(melody)
	return tokens
}
func ExtractAndFormatAttr(attr string) (string, int) {
	parts := strings.SplitN(attr, "=", 2)
	if len(parts) < 2 {
		return "", -1
	}
	valueParts := strings.Split(parts[1], "=")
	return parts[0] + "=" + valueParts[0][:len(valueParts[0])-2], len(valueParts[0][:len(valueParts[0])-2])
}
func ValidateMelody(melody string) int {
	notes := GetMelodyFields(melody)
	errorIndex := -1
	//Valido BPM
	if !ValidateBPM(notes[0], &errorIndex) {
		return errorIndex
	}
	//Valido Notas
	for _, note := range notes[1:] {
		if !ValidateNote(note, &errorIndex) {
			return errorIndex
		}
	}
	return -1
}
