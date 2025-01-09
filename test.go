package main

import (
	"fmt"
	"strings"
)

func main() {
	cmd := "cat 'dasd asd' 'sdsd'\n"
	command, args := parseCmd(cmd)

	fmt.Println("command: " + command)
	for i, arg := range args {
		fmt.Println(i+1, "arg:", arg)
	}
}

func parseCmd(cmd string) (command string, args []string) {
	var buffer strings.Builder
	var argString string

	for i, char := range cmd {
		if char == ' ' {
			command = buffer.String()
			argString = strings.TrimSpace(cmd[i+1:])
			break
		}
		if char == '\n' {
			command = buffer.String()
			return
		}
		buffer.WriteRune(char)
	}

	if command == "" {
		command = buffer.String()
		return
	}

	buffer.Reset()

	waitingForQuote := false
	for _, char := range argString {
		switch char {
		case '\'':
			waitingForQuote = !waitingForQuote
			if !waitingForQuote {
				args = append(args, buffer.String())
				buffer.Reset()
			}
		case ' ', '\n':
			if waitingForQuote {
				buffer.WriteRune(char)
			} else if buffer.Len() > 0 {
				args = append(args, buffer.String())
				buffer.Reset()
			}
		default:
			buffer.WriteRune(char)
		}
	}

	if buffer.Len() > 0 {
		args = append(args, buffer.String())
	}

	return
}
