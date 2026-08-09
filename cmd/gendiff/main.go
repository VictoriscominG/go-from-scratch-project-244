package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"code/internal/diff"
	myparser "code/internal/parser"

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
				Usage:   "output format",
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

			fmt.Sprintf("Format: %s\n", format)

			config1, _ := myparser.ParseFile(path1)
			config2, _ := myparser.ParseFile(path2)

			str, _ := diff.DiffFile(config1, config2)
			fmt.Println(str)

			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
