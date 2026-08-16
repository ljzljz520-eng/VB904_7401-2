package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"

	"example.com/videolab/internal/model"
)

type Runner interface {
	Run(context.Context, []string) ([]byte, error)
}

type OSRunner struct{}

func (OSRunner) Run(ctx context.Context, args []string) ([]byte, error) {
	var stderr bytes.Buffer
	command := exec.CommandContext(ctx, "ffmpeg", args...)
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return stderr.Bytes(), err
	}
	return stderr.Bytes(), nil
}

type Command struct {
	VideoPath  string
	OutputPath string
	Filter     string
}

func (c Command) Args() []string {
	return []string{
		"-y",
		"-ss", "00:00:01",
		"-i", c.VideoPath,
		"-vf", c.Filter,
		"-frames:v", "1",
		c.OutputPath,
	}
}

func (c Command) WindowsString() string {
	args := c.Args()
	quoted := make([]string, len(args)+1)
	quoted[0] = "ffmpeg"
	for i, arg := range args {
		quoted[i+1] = quoteWindows(arg)
	}
	return strings.Join(quoted, " ")
}

func quoteWindows(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
}

type Executor struct {
	runner Runner
}

func NewExecutor(runner Runner) *Executor {
	return &Executor{runner: runner}
}

func (e *Executor) ExtractFrame(ctx context.Context, command Command) error {
	if e.runner == nil {
		return errors.New("ffmpeg runner is not configured")
	}
	output, err := e.runner.Run(ctx, command.Args())
	if err != nil {
		detail := strings.TrimSpace(string(output))
		return &model.DecodeError{
			VideoPath: command.VideoPath,
			Detail:    detail,
			Cause:     err,
		}
	}
	return nil
}
