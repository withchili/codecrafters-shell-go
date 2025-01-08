package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
		default:
			fmt.Println(command + ": command not found")
		}

		
	}
}
