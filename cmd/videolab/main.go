package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"example.com/videolab/internal/ffmpeg"
	"example.com/videolab/internal/fixtures"
	"example.com/videolab/internal/model"
	"example.com/videolab/internal/repository"
	"example.com/videolab/internal/service"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	clips, err := fixtures.Clips()
	if err != nil {
		return err
	}
	presets, err := fixtures.Presets()
	if err != nil {
		return err
	}
	store := repository.NewMemory(clips, presets)
	workflow := service.NewWorkflow(store, ffmpeg.NewExecutor(ffmpeg.OSRunner{}))
	if len(args) == 0 || args[0] == "compare" {
		return printComparison(workflow)
	}
	if args[0] == "extract" {
		return runExtract(workflow, args[1:])
	}
	if args[0] == "preset" {
		return runPreset(workflow, args[1:])
	}
	return fmt.Errorf("unknown command %q", args[0])
}

func printComparison(workflow *service.Workflow) error {
	report, err := workflow.BuildComparison()
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func runExtract(workflow *service.Workflow, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: videolab extract CLIP_ID PRESET")
	}
	return workflow.Extract(context.Background(), args[0], args[1])
}

func runPreset(workflow *service.Workflow, args []string) error {
	flags := flag.NewFlagSet("preset", flag.ContinueOnError)
	name := flags.String("name", "", "preset name")
	filter := flags.String("filter", "", "ffmpeg video filter")
	description := flags.String("description", "", "preset description")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if err := workflow.SavePreset(model.Preset{Name: *name, VideoFilter: *filter, Description: *description}); err != nil {
		return err
	}
	fmt.Printf("saved preset %q\n", *name)
	return nil
}
