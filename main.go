package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)


type Question struct {
    prompt string
    expected string
}

type quiz struct {
    answers int
    correct int
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
    }
    return &quiz
}

func (quiz *quiz) run() {
    timer := time.NewTimer(time.Duration(*duration) * time.Second)
quizLoop:
    for _, question := range quiz.questions {
        // print prompt
        fmt.Print(question.prompt + " ")

        answerChan := make(chan string)

        go func() {
            // accept user input
            var answer string
            fmt.Scanln(&answer)
            answerChan <- answer
        }()

        select {
        case <-timer.C:
            fmt.Println()
            break quizLoop
        case answer := <-answerChan:
            // compare to question.expected, increment correct if needed
            if answer == question.expected {
                quiz.correct += 1
            }
            if answer != "" {
                quiz.answers += 1
            }
        }
    }
}

func (quiz *quiz) access() {
    fmt.Printf(
		"You answered %v questions out of a total of %v and got %v correct\n",
		quiz.answers,
        len(quiz.questions),
		quiz.correct,
	)
}

var (
    filePathPtr = flag.String("csv", "./problems.csv", "File Path for CSV")
    duration = flag.Int("time", 30, "Time Limit for Quiz")
)

func main() {
    flag.Parse()
    quiz := readQuiz(*filePathPtr)
    quiz.run()
    quiz.access()
}

