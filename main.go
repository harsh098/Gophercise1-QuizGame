package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var File *string
var TimeLimit *uint
const INFTIME uint = 1 << 64 -1
var score int = 0
var count int = 0
var questionSet [][]string = make([][]string, 0,1024) //For Performance Reasons
var answers = make(chan string)



func init(){
	//Initialise arguments
	File = flag.String("file","problems.csv","Filename of CSV file to use")
	TimeLimit = flag.Uint("time", INFTIME, "Time for Quiz")
	flag.Parse()
	fileObject , err := os.Open(*File)
	
	if err != nil{
		defer func(){
			if r:=recover(); r != nil{
				fmt.Println("Set to Defaults")
				*File = "problems.csv"
				*TimeLimit = INFTIME
				startQuiz()
			} 
		}()
	}
	defer fileObject.Close()
	
	//Create CSV Reader
	reader := csv.NewReader(fileObject)
	//Churn Questions
	questions, err := reader.ReadAll()

	if err != nil{
		panic("Malformed CSV File")
	}

	questionSet = append(questionSet, questions...)
	count = len(questionSet)	
}

func inputEvent(idx int, q, a string){
	//Take Answer as Input
	scanner := bufio.NewReader(os.Stdin)
	fmt.Printf("Problem #%v %v :",idx+1, q)
	input, err := scanner.ReadString('\n')
	if err != nil{
			panic("Error in Reading Input")
	}
	input = strings.TrimRight(input, "\n")

	answers<-input
	
}


func checkAnswer(in, a string){
	//Validate Answer
	if in == a{
		score++
	}
}





func startQuiz(){
	timer := time.NewTimer(time.Duration(*TimeLimit)*time.Second)
	//code for quizLogic
	
		for idx,x := range questionSet{
			q, a:= x[0], x[1]

			// readNextQuestion(idx, q, a)
			// updateScore(answers)
			
			go inputEvent(idx, q, a)
			select{
			case answer := <-answers:
				checkAnswer(answer, a)
			case <-timer.C:
				fmt.Println("\nTimeOut")
				return		
			}
		}
	
	
}

func cleanup(){
	close(answers)
}

func main(){
	//Logic for Quiz Controller
	startQuiz()
	//CleanUps if any
	cleanup()
	fmt.Println("Your Score was ", score, " out of ",count)
}