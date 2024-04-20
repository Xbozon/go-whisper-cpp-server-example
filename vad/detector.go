package vad

import (
	"log"
	"time"
)

const DefaultQuietTime = time.Millisecond * 1000

type Detector struct {
	lastFlux       float64
	sensitivity    float64
	start          time.Time
	quietTimeDelay time.Duration
	vad            *VAD
}

func NewDetector(sensitivity float64, delay time.Duration, width int) *Detector {
	return &Detector{
		sensitivity:    sensitivity,
		quietTimeDelay: delay,
		vad:            NewVAD(width),
	}
}

func (d *Detector) HearSomething(samples []byte) bool {
	flux := d.vad.Flux(bytesToInt16sLE(samples))

	if d.lastFlux == 0 {
		d.lastFlux = flux * d.sensitivity
		return false
	}

	if flux >= d.lastFlux {
		//log.Println(flux, ">=", d.lastFlux*detectCoefficient)
		d.start = time.Now()
		return true
	}

	if time.Since(d.start) < d.quietTimeDelay {
		log.Println("delay")
		return true
	}

	if flux*d.sensitivity <= d.lastFlux {
		return false
	}

	return false
}

func bytesToInt16sLE(bytes []byte) []int16 {
	// Ensure the byte slice length is even
	if len(bytes)%2 != 0 {
		panic("bytesToInt16sLE: input bytes slice has odd length, must be even")
	}

	int16s := make([]int16, len(bytes)/2)
	for i := 0; i < len(int16s); i++ {
		int16s[i] = int16(bytes[2*i]) | int16(bytes[2*i+1])<<8
	}
	return int16s
}
