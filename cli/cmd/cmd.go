package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syfar/parser"
	file "syfar/providers/file"
	http "syfar/providers/http"
	runner "syfar/runner"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "syfar",
		Short: "Syfar CLI",
		Long:  "Syfar is designed for efficient implementation and execution of integration tests.",
	}
	rootCmd.AddCommand(Run)
	rootCmd.AddCommand(Validate)
	Init.LocalFlags().StringP("template", "t", "", "Choose a project template")
	rootCmd.AddCommand(Init)
	return rootCmd
}

var Init = &cobra.Command{
	Use:   "init",
	Short: "Init syfar project",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		err := initProject(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "syfar: Error\n%v\n", err)
		} else {
			fmt.Println("\x1b[32mSyfar project create successfully!\x1b[0m\n\nUse the following to run your:")
			if len(args) == 1 {
				fmt.Printf("\ncd %s\nsyfar run\n", args[0])
			} else {
				fmt.Println("syfar run")
			}
		}

	},
}

var Validate = &cobra.Command{
	Use:   "validate",
	Short: "Validate a syfar file or a project",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, _, _, err := validate(args); err != nil {
			fmt.Fprintf(os.Stderr, "syfar: Error\n%v\n", err)
		} else {
			fmt.Println("Well done, everything is OK")
		}
	},
}

var Run = &cobra.Command{
	Use:   "run",
	Short: "Run a test or a project",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(args); err != nil {
			fmt.Fprintf(os.Stderr, "syfar: Error\n%v\n", err)
		}
	},
}

func initProject(args []string) error {
	var projectDir string

	wdir, err := os.Getwd()
	if err != nil {
		return err
	}

	if len(args) >= 1 {
		projectDir = filepath.Join(wdir, args[0])
	} else {
		projectDir = wdir
	}

	err = initSyfarJson(projectDir)
	if err != nil {
		return err
	}
	return CopyDir("./example/get_started", projectDir)
}

func validate(args []string) (*runner.Syfar, *parser.SyfarFile, context.Context, error) {
	filedir, filename, err := buildPath(args)
	if err != nil {
		return nil, nil, nil, err
	}

	syfar := runner.NewSyfar()

	//File provider
	fileProvider := file.ActionProvider{}
	syfar.RegisterActionProvider("file", &fileProvider)

	//HTTP Provider
	httpProvider := http.ActionProvider{}
	syfar.RegisterActionProvider("http", &httpProvider)

	syfar.Init()

	ast, ctx, err := syfar.Validate(filedir, filename)
	if err != nil {
		return nil, nil, nil, err
	}

	return &syfar, ast, ctx, nil
}

func run(args []string) error {
	syfar, ast, ctx, err := validate(args)
	if err != nil {
		return err
	}
	return syfar.Run(ast, ctx)
}
