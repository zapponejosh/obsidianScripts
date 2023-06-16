package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := "/Users/joshzappone/Library/Mobile Documents/iCloud~md~obsidian/Documents/personal"

	destDir := "/Users/joshzappone/Library/Mobile Documents/iCloud~md~obsidian/Documents/personal/cleanup"

	// Define and parse the command-line flag
	searchString := flag.String("search", "", "String to search for in the first 10 lines of .md files")
	deleteUntitled := flag.Bool("untitled", false, "Delete empty Untitled.md files")
	flag.Parse()

	if *searchString == "" && !*deleteUntitled {
		fmt.Println("Please provide a search string using the -search flag or use -untitled to delete empty Untitled.md files")
		return
	}

	if *deleteUntitled {
		err := deleteUntitledFiles(dir)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Untitled files deleted successfully.")
		}
	}

	if *searchString != "" {
		matchingFiles := findFilesWithString(dir, *searchString)

		err := moveFiles(matchingFiles, destDir)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("File movement completed successfully.")
		}
	}

}

func findFilesWithString(dirPath, searchString string) []string {
	var matchingFiles []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only consider .md files
		if filepath.Ext(path) == ".md" {
			matchFound, err := searchFirstLines(path, searchString)
			if err != nil {
				return err
			}

			if matchFound {
				matchingFiles = append(matchingFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return matchingFiles
}

func searchFirstLines(filePath, searchString string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()

		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		words := strings.Fields(line)
		for _, word := range words {
			if word == searchString {
				return true, nil
			}
		}

		if lineCount >= 10 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func moveFiles(files []string, destinationDirectory string) error {
	// Create the destination directory if it doesn't exist
	err := os.MkdirAll(destinationDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := filepath.Base(file)
		destinationPath := filepath.Join(destinationDirectory, fileName)

		// Move the file to the destination directory
		err = os.Rename(file, destinationPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteUntitledFiles(directoryPath string) error {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file name starts with "Untitled" and ends with ".md"
		if strings.HasPrefix(info.Name(), "Untitled") && strings.HasSuffix(info.Name(), ".md") && info.Size() == 0 {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
