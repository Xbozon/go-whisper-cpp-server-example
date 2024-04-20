package vad

import (
	"fmt"

	"github.com/go-audio/audio"
	"github.com/streamer45/silero-vad-go/speech"
)

type SileroDetector struct {
	detector *speech.Detector
}

func NewSileroDetector(filepath string) (*SileroDetector, error) {
	sd, err := speech.NewDetector(speech.DetectorConfig{
		ModelPath:            filepath,
		SampleRate:           16000,
		WindowSize:           1024,
		Threshold:            0.5,
		MinSilenceDurationMs: 0,
		SpeechPadMs:          0,
	})
	if err != nil {
		return nil, fmt.Errorf("create silero detector: %w", err)
	}

	return &SileroDetector{
		detector: sd,
	}, nil
}

// DetectVoice tries to identify the segment in which the voice is present.
// You can also use a set of segments by iterating over it.
//
//	for _, s := range segments {
//	  log.Printf("speech starts at %0.2fs", s.SpeechStartAt)
//	  if s.SpeechEndAt > 0 {
//	 	log.Printf("speech ends at %0.2fs", s.SpeechEndAt)
//	  }
//	}
func (s *SileroDetector) DetectVoice(buffer *audio.IntBuffer) (bool, error) {
	pcmBuf := buffer.AsFloat32Buffer()

	segments, err := s.detector.Detect(pcmBuf.Data)
	if err != nil {
		return false, fmt.Errorf("detect: %w", err)
	}

	return len(segments) > 0, nil
}
