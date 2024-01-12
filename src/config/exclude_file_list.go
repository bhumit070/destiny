package config

var ExcludedFileList = map[string]string{
	// node related files
	"package.json":      "package.json",
	"package-lock.json": "package-lock.json",
	"yarn.lock":         "yarn.lock",
	"node_modules":      "node_modules",
	"npm-debug.log":     "npm-debug.log",
	"yarn-error.log":    "yarn-error.log",
	"pnpm-debug.log":    "pnpm-debug.log",
	"pnpm-lock.yaml":    "pnpm-lock.yaml",
	"pnpmfile.js":       "pnpmfile.js",

	// git
	".gitignore": ".gitignore",

	// zsh
	".zshrc":       ".zshrc",
	".zsh_history": ".zsh_history",
	".zprofile":    ".zprofile",
	".zcompdump":   ".zcompdump",

	// bash
	".bashrc":       ".bashrc",
	".bash_history": ".bash_history",
	".bash_profile": ".bash_profile",

	// vim
	".vimrc":   ".vimrc",
	".viminfo": ".viminfo",
	".vim":     ".vim",

	// tmux
	".tmux.conf":       ".tmux.conf",
	".tmux.conf.local": ".tmux.conf.local",
	".tmux":            ".tmux",
	".tmuxinator":      ".tmuxinator",

	// go
	"go.mod": "go.mod",
	"go.sum": "go.sum",

	// Linux
	"Makefile": "Makefile",
}
