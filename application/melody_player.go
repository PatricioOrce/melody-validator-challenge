package application

import (
	"fmt"
	"log"
	"melody-validator-challenge/cmd/server/models"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
)

func PlayMelody(body models.PlayMelodyRequest) error {
	const sampleDuration = 48000
	const sampleRate = beep.SampleRate(sampleDuration)
	// Inicializamos el altavoz
	if err := speaker.Init(sampleRate, sampleDuration/10); err != nil {
		log.Printf("Error initializing speaker: %v", err)
		return fmt.Errorf("Error initializing speaker: %v", err)
	}

	// chan pa esperar hasta que termine la reproducción
	ch := make(chan struct{})

	var sounds []beep.Streamer
	for _, note := range body.Notes {
		freq := note.Frequency

		if note.Type == "silence" {
			duration := calculateNoteDuration(note.Duration, body.Tempo.Value)
			sounds = append(sounds, beep.Silence(duration))
			continue
		}
		sineWave, err := generators.SinTone(sampleRate, int(freq))
		if err != nil {
			log.Printf("Error generating sine wave: %v", err)
			return fmt.Errorf("Error generating sine wave: %v", err)

		}
		duration := calculateNoteDuration(note.Duration, body.Tempo.Value)
		sounds = append(sounds, beep.Take(duration, sineWave))
	}

	sounds = append(sounds, beep.Callback(func() {
		ch <- struct{}{}
	}))

	speaker.Play(beep.Seq(sounds...))

	<-ch
	return nil
}

func calculateNoteDuration(noteDuration float64, tempoBPM int) int {
	// Convertimos la duración de la nota  a segundos
	durationInSeconds := (noteDuration * 60) / float64(tempoBPM)
	// Convertimos a la cantidad correspondiente de muestras
	return beep.SampleRate(48000).N(time.Duration(durationInSeconds * float64(time.Second)))
}
