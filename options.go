package gocli

type AutoCompleter interface {
	// Readline will pass the whole line and current offset to it
	// Completer need to pass all the candidates, and how long they shared the same characters in line
	// Example:
	//   [go, git, git-shell, grep]
	//   Do("g", 1) => ["o", "it", "it-shell", "rep"], 1
	//   Do("gi", 2) => ["t", "t-shell"], 2
	//   Do("git", 3) => ["", "-shell"], 3
	Do(line []rune, pos int) (newLine [][]rune, length int)
}

type Option func(config *Config) error

func WithConfig(config *Config) Option {
	return func(c *Config) error {
		*c = *config
		return nil
	}
}

func WithCommands(commands []Command) Option {
	return func(c *Config) error {
		c.Commands = commands
		return nil
	}
}

func WithPrompt(prompt string) Option {
	return func(c *Config) error {
		c.Cfg.Prompt = prompt
		return nil
	}
}

func WithHistoryFile(historyFile string) Option {
	return func(c *Config) error {
		c.Cfg.HistoryFile = historyFile
		return nil
	}
}

// WithHistoryLimit specify the max length of historys,
// it's 500 by default, set it to -1 to disable history
func WithHistoryLimit(historyLimit int) Option {
	return func(c *Config) error {
		c.Cfg.HistoryLimit = historyLimit
		return nil
	}
}

func WithDisableAutoSaveHistory(disableAutoSaveHistory bool) Option {
	return func(c *Config) error {
		c.Cfg.DisableAutoSaveHistory = disableAutoSaveHistory
		return nil
	}
}

// WithHistorySearchFold enable case-insensitive history searching
func WithHistorySearchFold(historySearchFold bool) Option {
	return func(c *Config) error {
		c.Cfg.HistorySearchFold = historySearchFold
		return nil
	}
}

// WithAutoComplete AutoCompleter will called once user press TAB
func WithAutoComplete(autoComplete AutoCompleter) Option {
	return func(c *Config) error {
		c.Cfg.AutoComplete = autoComplete
		return nil
	}
}

func WithInterruptPrompt(interruptPrompt string) Option {
	return func(c *Config) error {
		c.Cfg.InterruptPrompt = interruptPrompt
		return nil
	}
}

func WithEOFPrompt(EOFPrompt string) Option {
	return func(c *Config) error {
		c.Cfg.EOFPrompt = EOFPrompt
		return nil
	}
}
