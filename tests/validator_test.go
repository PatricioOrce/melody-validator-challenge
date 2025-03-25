package application_test

import (
	"melody-validator-challenge/application"
	"testing"
)

// TestCase define un caso de prueba para la validación de melodías.
type TestCase struct {
	melody      string // La melodía a validar.
	expectedPos int    // Posición del primer error; -1 indica que es válida.
	description string // Descripción del caso.
}

var testCases = []TestCase{
	{
		melody:      "60 A{d=7/4;o=3;a=#} B{o=2;d=1/2} S A{d=2;a=n} G{a=b} B S{d=1/3} C{d=1} D{d=1}",
		expectedPos: -1, // Válida
		description: "Melodía válida extendida con dos notas adicionales (C y D).",
	},
	{
		melody:      "A{d=1;o=3} C{d=1;o=4} B{d=1} S{d=1} A{d=1}",
		expectedPos: 0, // Error en posición 0: falta BPM.
		description: "Error: falta BPM.",
	},
	{
		melody:      "60 A{d=7/4,o=3;a=#} B{o=2;d=1/2} S A{d=2;a=n} G{a=b} B S{d=1/3} C{d=1} D{d=1}",
		expectedPos: 10, // Error en posición 4: separador incorrecto en A (coma en vez de ;)  --Chequear
		description: "Error: separador incorrecto en A (se usa coma en vez de punto y coma).",
	},
	{
		melody:      "60 C{d=0.3;o=4} D{d=1/2} E{d=1} F{d=1} G{d=1} A{d=1}",
		expectedPos: 8, // Error en posición 4: duración en notación decimal.
		description: "Error: notación decimal en duración (0.3 no permitida).",
	},
	{
		melody:      "60 S{o=2} A{d=1;o=4;a=n} B{d=1/2} S{d=1} C{d=1}",
		expectedPos: 5, // Error en posición 4: silencio con atributo 'o'
		description: "Error: silencio (S) no puede tener atributo 'o'.",
	},
	{
		melody:      "60 G{d=1;o=3;a=#}} B{d=1} S{d=1} A{d=1} B{d=1}",
		expectedPos: 17, // Error en posición 7: llave de cierre extra en G.
		description: "Error: llave de cierre extra en G.",
	},
	{
		melody:      "60 Z{d=1} A{d=1} S{d=1} B{d=1}",
		expectedPos: 3, // Error en posición 3: 'Z' no es nota válida.
		description: "Error: nota inválida ('Z' no es válida).",
	},
	{
		melody:      "60 D{p=f;d=1} C{d=1} B{d=1} S{d=1}",
		expectedPos: 5, // Error en posición 4: atributo desconocido 'p'.
		description: "Error: atributo desconocido ('p') en D.",
	},
	{
		melody:      "60 F{d=1/16 A{d=1} S{d=1} G{d=1}",
		expectedPos: 11, // Error en posición 4: falta llave de cierre en F.
		description: "Error: falta llave de cierre en F.",
	},
	{
		melody:      "60 A{d=7/4 o=3;a=#} B{d=1/2} S{d=1} C{d=1} D{d=1}",
		expectedPos: 10, // Error en posición 8: falta punto y coma entre atributos en A.
		description: "Error: falta punto y coma entre atributos en A.",
	},
	{
		melody:      "60 A{d=7/4;d=1;a=#} B{d=1/2} S{d=1} E{d=1} F{d=1}",
		expectedPos: 11, // Error en posición 4: atributo 'd' duplicado en A.
		description: "Error: atributo 'd' duplicado en A.",
	},
	{
		melody:      "60 A{o=9;d=1} B{d=1/2} S{d=1} C{d=1} D{d=1}",
		expectedPos: 7, // Error en posición 4: octava fuera de rango en A.
		description: "Error: octava fuera de rango en A (9 > 8).",
	},
	{
		melody:      "60 A{d=5;o=4;a=n} B{d=1/2} S{d=1} D{d=1} E{d=1}",
		expectedPos: 7, // Error en posición 4: duración fuera de rango en A.
		description: "Error: duración fuera de rango en A (5 > 4).",
	},
	{
		melody:      "60 A{d=0/3;o=4;a=n} B{d=1/2} S{d=1} F{d=1} G{d=1}",
		expectedPos: 7, // Error en posición 4: duración resulta en 0.
		description: "Error: duración resulta en 0 en A (0/3).",
	},
	{
		melody:      "60 A{d=1/0;o=4;a=n} B{d=1/2} S{d=1} H{d=1} I{d=1}",
		expectedPos: 9, // Error en posición 4: denominador cero en A.
		description: "Error: denominador cero en A (1/0).",
	},
	{
		melody:      "60 A{d=1.5;o=4;a=n} B{d=1/2} S{d=1} L{d=1} M{d=1}",
		expectedPos: 8, // Error en posición 4: notación decimal en duración.
		description: "Error: notación decimal en duración en A ('1.5' no permitida).",
	},
	{
		melody:      "60 A d=1 B{d=1/2} S{d=1} N{d=1} O{d=1}",
		expectedPos: 5, // Error en posición 3: faltan llaves para atributos en A.
		description: "Error: faltan llaves para atributos en A.",
	},
	{
		melody:      "60 A{d=7/4;o=3;a=#} B{d=1/2} S A{d=2;a=n} G{a=b} B S{d=1/3} T{d=1} U{d=1}",
		expectedPos: 60, // Válida extendida
		description: "Melodía invalida con notas adicionales T y U.",
	},
	{
		melody:      "60 A{d=7/4;o=3;a=# B{o=2;d=1/4} S G{d=2} R{d=1} S{d=1}",
		expectedPos: 18, // Error en posición 18: error en bloque de atributos de A.
		description: "Error: bloque de atributos de A mal formado (falta cerrar correctamente, error en pos 18).",
	},
}

func TestValidateMelody(t *testing.T) {
	for i, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got := application.ValidateMelody(tc.melody)
			if tc.expectedPos == -1 {
				// Se espera que la melodía sea válida.
				if got > 0 {
					t.Errorf("Test %d: Para la melodía %q se esperaba que fuera válida (-1) pero se obtuvo error en la posición %d", i+1, tc.melody, got)
				}
			} else {
				// Se espera un error en una posición específica.
				if got != tc.expectedPos {
					t.Errorf("Test %d: Para la melodía %q se esperaba error en la posición %d, pero se obtuvo %d", i+1, tc.melody, tc.expectedPos, got)
				}
			}
		})
	}
}
