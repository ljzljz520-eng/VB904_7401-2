package service

import (
	"context"
	"fmt"

	"example.com/videolab/internal/ffmpeg"
	"example.com/videolab/internal/model"
)

type Store interface {
	ListClips() []model.Clip
	GetClip(string) (model.Clip, bool)
	SavePreset(model.Preset) error
	ListPresets() []model.Preset
}

type Workflow struct {
	store    Store
	executor *ffmpeg.Executor
}

func NewWorkflow(store Store, executor *ffmpeg.Executor) *Workflow {
	return &Workflow{store: store, executor: executor}
}

func (w *Workflow) BuildComparison() (model.ComparisonReport, error) {
	clips := w.store.ListClips()
	presets := w.store.ListPresets()
	if len(clips) == 0 {
		return model.ComparisonReport{}, fmt.Errorf("no clips are available")
	}
	if len(presets) == 0 {
		return model.ComparisonReport{}, fmt.Errorf("no color presets are available")
	}
	report := model.ComparisonReport{
		Clips:   append([]model.Clip(nil), clips...),
		Presets: append([]model.Preset(nil), presets...),
	}
	for _, clip := range clips {
		for _, preset := range presets {
			command := ffmpeg.Command{
				VideoPath:  clip.VideoPath,
				OutputPath: clip.FrameOutputFor(preset.Name),
				Filter:     preset.VideoFilter,
			}
			report.Extractions = append(report.Extractions, model.Extraction{
				ClipID:     clip.ID,
				PresetName: preset.Name,
				VideoPath:  clip.VideoPath,
				OutputPath: command.OutputPath,
				Filter:     preset.VideoFilter,
				Command:    command.WindowsString(),
			})
		}
	}
	return report, nil
}

func (w *Workflow) Extract(ctx context.Context, clipID, presetName string) error {
	clip, ok := w.store.GetClip(clipID)
	if !ok {
		return fmt.Errorf("clip %q was not found", clipID)
	}
	var preset model.Preset
	for _, candidate := range w.store.ListPresets() {
		if candidate.Name == presetName {
			preset = candidate
			break
		}
	}
	if preset.Name == "" {
		return fmt.Errorf("preset %q was not found", presetName)
	}
	command := ffmpeg.Command{
		VideoPath:  clip.VideoPath,
		OutputPath: clip.FrameOutputFor(preset.Name),
		Filter:     preset.VideoFilter,
	}
	if err := w.executor.ExtractFrame(ctx, command); err != nil {
		return fmt.Errorf("extract frame for %s with %s: %w", clip.Scene, preset.Name, err)
	}
	return nil
}

func (w *Workflow) SavePreset(preset model.Preset) error {
	return w.store.SavePreset(preset)
}
