package repository

import (
	"fmt"
	"sort"
	"sync"

	"example.com/videolab/internal/model"
)

type Memory struct {
	mu      sync.RWMutex
	clips   map[string]model.Clip
	presets map[string]model.Preset
}

func NewMemory(clips []model.Clip, presets []model.Preset) *Memory {
	store := &Memory{
		clips:   make(map[string]model.Clip, len(clips)),
		presets: make(map[string]model.Preset, len(presets)),
	}
	for _, clip := range clips {
		store.clips[clip.ID] = clip
	}
	for _, preset := range presets {
		store.presets[preset.Name] = preset
	}
	return store
}

func (m *Memory) ListClips() []model.Clip {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]model.Clip, 0, len(m.clips))
	for _, clip := range m.clips {
		result = append(result, clip)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (m *Memory) GetClip(id string) (model.Clip, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	clip, ok := m.clips[id]
	return clip, ok
}

func (m *Memory) SavePreset(preset model.Preset) error {
	if preset.Name == "" || preset.VideoFilter == "" {
		return fmt.Errorf("preset name and video filter are required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.presets[preset.Name] = preset
	return nil
}

func (m *Memory) ListPresets() []model.Preset {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]model.Preset, 0, len(m.presets))
	for _, preset := range m.presets {
		result = append(result, preset)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
