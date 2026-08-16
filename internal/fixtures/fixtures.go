package fixtures

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"example.com/videolab/internal/model"
)

//go:embed clips.json
var clipsData []byte

//go:embed presets.json
var presetsData []byte

func Clips() ([]model.Clip, error) {
	var clips []model.Clip
	if err := json.Unmarshal(clipsData, &clips); err != nil {
		return nil, fmt.Errorf("load clip fixtures: %w", err)
	}
	return clips, nil
}

func Presets() ([]model.Preset, error) {
	var presets []model.Preset
	if err := json.Unmarshal(presetsData, &presets); err != nil {
		return nil, fmt.Errorf("load preset fixtures: %w", err)
	}
	return presets, nil
}
