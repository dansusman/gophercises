package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)


type question struct {
    prompt string
    expected string
}

type quiz struct {
    answers int
    correct int
    questions []question
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

        question := question{
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

            // ensure extra whitespace, capitalization, etc. are 
            // not considered incorrect
            answer = strings.ToLower(answer)
            answer = strings.TrimSpace(answer)
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

func (quiz *quiz) shuffle() {
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(quiz.questions), func(i, j int) { 
        quiz.questions[i], quiz.questions[j] = quiz.questions[j], quiz.questions[i]
    })

}

func (quiz *quiz) assess() {
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
    shuffle = flag.Bool("shuffle", false, "Shuffle the Quiz Questions")
)

func main() {
    flag.Parse()
    quiz := readQuiz(*filePathPtr)
    if *shuffle {
        quiz.shuffle()
    }
    quiz.run()
    quiz.assess()
}

