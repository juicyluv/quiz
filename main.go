package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {
	rightAnswers := 0

	// set flag to input filename with quiz questions
	csvFileName := flag.String("csv", "questions.csv", "a csv file in (question,answer) format")
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

	for i, q := range questions {
		fmt.Printf("Question #%d: %s = ", i+1, q.question)
		answer := ""
		fmt.Scanf("%s\n", &answer)
		if answer == q.answer {
			rightAnswers += 1
			fmt.Printf("Right!\n\n")
		} else {
			fmt.Printf("Wrong.\n\n")
		}
	}

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
			answer:   line[1],
		}
	}

	return quesions
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
