package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cygnetdigital/roachdiff/pkg/diff"
	"github.com/cygnetdigital/roachdiff/pkg/git"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "roachdiff",
		Usage:   "diff a cockroachdb sql schema and produce the migration steps",
		Version: "dev",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "gitref",
				Usage: "a git reference (branch, tag, commit) where the schema exists in the previous state",
			},
			&cli.StringFlag{
				Name:    "previous",
				Aliases: []string{"prev"},
				Usage:   "a path to the previous version of the schema",
			},
		},
		ArgsUsage: "[path_to_schema]",
		Commands:  []*cli.Command{},
		Action: func(c *cli.Context) error {
			gitref := c.String("gitref")
			prev := c.String("previous")
			schemaPath := c.Args().First()

			if gitref == "" && prev == "" {
				gitref = "main" // default to main branch
			}

			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get pwd: %w", err)
			}

			currentSchema, err := os.ReadFile(resolvePath(pwd, schemaPath))
			if err != nil {
				return fmt.Errorf("failed to read current schema: %w", err)
			}

			if gitref != "" {
				prevSchema, err := git.OpenFile(pwd, schemaPath, gitref)
				if err != nil {
					return fmt.Errorf("failed to open git schema: %w", err)
				}

				return invokeDiff(
					diff.NewDiffer(string(prevSchema), string(currentSchema)),
				)
			}

			if prev != "" {
				prevSchema, err := os.ReadFile(resolvePath(pwd, prev))
				if err != nil {
					return fmt.Errorf("failed to read previous schema: %w", err)
				}

				return invokeDiff(
					diff.NewDiffer(string(prevSchema), string(currentSchema)),
				)
			}

			return fmt.Errorf("invalid arguments")
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func invokeDiff(df *diff.Differ) error {
	df.Generator.Warnings = true

	result, err := df.Run()
	if err != nil {
		return fmt.Errorf("failed to diff: %w", err)
	}

	fmt.Print(result.String())
	return nil
}

func resolvePath(pwd string, p string) string {
	if path.IsAbs(p) {
		return p
	}

	return path.Join(pwd, p)
}
