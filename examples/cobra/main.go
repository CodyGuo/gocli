package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/CodyGuo/gocli"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func main() {
	initCmd()
	options := []gocli.Option{
		gocli.WithPrompt(rootCmd.Name() + " >>> "),
		gocli.WithHistoryFile("./.history_" + rootCmd.Name()),
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
	if err := shell.ParseCommands(rootCmd); err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 1 {
		shell.Run(rootCmd.Name(), rootCmd)
	} else {
		rootCmd.Execute()
	}
}

func initCmd() {
	var echoTimes int

	var cmdPrint = &cobra.Command{
		Use:   "print [string to print]",
		Short: "Print anything to the screen",
		Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + strings.Join(args, " "))
		},
	}

	var cmdTimes = &cobra.Command{
		Use:   "times [string to echo]",
		Short: "Echo anything to the screen more times",
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}

	var cmdExit = &cobra.Command{
		Use:   "exit",
		Short: "exit",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Bye")
			os.Exit(0)
		},
	}

	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdPrint, cmdEcho, cmdExit)
	cmdEcho.AddCommand(cmdTimes)
}
