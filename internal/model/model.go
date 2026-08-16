package model

import "fmt"

type Clip struct {
	ID          string `json:"id"`
	Scene       string `json:"scene"`
	VideoPath   string `json:"video_path"`
	FrameOutput string `json:"frame_output"`
}

func (c Clip) FrameOutputFor(presetName string) string {
	if len(c.FrameOutput) < len(".jpg") {
		return c.FrameOutput + "-" + presetName
	}
	return c.FrameOutput[:len(c.FrameOutput)-len(".jpg")] + "-" + presetName + ".jpg"
}

type Preset struct {
	Name        string `json:"name"`
	VideoFilter string `json:"video_filter"`
	Description string `json:"description"`
}

type Extraction struct {
	ClipID     string
	PresetName string
	VideoPath  string
	OutputPath string
	Filter     string
	Command    string
}

type ComparisonReport struct {
	Clips       []Clip
	Presets     []Preset
	Extractions []Extraction
}

type DecodeError struct {
	VideoPath string
	Detail    string
	Cause     error
}

func (e *DecodeError) Error() string {
	if e.Detail == "" {
		return fmt.Sprintf("ffmpeg could not decode %s", e.VideoPath)
	}
	return fmt.Sprintf("ffmpeg could not decode %s: %s", e.VideoPath, e.Detail)
}

func (e *DecodeError) Unwrap() error { return e.Cause }
