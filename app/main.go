package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage

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
		case "echo":
			fmt.Println(strings.Join(args, " "))
		default:
			fmt.Println(strings.TrimSpace(command) + ": command not found")
		}


	}
}
