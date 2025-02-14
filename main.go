package main

import (
	"flag"
	"fmt"
	"log"
	"ok/cmd/run"
	"ok/cmd/version"
	"os"
	"sort"
)

type command interface {
	Run(args []string)
	Description() string
}

var commands = map[string]command{
	"run":     &run.Command{},
	"version": &version.Command{},
}

func main() {
	flag.Usage = func() {
		fmt.Println("ok is a tool for managing ok source code.")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Println("\tok <command> [arguments]")
		fmt.Println("")
		fmt.Println("The commands are:")
		fmt.Println("")

		var keys []string
		for name := range commands {
			keys = append(keys, name)
		}

		sort.Strings(keys)

		for _, name := range keys {
			fmt.Printf("\t%s\t\t%s\n", name, commands[name].Description())
		}

		fmt.Println("")
		fmt.Println(`Use "ok <command> -help" for more information about a command.`)
		fmt.Println("")
	}

	// This is just for the -help (usage).
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalln("missing command")
	}

	cmdName := os.Args[1]
	if cmd, ok := commands[cmdName]; ok {
		cmd.Run(os.Args[2:])
		return
	}

	log.Fatalln("unknown command:", cmdName)
}
