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

		cmdArgs := strings.Split(strings.TrimSpace(command), " ")

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
			if file, ok := findBinInPath(args); ok {
				fmt.Println("Entra file ok para handleExe")
				handleExeInPath(file, args)
				return
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

func handleExeInPath(file string, args []string) {
	restArgs := args[1:]
	arg := args[0]

	fmt.Println("cmd:", arg, restArgs)
	if file != "" {
		cmd := exec.Command(arg, restArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
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
	fmt.Println(strings.Join(args, " "))
}