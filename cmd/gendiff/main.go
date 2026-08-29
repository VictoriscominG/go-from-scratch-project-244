package main

import (
	"code"
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
				Usage:   "output format (stylish, plain, json)",
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

			str, err := code.GenDiff(path1, path2, format)
			if err != nil {
				return err
			}

			fmt.Println(str)
			return nil
		},
	}

	if err := command.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
