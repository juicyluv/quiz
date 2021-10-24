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
	// set flag to input filename with quiz questions
	csvFileName := flag.String("csv", "questions.csv", "a csv file in (question,answer) format")
	timeLimit := flag.Int("limit", 5, "a time limit for every question in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open file \"%s\"\n", *csvFileName))
	}

	// read all the lines from .csv file
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to read provided file\n")
	}
	questions := parseLines(lines)

	rightAnswers := startQuiz(questions, timeLimit)

	fmt.Printf("Result: %d/%d.\n", rightAnswers, len(questions))
}

type question struct {
	question string
	answer   string
}

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

	for i, q := range questions {
		fmt.Printf("Question #%d: %s = ", i+1, q.question)
		timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
		answerChannel := make(chan string)

		go func() {
			answer := ""
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou run out of time.\n")
			return rightAnswers
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
