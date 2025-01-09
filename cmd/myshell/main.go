package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtinCommands = map[string]bool{
	"type": true,
	"exit": true,
	"echo": true,
	"pwd":  true,
	"cd":   true,
}

func exitCommand(arguments []string) {
	if arguments[0] == "0" {
		os.Exit(0)
	}
}

func echoCommand(arguments []string) {
	for _, arg := range arguments {
		fmt.Print(arg, " ")
	}
	fmt.Print("\n")
}

func inPath(fileName string) (string, bool) {
	pathDirectories := strings.Split(os.Getenv("PATH"), ":")

	for _, dir := range pathDirectories {
		fullPath := dir + "/" + fileName
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, true
		}
	}

	return "", false
}

func typeCommand(arguments []string) {
	arg := arguments[0]

	if _, ok := builtinCommands[arg]; ok {
		fmt.Println(arg, "is a shell builtin")
		return
	}

	if fullPath, ok := inPath(arg); ok {
		fmt.Println(arg, "is", fullPath)
		return
	}

	fmt.Println(arg + ": not found")
}

func pwdCommand() {
	dir, _ := os.Getwd()
	fmt.Println(dir)
}

func cdCommand(arguments []string) {
	if len(arguments) == 0 {
		return
	}

	path := arguments[0]

	if path == "~" {
		_ = os.Chdir(os.Getenv("HOME"))
		return
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Println("cd: " + path + ": No such file or directory")
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmdString, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		// cmdString = strings.TrimSuffix(cmdString, "\n")
		// cmdStringParts := strings.Split(cmdString, " ")

		command, arguments := parseCmd(cmdString)

		out, err := exec.Command(command, arguments...).Output()
		if err == nil {
			fmt.Print(string(out))
			continue
		}

		switch command {
		case "exit":
			exitCommand(arguments)
		case "echo":
			echoCommand(arguments)
		case "type":
			typeCommand(arguments)
		case "pwd":
			pwdCommand()
		case "cd":
			cdCommand(arguments)
		default:
			fmt.Println(command + ": command not found")
		}

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
	for i, char := range argString {
		switch char {
		case '\'':
			waitingForQuote = !waitingForQuote
			if !waitingForQuote && i+1 < len(argString)-1 && rune(argString[i+1]) != '\'' {
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
