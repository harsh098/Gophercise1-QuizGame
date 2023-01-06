package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	//"time"
)

var File *string
var TimeLimit *uint
const INFTIME uint = 1 << 64 -1
var score uint = 0
var count uint = 0
var questionSet [][]string = make([][]string, 0,1024) //For Performance Reasons
var answers = make(chan bool)
var wg = sync.WaitGroup{}


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
	
}



func readNextQuestion(idx int,q, a string){
	scanner := bufio.NewReader(os.Stdin)
	fmt.Printf("Problem #%v %v :",idx+1, q)
	input, err := scanner.ReadString('\n')
	if err != nil{
			panic("Error in Reading Input")
	}
	input = strings.TrimRight(input, "\n")
	
	if input == a{
		answers<-true
	}else{
		answers<-false
	}
	
	count++
	wg.Done()
}

func updateScore(x <-chan bool){
	if y:=<-x; y{
		score++
	}
	wg.Done()
}




func startQuiz(){
	//code for quizLogic
	for idx,x := range questionSet{
		q, a:= x[0], x[1]
		wg.Add(2)
		go readNextQuestion(idx, q, a)
		go updateScore(answers)
		wg.Wait()
	}
	close(answers)
	fmt.Println("Your Score was ", score, " out of ",count)
}

func main(){
	//Logic for Quiz Controller
	startQuiz()
	//CleanUps if any

}