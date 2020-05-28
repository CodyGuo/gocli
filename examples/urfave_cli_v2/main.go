package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/CodyGuo/gocli"
	cli "github.com/urfave/cli/v2"
)

var (
	app *cli.App
)

func main() {
	initCmd()
	options := []gocli.Option{
		gocli.WithPrompt(app.Name + " >>> "),
		gocli.WithHistoryFile("./.history_" + app.Name),
	}
	shell := gocli.New(options...)
	shell.PreRun(func(args []string) {
		fmt.Println("           _____  _____  ")
		fmt.Println("     /\\   |  __ \\|  __ \\ ")
		fmt.Println("    /  \\  | |__) | |__) |")
		fmt.Println("   / /\\ \\ |  ___/|  ___/ ")
		fmt.Println("  / ____ \\| |    | |     ")
		fmt.Println(" /_/    \\_\\_|    |_|     ")
		fmt.Println("                         ")
		fmt.Println("                         ")
	})
	shell.PostRun(func(args []string) {
		fmt.Println("exit...")
	})
	if err := shell.ParseCommands(app.Commands); err != nil {
		log.Fatal(err)
	}
	app.Action = func(ctx *cli.Context) error {
		if ctx.NArg() == 0 {
			return shell.Run(app.Name, app)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initCmd() {
	var commands = []*cli.Command{
		{
			Name:    "exit",
			Aliases: []string{"quit"},
			Usage:   "exit",
			Action: func(ctx *cli.Context) error {
				fmt.Println("Bye")
				os.Exit(0)
				return nil
			},
		},
		{
			Name:      "print",
			ArgsUsage: "print [string to print]",
			Usage:     "Print anything to the screen",
			UsageText: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
			Action: func(ctx *cli.Context) error {
				args := ctx.Args().Slice()
				if err := gocli.MinimumNArgs(args, 1); err != nil {
					fmt.Println(err)
					return nil
				}
				fmt.Println("Print: " + strings.Join(args, " "))
				return nil
			},
		},
		{
			Name:      "echo",
			ArgsUsage: "echo [string to echo]",
			Usage:     "Echo anything to the screen",
			UsageText: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
			Action: func(ctx *cli.Context) error {
				args := ctx.Args().Slice()
				if err := gocli.MinimumNArgs(args, 1); err != nil {
					fmt.Println(err)
					return nil
				}
				fmt.Println("Echo: " + strings.Join(args, " "))
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name:      "times",
					Aliases:   []string{"t"},
					ArgsUsage: "times [string to echo]",
					Usage:     "Echo anything to the screen more times",
					UsageText: `echo things multiple times back to the user by providing
a count and a string.`,
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:    "times",
							Aliases: []string{"t"},
							Value:   1,
							Usage:   "times to echo the input",
						},
					},
					Action: func(ctx *cli.Context) error {
						args := ctx.Args().Slice()
						if err := gocli.MinimumNArgs(args, 1); err != nil {
							fmt.Println(err)
							return nil
						}
						echoTimes := ctx.Int("times")
						for i := 0; i < echoTimes; i++ {
							fmt.Println("Echo: " + strings.Join(args, " "))
						}
						return nil
					},
				},
			},
		},
	}
	app = cli.NewApp()
	app.Name = "app"
	app.Commands = commands
}
