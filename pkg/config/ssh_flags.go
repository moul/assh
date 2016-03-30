package config

import "github.com/codegangsta/cli"

var SSHFlags = []cli.Flag{}
var SSHBoolFlags = []string{"1", "2", "4", "6", "A", "a", "C", "f", "G", "g", "K", "k", "M", "N", "n", "q", "s", "T", "t", "V", "v", "X", "x", "Y", "y"}
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
