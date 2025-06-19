package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {

	mapCommands := map[string]string {
		"echo":"is a shell builtin",
		"exit":"is a shell builtin",
		"type":"is a shell builtin",
	}

	for {
		fmt.Fprint(os.Stdout, "$ ")
	
		// Wait for user input
		command,err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("error", err.Error())
			os.Exit(1)
		}

		// cmdArgs := strings.Split(strings.TrimSpace(command), " ")
		cmdArgs, err := parseArgs(command)
		if err != nil {
			fmt.Println("error parsing arguments", err)
			continue
		}

		cmd := cmdArgs[0]
		args := cmdArgs[1:]

		switch cmd {
		case "exit":
			exitCommand(args)
		case "echo":
			echoCommand(args)
		case "type":
			typeCommand(args, mapCommands)
		default:
			if file, ok := findExeInPath(cmd); ok {
				handleExeInPath(cmd, args, file)
			} else {
				fmt.Println(strings.TrimSpace(command) + ": command not found")
			}
		}
	}
}

func typeCommand (args []string, mapCommands map[string]string) {
	arg := args[0]

	if desc, ok := mapCommands[arg]; ok {
		fmt.Println(arg, desc)
		return
	}

	if file, exists := findBinInPath(args); exists {
		fmt.Printf("%s is %s \n", arg, file)
		return
	}
	fmt.Println(arg + ": not found")

}

func handleExeInPath(cmd string, args []string, file string) {

	if file != "" {
		execCmd := exec.Command(cmd, args...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Run()
	}
	return
}

func findExeInPath(cmd string) (string, bool)  {
	paths := os.Getenv("PATH")

	for _, path := range strings.Split(paths, ":") {
		file := path + "/" + cmd

		if _,err := os.Stat(file); err == nil {
			return file, true
		}
	}
	return "", false
}

func findBinInPath(args []string) (string, bool)  {
	bin := args[0]

	paths := os.Getenv("PATH")
	
	for _, path := range strings.Split(paths, ":") {
		file := path + "/" + bin

		if _,err := os.Stat(file); err == nil {
			return file, true
		}
	}
	return "", false
}

func exitCommand (args []string) {
	returnCode := 0
	if len(args) > 0 {
		returnCode2, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(os.Stderr, "Error reading 'exit' command argument", err)
			os.Exit(1)
		}
		returnCode = returnCode2
	}
	os.Exit(returnCode)
}

func echoCommand (args []string) {
	unquoted := make([]string, 0, len(args))

	for _, arg := range args {
		// Remove surrounding single or double quotes if present
		if len(arg) >= 2 {
			if (arg[0] == '"' && arg[len(arg)-1] == '"') || (arg[0] == '\'' && arg[len(arg)-1] == '\'') {
				arg = arg[1 : len(arg)-1]
			}
		}
		unquoted = append(unquoted, arg)
	}

	fmt.Println(strings.Join(unquoted, " "))
}

func parseArgs(input string) ([]string , error){
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	words := []string{}
	current := ""
	inQuote := false
	var quoteChar rune

	for _, ch := range input {
		switch ch {
		case ' ', '\t':
			if inQuote {
				current += string(ch)
			} else if current != "" {
				words = append(words, current)
				current = ""
			}
		case '\'', '"':
			if inQuote && ch == quoteChar {
				inQuote = false
				words = append(words, current)
				current = ""
			} else if !inQuote {
				inQuote = true
				quoteChar = ch
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}

	if current != "" {
		words = append(words, current)
	}

	if inQuote {
		return nil, fmt.Errorf("uncloused quote in input")
	}
	return words, nil
}