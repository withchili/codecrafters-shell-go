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

func parseCmd(cmd string) (command string, args []string, err error) {
	var stringBuffer strings.Builder
	var remainingCmd string

	if cmd == "" {
		return "", nil, nil
	}

	for _, rn := range cmd {
		if rn == ' ' {
			command = stringBuffer.String()
			remainingCmd = strings.TrimPrefix(cmd, command+" ")
			break
		}
		if rn == '\n' {
			command = stringBuffer.String()
			return
		}
		stringBuffer.WriteRune(rn)
	}

	stringBuffer.Reset()

	remainingCmd = strings.TrimSuffix(remainingCmd, "\n")

	expectingDoubleQuote := false
	expectingSingleQuote := false
	backslashIsOn := false

	for _, rn := range remainingCmd {
		if backslashIsOn {
			stringBuffer.WriteRune(rn)
			backslashIsOn = false
			continue
		}

		switch rn {
		case '"':
			if !expectingSingleQuote {
				expectingDoubleQuote = !expectingDoubleQuote
				continue
			}
			stringBuffer.WriteRune(rn)
		case '\'':
			if !expectingDoubleQuote {
				expectingSingleQuote = !expectingSingleQuote
				continue
			}
			stringBuffer.WriteRune(rn)
		case '\\':
			if backslashIsOn {
				stringBuffer.WriteRune(rn)
				backslashIsOn = false
			} else if expectingSingleQuote || expectingDoubleQuote {
				stringBuffer.WriteRune(rn)
			} else {
				backslashIsOn = true
			}
		case ' ':
			if backslashIsOn {
				stringBuffer.WriteRune(rn)
				backslashIsOn = false
			} else if expectingDoubleQuote || expectingSingleQuote {
				stringBuffer.WriteRune(rn)
			} else {
				if stringBuffer.Len() > 0 {
					args = append(args, stringBuffer.String())
					stringBuffer.Reset()
				}
			}
		default:
			if backslashIsOn {
				stringBuffer.WriteRune('\\')
				backslashIsOn = false
			}
			stringBuffer.WriteRune(rn)
		}
	}

	if expectingDoubleQuote || expectingSingleQuote {
		return "", nil, fmt.Errorf("unmatched quote in input string")
	}

	if stringBuffer.Len() > 0 {
		args = append(args, stringBuffer.String())
	}

	return command, args, nil
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmdString, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		command, arguments, _ := parseCmd(cmdString)

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
			cmd := exec.Command(command, arguments...)
			cmd.Stdout = os.Stdout
			if err := cmd.Run(); err != nil {
				fmt.Println(command + ": command not found")
			}
		}
	}
}
