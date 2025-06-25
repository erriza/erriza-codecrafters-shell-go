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
		cmdArgs := parseCommand(strings.TrimSpace(command))
		
		if len(cmdArgs) == 0 {
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
	if len(args) == 0 {
		return
	}

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
	if len(args) == 0 {
		return "", false
	}

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

func parseCommand(input string) []string {
    var result []string
    var current strings.Builder
    var quoteChar rune = 0 // 0 means not in quotes, '\'' or '"' means in quotes
    
    for i := 0; i < len(input); i++ {
        char := rune(input[i])
        
        if quoteChar == 0 {
            // Not currently in quotes
            if char == '\'' || char == '"' {
                // Start of quoted section
                quoteChar = char
            } else if char == ' ' {
                // Space outside quotes - end current argument
                if current.Len() > 0 {
                    result = append(result, current.String())
                    current.Reset()
                }
                // Skip consecutive spaces
                for i+1 < len(input) && input[i+1] == ' ' {
                    i++
                }
            } else {
                // Regular character outside quotes
                current.WriteRune(char)
            }
        } else {
            // Currently in quotes
            if char == quoteChar {
                // End of quoted section
                quoteChar = 0
            } else {
                // Character inside quotes (treated literally)
                current.WriteRune(char)
            }
        }
    }
    
    // Add the last argument if there is one
    if current.Len() > 0 {
        result = append(result, current.String())
    }
    
    return result
}