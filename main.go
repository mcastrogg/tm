package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"time"
)

const dateFormat = "20060102"

var (
	editFlag bool
)

func init() {
	flag.BoolVar(&editFlag, "e", false, "Open the most recent file in the editor")
	flag.Parse()
}

func main() {

	directory := "/home/matias/Dropbox/vault/tasks/"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	var matchingFiles []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Use a regular expression to match files with the specified date format in the name
		match, err := regexp.MatchString(`^\d{8}\.md$`, file.Name())
		if err != nil {
			fmt.Println("Error matching file name:", err)
			os.Exit(1)
		}

		if match {
			matchingFiles = append(matchingFiles, file.Name())
		}
	}

	if len(matchingFiles) == 0 {
		fmt.Println("No matching files found.")
		os.Exit(0)
	}

	// Sort the matching files by modification time
	sort.Slice(matchingFiles, func(i, j int) bool {
		timeI, _ := time.Parse(dateFormat+".md", matchingFiles[i])
		timeJ, _ := time.Parse(dateFormat+".md", matchingFiles[j])
		return timeI.After(timeJ)
	})

	// Get the most recent file
	mostRecentFile := matchingFiles[0]
	mostRecentFilePath := fmt.Sprintf("%s/%s", directory, mostRecentFile)

	fmt.Println("Most recent file:", mostRecentFilePath)

	if editFlag {
		err = openInTerminalEditor(mostRecentFilePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
	}
	openPreview(mostRecentFilePath)
}

func openPreview(filePath string) error {
	cmd := exec.Command("glow", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func openInTerminalEditor(filePath string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", os.Getenv("EDITOR"), filePath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

