package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName, timeLimit := setFlags()

	lines := readFile(csvFileName)

	questions := parseLines(lines)

	rightAnswers := startQuiz(questions, timeLimit)

	fmt.Printf("Result: %d/%d.\n", rightAnswers, len(questions))
}

type question struct {
	question string
	answer   string
}

func setFlags() (*string, *int) {
	csvFileName := flag.String("csv", "questions.csv", "a csv file in (question,answer) format")
	timeLimit := flag.Int("limit", 5, "a time limit for every question in seconds")
	flag.Parse()
	return csvFileName, timeLimit
}

func readFile(filename *string) [][]string {
	file, err := os.Open(*filename)

	if err != nil {
		exit(fmt.Sprintf("Failed to open file \"%s\"\n", *filename))
	}

	// read all the lines from .csv file
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to read provided file\n")
	}
	return lines
}

// parses csv file to question struct array
func parseLines(lines [][]string) []question {
	quesions := make([]question, len(lines))

	for i, line := range lines {
		quesions[i] = question{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return quesions
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func startQuiz(questions []question, timeLimit *int) int {
	rightAnswers := 0

	// start quiz timer
	defer elapsed()()

	for i, q := range questions {
		fmt.Printf("Question #%d: %s = ", i+1, q.question)
		// start timer for each question
		timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
		answerChannel := make(chan string)

		// read input async
		go func() {
			answer := ""
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		// timer is done
		case <-timer.C:
			fmt.Printf("\nYou run out of time.\n")
			return rightAnswers
		// user typed an answer
		case answer := <-answerChannel:
			if answer == q.answer {
				rightAnswers++
				fmt.Printf("Right!\n\n")
			} else {
				fmt.Printf("Wrong.\n\n")
			}
		}
	}

	return rightAnswers
}

// total quiz time
func elapsed() func() {
	start := time.Now()
	return func() {
		fmt.Printf("Total quiz time: %.2fs\n", time.Since(start).Seconds())
	}
}
