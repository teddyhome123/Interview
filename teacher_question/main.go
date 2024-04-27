package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Student represents a student with a name.
type Student struct {
	Name      string
	isCorrect bool
	stdAnswer float64
}

// Question struct represents a math question.
type Question struct {
	Text   string
	Answer float64
}

func main() {
	//rand.Seed(time.Now().UnixNano())

	// Create students
	students := []Student{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
	}

	// Use a channel to communicate the question to the students
	questionCh := make(chan Question, 20)
	answerCh := make(chan Student, 20)

	// Teacher and students goroutines
	go teacher(questionCh, answerCh)
	go studentsGroup(students, questionCh, answerCh)

	// Wait for the simulation to end (which it won't in this simple example)
	select {}
}

// teacher function models the teacher's behavior
func teacher(questionCh chan<- Question, answerCh <-chan Student) {
	for {
		// Teacher starts the class
		fmt.Println("Teacher: Guys, are you ready?")
		time.Sleep(3 * time.Second) // Warm-up time

		question := generateQuestion()
		fmt.Printf("Teacher: %s = ?\n", question.Text)

		// Send the question to the students
		questionCh <- question

		// Wait for an answer from a student
		for i := 0; i < 5; i++ {
			winner := <-answerCh
			if winner.isCorrect { //ok
				fmt.Printf("Student %s: %s = %.2f!\n", winner.Name, question.Text, winner.stdAnswer)
				fmt.Printf("Teacher: %s, you are right!\n", winner.Name)
				for _, name := range []string{"A", "B", "C", "D", "E"} {
					if name != winner.Name {
						fmt.Printf("Student %s: %s, you win.\n", name, winner.Name)
					}
				}
				clearChannel(answerCh)
				break
			} else { //not ok
				fmt.Printf("Student %s: %s = %.2f!\n", winner.Name, question.Text, winner.stdAnswer)
				fmt.Printf("Teacher: %s, you are wrong!\n", winner.Name)
			}
		}
		fmt.Printf("Teacher: Boooo~ Answer is %.2f.\n", question.Answer)
	}
}

// studentsGroup handles the group of students
func studentsGroup(students []Student, questionCh <-chan Question, answerCh chan<- Student) {
	var wg sync.WaitGroup

	for question := range questionCh {
		wg.Add(len(students))
		// Each student will try to answer
		for _, student := range students {
			go student.answer(question, &wg, answerCh)
		}
		wg.Wait()
	}
}

// answer simulates a student attempting to answer a question
func (s Student) answer(question Question, wg *sync.WaitGroup, answerCh chan<- Student) {
	defer wg.Done()
	// Random thinking time between 1 and 3 seconds
	c := rand.Intn(101)

	if rand.Intn(10) < 3 {
		s.stdAnswer = float64(c)
		s.isCorrect = false
	} else {
		s.isCorrect = true
		s.stdAnswer = question.Answer
	}

	time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)

	answerCh <- s
}

// generateQuestion creates a random math question
func generateQuestion() Question {
	a := rand.Intn(101)
	b := rand.Intn(101)

	operators := []string{"+", "-", "*", "/"}
	op := operators[rand.Intn(len(operators))]

	// Formulate the question text
	questionText := fmt.Sprintf("%d %s %d", a, op, b)
	answer := evaluateExpression(a, b, op)
	return Question{Text: questionText, Answer: answer}
}

// evaluateExpression calculates the answer to a given question
func evaluateExpression(a, b int, op string) float64 {
	switch op {
	case "+":
		return float64(a + b)
	case "-":
		return float64(a - b)
	case "*":
		return float64(a * b)
	case "/":
		// Handle division by zero by returning zero
		if b == 0 {
			return 0
		}
		return float64(a) / float64(b)
	default:
		return 0
	}
}

func clearChannel(ch <-chan Student) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}
