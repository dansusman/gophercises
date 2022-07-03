package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)


type Question struct {
    prompt string
    expected string
}

type quiz struct {
    answers int
    correct int
    total int
    questions []Question
}


func readQuiz(filePath string) *quiz {

    var quiz quiz

    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }

    defer file.Close()

    reader := csv.NewReader(file)

    for {
        record, err := reader.Read()
        // stop at EOF
        if err == io.EOF {
            break
        }

        if err != nil {
            panic(err)
        }

        question := Question{
            prompt: record[0],
            expected: record[1],
        }
        quiz.questions = append(quiz.questions, question)
        quiz.total = len(quiz.questions)
    }
    return &quiz
}

func (quiz *quiz) run() {
    for _, question := range quiz.questions {
        // print prompt
        fmt.Print(question.prompt + " ")

        // accept user input
        var answer string
        fmt.Scanln(&answer)

        // convert string -> int
        strconv.Atoi(answer)

        // compare to question.expected, increment correct if needed
        if answer == question.expected {
            quiz.correct += 1
        }
        if answer != "" {
            quiz.answers += 1
        }
    }
}

func (quiz *quiz) access() {
    fmt.Printf(
		"You answered %v questions out of a total of %v and got %v correct",
		quiz.answers,
        quiz.total,
		quiz.correct,
	)
}

func main() {
    filePath := "./problems.csv"

    if len(os.Args) > 1 {
        filePath = os.Args[1]
    }

    quiz := readQuiz(filePath)

    quiz.run()
    quiz.access()
}

