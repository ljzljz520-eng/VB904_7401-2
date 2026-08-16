package ffmpeg

import (
	"context"
	"testing"
)

type successRunner struct{}

func (successRunner) Run(_ context.Context, args []string) ([]byte, error) {
	if len(args) == 0 {
		return nil, context.Canceled
	}
	return nil, nil
}

func TestCommandArgsKeepWindowsPathAsOneArgument(t *testing.T) {
	command := Command{VideoPath: `C:\Video Lab\clip.mp4`, OutputPath: `C:\Video Lab\frame.jpg`, Filter: "eq=brightness=0.1"}
	args := command.Args()
	if args[4] != command.VideoPath {
		t.Fatalf("unexpected input path: %q", args[4])
	}
	if command.WindowsString() == "" {
		t.Fatal("expected command")
	}
}

func TestExecutorAcceptsSuccessfulRun(t *testing.T) {
	if err := NewExecutor(successRunner{}).ExtractFrame(context.Background(), Command{VideoPath: "clip.mp4", OutputPath: "frame.jpg", Filter: "null"}); err != nil {
		t.Fatal(err)
	}
}
