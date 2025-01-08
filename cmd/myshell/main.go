package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtinCommands = map[string]bool{
	"type": true,
	"exit": true,
	"echo": true,
}

func exit(arguments []string) {
	if arguments[0] == "0" {
		os.Exit(0)
	}
}

func echo(arguments []string) {
	for _, arg := range arguments {
		fmt.Print(arg, " ")
	}
	fmt.Print("\n")
}

func typeCommand(arguments []string) {
	arg := arguments[0]
	if _, ok := builtinCommands[arg]; ok {
		fmt.Println(arg, "is a shell builtin")
	} else {
		fmt.Println(arg + ": not found")
	}
}

func main() {
	

	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		cmdString, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		cmdString = strings.TrimSuffix(cmdString, "\n")
		cmdStringParts := strings.Split(cmdString, " ")
		
		
		command := cmdStringParts[0]
		arguments := cmdStringParts[1:]
		
		switch command {
		case "exit":
			exit(arguments)
		case "echo":
			echo(arguments)
		case "type":
			typeCommand(arguments)
		default:
			fmt.Println(command + ": command not found")
		}

		
	}
}
