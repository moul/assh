package config

import "github.com/urfave/cli"

// SSHFlags is built at init(), it contains urfave/cli string & bool flags for SSH
var SSHFlags = []cli.Flag{}

// SSHBoolFlags contains list of available SSH boolean options
var SSHBoolFlags = []string{"1", "2", "4", "6", "A", "a", "C", "f", "G", "g", "K", "k", "M", "N", "n", "q", "s", "T", "t", "V", "v", "X", "x", "Y", "y"}

// SSHStringFlags contains list of available SSH string options
var SSHStringFlags = []string{"b", "c", "D", "E", "e", "F", "I", "i", "L", "l", "m", "O", "o", "p", "Q", "R", "S", "W", "w"}

func init() {
	// Populate SSHFlags
	// FIXME: support slice flags (-O a -O b -O c === []string{"a", "b", "c"}
	// FIXME: support count flags (-vvv == -v -v -v)
	// FIXME: support joined bool flags (-it == -i -t)
	for _, flag := range SSHBoolFlags {
		SSHFlags = append(SSHFlags, cli.BoolFlag{
			Name: flag,
		})
	}
	for _, flag := range SSHStringFlags {
		SSHFlags = append(SSHFlags, cli.StringFlag{
			Name:  flag,
			Value: "",
		})
	}
}
