package main

import (
	mydiff "code/internal/diff"
	myformatter "code/internal/formatters"
	myparser "code/internal/parser"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	command := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format (stylish, plain)",
				Value:   "stylish",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := cmd.Args().Slice() // получаем аргументы коммандной строки
			if len(args) != 2 {
				return (cli.ShowAppHelp(cmd))
			}
			path1 := args[0]
			path2 := args[1]

			format := cmd.String("format")

			config1, err := myparser.ParseFile(path1)
			if err != nil {
				return err
			}
			config2, err := myparser.ParseFile(path2)
			if err != nil {
				return err
			}

			diffResult, err := mydiff.DiffFile(config1, config2)
			if err != nil {
				return err
			}

			str := myformatter.Formatters(diffResult, format)
			fmt.Println(str)

			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
