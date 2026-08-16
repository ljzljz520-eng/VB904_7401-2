package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"example.com/videolab/internal/ffmpeg"
	"example.com/videolab/internal/fixtures"
	"example.com/videolab/internal/model"
	"example.com/videolab/internal/repository"
)

type recordingRunner struct {
	args   []string
	output []byte
	err    error
}

func (r *recordingRunner) Run(_ context.Context, args []string) ([]byte, error) {
	r.args = append([]string(nil), args...)
	return r.output, r.err
}

func newWorkflow(t *testing.T, runner ffmpeg.Runner) *Workflow {
	t.Helper()
	clips, err := fixtures.Clips()
	if err != nil {
		t.Fatal(err)
	}
	presets, err := fixtures.Presets()
	if err != nil {
		t.Fatal(err)
	}
	return NewWorkflow(repository.NewMemory(clips, presets), ffmpeg.NewExecutor(runner))
}

func TestBuildComparisonUsesStableWindowsCommands(t *testing.T) {
	workflow := newWorkflow(t, &recordingRunner{})
	report, err := workflow.BuildComparison()
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Extractions) != 9 {
		t.Fatalf("got %d extractions", len(report.Extractions))
	}
	want := `ffmpeg "-y" "-ss" "00:00:01" "-i" "C:\Video Lab\fixtures\seaside beach.mp4" "-vf" "eq=contrast=1.04:brightness=0.02:saturation=1.05" "-frames:v" "1" "C:\Video Lab\frames\seaside-clean.jpg"`
	var actual string
	for _, extraction := range report.Extractions {
		if extraction.ClipID == "seaside" && extraction.PresetName == "clean" {
			actual = extraction.Command
			break
		}
	}
	if diff := cmp.Diff(want, actual); diff != "" {
		t.Fatal(diff)
	}
}

func TestSavePresetAddsDeterministicOption(t *testing.T) {
	workflow := newWorkflow(t, &recordingRunner{})
	if err := workflow.SavePreset(model.Preset{Name: "warm", VideoFilter: "eq=saturation=1.2"}); err != nil {
		t.Fatal(err)
	}
	report, err := workflow.BuildComparison()
	if err != nil {
		t.Fatal(err)
	}
	if report.Presets[len(report.Presets)-1].Name != "warm" {
		t.Fatalf("unexpected last preset: %q", report.Presets[len(report.Presets)-1].Name)
	}
	if len(report.Extractions) != 12 {
		t.Fatalf("got %d extractions", len(report.Extractions))
	}
}

func TestCompareReportsDecodeFailureDetails(t *testing.T) {
	path := `C:\Video Lab\fixtures\seaside beach.mp4`
	detail := "Invalid data found when processing input"
	runner := &recordingRunner{output: []byte(detail), err: errors.New("exit status 1")}
	workflow := newWorkflow(t, runner)
	err := workflow.Extract(context.Background(), "seaside", "clean")
	if err == nil {
		t.Fatal("expected extraction to fail")
	}
	if !strings.Contains(err.Error(), path) {
		t.Fatalf("error does not include video path: %v", err)
	}
	if !strings.Contains(err.Error(), detail) {
		t.Fatalf("error does not include ffmpeg detail: %v", err)
	}
}
