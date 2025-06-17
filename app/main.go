package main

import (
	"bufio"
	"fmt"
	"os"
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
		command = strings.TrimSpace(command)
		fmt.Println(command + ": command not found")
	}
}
