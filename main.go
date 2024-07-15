package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

func getMinuteSecond(lineText string) (mnt int, sec int) {
	re := regexp.MustCompile(`\d{2}:(\d{2}):(\d{2})`)
	x := re.FindStringSubmatch(lineText)

	if len(x) == 3 {
		mnt, _ = strconv.Atoi(x[1])
		sec, _ = strconv.Atoi(x[2])
	}

	return mnt, sec
}

func analyzeFile(fileName string) (requestsCount int, maxReqPerSec int, err error) {
	fmt.Printf("Analyzing '%s' ....\n", fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return requestsCount, maxReqPerSec, err
	}
	defer file.Close()

	maxReqPerSec = 0
	curReqPerSec := 0
	prevSec := 0
	prevMin := 0

	s := bufio.NewScanner(file)
	for s.Scan() {
		lineText := s.Text()

		curMin, curSec := getMinuteSecond(lineText)
		if curSec == prevSec && curMin == prevMin {
			curReqPerSec += 1
		} else {
			if curReqPerSec > maxReqPerSec {
				maxReqPerSec = curReqPerSec
			}
			curReqPerSec = 1
			prevSec = curSec
			prevMin = curMin
		}

		if strings.Contains(lineText, "HTTP") {
			requestsCount += 1
		}
	}

	return requestsCount, maxReqPerSec, nil
}

func main() {
	logFileNames, err := getLogFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	maxPerSec := 0

	sum := 0
	for _, logFileName := range logFileNames {
		x, curMaxPerSecond, _ := analyzeFile(logFileName)
		sum += x
		if curMaxPerSecond > maxPerSec {
			maxPerSec = curMaxPerSecond
		}
	}

	fmt.Println("\n======")
	fmt.Println("Number of files: ", len(logFileNames))
	fmt.Println("Number of lines: ", sum)
	fmt.Println("Max per second: ", maxPerSec)
}
