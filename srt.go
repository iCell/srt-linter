package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iCell/srt/linter"
	"github.com/urfave/cli"
)

const version = "0.1.0"

func filesFromArgs(args cli.Args) ([]string, error) {
	var files []string
	for _, arg := range args {
		if _, err := os.Stat(arg); err != nil {
			return nil, err
		}

		err := filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func lint(files []string) {
	for _, f := range files {
		lint := linter.NewLinter(f)
		results := lint.Lint()
		if results != nil {
			fmt.Println(".......", f, ".......")
			for _, v := range results {
				fmt.Println("error:", v.Error.Error(), "near line:", v.LineNum)
			}
			fmt.Println(".....................")
		} else {
			fmt.Println(".......", f, ".......")
			fmt.Println("No errors found")
			fmt.Println(".....................")
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "srt"
	app.Usage = "lint srt files"
	app.Version = version
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "iCell",
			Email: "i@icell.io",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "lint",
			Usage: "lint the given files, or the files within the given directory",
			Action: func(c *cli.Context) error {
				files, err := filesFromArgs(c.Args())
				if err != nil {
					return err
				}
				lint(files)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
