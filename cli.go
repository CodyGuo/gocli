package gocli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"

	cliV1 "github.com/urfave/cli"
	cliV2 "github.com/urfave/cli/v2"

	"github.com/chzyer/readline"
)

type Command struct {
	Name        string
	SubCommands []Command
}

type Config struct {
	Cfg      *readline.Config
	Commands []Command
}

type Cli struct {
	mu       sync.Mutex
	config   *Config
	preRun   func(args []string)
	postRun  func(args []string)
	preRunE  func(args []string) error
	postRunE func(args []string) error
}

func New(options ...Option) *Cli {
	cfg := &readline.Config{
		Prompt:          ">>> ",
		InterruptPrompt: "^C",
	}
	conf := &Config{
		Cfg: cfg,
	}
	for _, o := range options {
		o(conf)
	}
	return &Cli{config: conf}
}

func (c *Cli) Config() *Config {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.config
}

func (c *Cli) SetConfig(config *Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = config
}

func (c *Cli) Commands() []Command {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.config.Commands
}

func (c *Cli) SetCommands(commands []Command) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Commands = commands
}

func (c *Cli) PreRun(preRun func(args []string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.preRun = preRun
}

func (c *Cli) PostRun(postRun func(args []string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.postRun = postRun
}

func (c *Cli) run(app interface{}, args []string) error {
	switch cmd := app.(type) {
	case *cobra.Command:
		return cobraRun(cmd, args)
	default:
		return cliRun(app, args)
	}
}

func (c *Cli) Run(name string, app interface{}) error {
	args := os.Args
	if c.preRun != nil {
		c.preRun(args)
	}
	if c.preRunE != nil {
		if err := c.preRunE(args); err != nil {
			return err
		}
	}
	defer func() error {
		if c.postRun != nil {
			c.postRun(args)
		}
		if c.postRunE != nil {
			if err := c.postRunE(args); err != nil {
				return err
			}
		}
		return nil
	}()

	instance, err := readline.NewEx(c.config.Cfg)
	if err != nil {
		return err
	}
	defer instance.Close()
	for {
		line, err := instance.Readline()
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				break
			}
			fmt.Printf("Run failed, error: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cmd := strings.Fields(line)
		_commands := append([]string{name}, cmd...)
		err = c.run(app, _commands)
		if err != nil && !errors.Is(err, readline.ErrInterrupt) {
			return err
		}
	}
	return nil
}

func (c *Cli) ParseCommands(commander interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch _commands := commander.(type) {
	case []cliV1.Command:
		c.config.Commands = parseUrFaveCliCommandsV1(_commands)
	case []*cliV2.Command:
		c.config.Commands = parseUrFaveCliCommandsV2(_commands)
	case *cobra.Command:
		c.config.Commands = parseUrCobraCommands(_commands.Commands())
	default:
		return fmt.Errorf("ParseCommands failed, not support: %T\n", _commands)
	}
	completes := parseAutoComplete(c.config.Commands)
	if len(completes) > 0 {
		c.config.Cfg.AutoComplete = readline.NewPrefixCompleter(completes...)
	}
	return nil
}

func parseAutoComplete(commands []Command) []readline.PrefixCompleterInterface {
	completes := []readline.PrefixCompleterInterface{}
	for _, command := range commands {
		sbcCompletes := parseAutoComplete(command.SubCommands)
		complete := readline.PcItem(command.Name, sbcCompletes...)
		completes = append(completes, complete)
	}
	return completes
}

func parseUrFaveCliCommandsV1(commands cliV1.Commands) []Command {
	_commands := []Command{}
	for _, c := range commands {
		subCommands := parseUrFaveCliCommandsV1(c.Subcommands)
		_command := Command{Name: c.Name, SubCommands: subCommands}
		for _, flag := range c.Flags {
			name := "-" + strings.Split(flag.GetName(), ",")[0]
			flagCommand := Command{Name: name, SubCommands: subCommands}
			subCommands = append([]Command{flagCommand}, subCommands...)
			_command = Command{Name: c.Name, SubCommands: subCommands}
		}
		_commands = append(_commands, _command)
	}
	return _commands
}

func parseUrFaveCliCommandsV2(commands []*cliV2.Command) []Command {
	_commands := []Command{}
	for _, c := range commands {
		subCommands := parseUrFaveCliCommandsV2(c.Subcommands)
		_command := Command{Name: c.Name, SubCommands: subCommands}
		for _, flag := range c.Flags {
			if len(flag.Names()) == 0 {
				continue
			}
			name := "-" + flag.Names()[0]
			flagCommand := Command{Name: name, SubCommands: subCommands}
			subCommands = append([]Command{flagCommand}, subCommands...)
			_command = Command{Name: c.Name, SubCommands: subCommands}
		}
		_commands = append(_commands, _command)
	}
	return _commands
}

func parseUrCobraCommands(commands []*cobra.Command) []Command {
	_commands := []Command{}
	for _, c := range commands {
		subCommands := parseUrCobraCommands(c.Commands())
		_command := Command{Name: c.Name(), SubCommands: subCommands}
		flags := c.Flags()
		flags.VisitAll(func(flag *pflag.Flag) {
			_name := "--" + flag.Name
			flagCommand := Command{Name: _name, SubCommands: subCommands}
			subCommands = append([]Command{flagCommand}, subCommands...)
			_command = Command{Name: c.Name(), SubCommands: subCommands}
		})
		_commands = append(_commands, _command)
	}
	return _commands
}

func cliRun(command interface{}, args []string) error {
	skips := []string{"Usage", "argument", "invalid"}
	switch cmd := command.(type) {
	case *cliV1.App:
		if err := cmd.Run(args); err != nil {
			return skipError(err, skips...)
		}
	case *cliV2.App:
		if err := cmd.Run(args); err != nil {
			return skipError(err, skips...)
		}
	}
	return nil
}

func cobraRun(cmd *cobra.Command, args []string) error {
	skips := []string{"argument", "requires", "unknown"}
	cmd.SetArgs(args[1:])
	err := cmd.Execute()
	if err != nil {
		return skipError(err, skips...)
	}
	return nil
}

func skipError(err error, v ...string) error {
	for _, s := range v {
		if strings.Contains(err.Error(), s) {
			return nil
		}
	}
	return err
}

func MinimumNArgs(args []string, n int) error {
	if len(args) < n {
		return fmt.Errorf("requires at least %d arg(s), only received %d", n, len(args))
	}
	return nil
}
