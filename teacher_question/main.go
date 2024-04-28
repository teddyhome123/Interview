package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Student struct {
	Name      string
	isCorrect bool
	stdAnswer float64
}

type Question struct {
	Text   string
	Answer float64
}

func main() {
	students := []Student{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
		{Name: "E"},
	}

	questionCh := make(chan Question, 5)
	answerCh := make(chan Student, 5)

	go teacher(questionCh, answerCh)
	go studentsGroup(students, questionCh, answerCh)

	select {}
}

// teacher function models the teacher's behavior
func teacher(questionCh chan<- Question, answerCh <-chan Student) {
	for {
		wrongfreq := 0
		fmt.Println("Teacher: Guys, are you ready?")
		time.Sleep(3 * time.Second)

		question := generateQuestion()
		fmt.Printf("Teacher: %s = ?\n", question.Text)

		questionCh <- question

		for i := 0; i < 5; i++ {
			ans := <-answerCh
			if ans.isCorrect { //ok
				fmt.Printf("Student %s: %s = %.2f!\n", ans.Name, question.Text, question.Answer)
				fmt.Printf("Teacher: %s, you are right!\n", ans.Name)
				for _, name := range []string{"A", "B", "C", "D", "E"} {
					if name != ans.Name {
						fmt.Printf("Student %s: %s, you win.\n", name, ans.Name)
					}
				}
				clearChannel(answerCh)
				break
			} else { //not ok
				fmt.Printf("Student %s: %s = %.2f!\n", ans.Name, question.Text, ans.stdAnswer)
				fmt.Printf("Teacher: %s, you are wrong!\n", ans.Name)
				wrongfreq++
			}
			if wrongfreq == 5 {
				wrongfreq = 0
				fmt.Printf("Teacher: Boooo~ Answer is %.2f.\n", question.Answer)
				break
			}
		}
	}
}

// studentsGroup handles the group of students
func studentsGroup(students []Student, questionCh <-chan Question, answerCh chan<- Student) {
	var wg sync.WaitGroup

	for question := range questionCh {
		wg.Add(len(students))
		for _, student := range students {
			go student.answer(question, &wg, answerCh)
		}
		wg.Wait()
	}
}

// answer simulates a student attempting to answer a question
func (s Student) answer(question Question, wg *sync.WaitGroup, answerCh chan<- Student) {
	defer wg.Done()
	c := rand.Intn(101)

	if rand.Intn(10) < 3 {
		s.isCorrect = false
		s.stdAnswer = float64(c)
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
