package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func getLogFiles() (logFileNames []string, err error) {
	directory := "./"
	files, err := os.Open(directory)
	if err != nil {
		return logFileNames, errors.New("error opening directory")
	}
	defer files.Close()

	fileInfos, err := files.ReadDir(-1)
	if err != nil {
		errStr := fmt.Sprintf("error reading directory: %v", err) //if directory is not read properly print error message
		return logFileNames, errors.New(errStr)
	}

	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()
		if strings.Contains(fileName, ".log") && strings.Contains(fileName, "2024-06") {
			logFileNames = append(logFileNames, fileInfo.Name())
		}
	}

	return logFileNames, err
}

func analyzeFile(fileName string) (requestsCount int, err error) {
	fmt.Printf("Analyzing '%s' ....\n", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return requestsCount, err
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		lineText := s.Text()
		if strings.Contains(lineText, "HTTP") {
			requestsCount += 1
		}
	}

	return requestsCount, nil
}

func main() {
	logFileNames, err := getLogFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	sum := 0
	for _, logFileName := range logFileNames {
		c, _ := analyzeFile(logFileName)
		sum += c
	}

	fmt.Println("\n======")
	fmt.Println("Number of files: ", len(logFileNames))
	fmt.Println("Number of lines: ", sum)
}
