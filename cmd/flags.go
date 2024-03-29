package cmd

import (
	"flag"
)

type Flags struct {
	ConfigFile string
	Version    bool
}

func ParseFlags() Flags {
	config := flag.String("config-file", "~/.hermes/config.toml", "path to configuration file")
	version := flag.Bool("version", false, "print version")
	flag.Parse()
	return Flags{
		ConfigFile: *config,
		Version:    *version,
	}
}
