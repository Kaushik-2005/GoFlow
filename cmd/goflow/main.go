package main

import (
	"fmt"
	"os"
)

func requireArg(args []string, index int, message string) (string, bool) {
	if len(args) <= index {
		fmt.Println(message)
		return "", false
	}
	return args[index], true
}

func printHelp() {
	fmt.Println("Usage: goflow <command>")
	fmt.Println("Commands: list, create, get, process")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "list":
		fmt.Println("listing jobs")
	case "create":
		jobType, ok := requireArg(os.Args, 2, "missing job type")
		if !ok {
			return
		}
		fmt.Printf("creating job of type: %s\n", jobType)
	case "get":
		jobID, ok := requireArg(os.Args, 2, "missing job id")
		if !ok {
			return
		}
		fmt.Printf("getting job: %s\n", jobID)
	case "process":
		jobID, ok := requireArg(os.Args, 2, "missing job id")
		if !ok {
			return
		}
		fmt.Printf("processing job: %s\n", jobID)
	default:
		fmt.Printf("unknown command: %s\n", command)
	}
}
