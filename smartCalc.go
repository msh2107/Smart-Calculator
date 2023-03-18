package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Stack[T any] struct {
	storage []T
}

func (s *Stack[T]) Push(value T) {
	s.storage = append(s.storage, value)
}

func (s *Stack[T]) Pop() (el T) {
	last := len(s.storage) - 1

	value := s.storage[last]     // save the value
	s.storage = s.storage[:last] // remove the last element

	return value // return saved value and nil error
}

func (s *Stack[T]) Top() (el T) {
	last := len(s.storage) - 1

	el = s.storage[last] // save the value

	return // return saved value and nil error
}

func (s *Stack[T]) IsEmpty() bool {
	if len(s.storage) <= 0 {
		return true
	}
	return false
}

const HELP = "The program calculates the sum and diff of numbers"

var variables = map[string]int{}

func main() {
	workWithInput()
	fmt.Println("Bye!")
}

func workWithInput() {
	scanner := bufio.NewScanner(os.Stdin)
	processInput(scanner)
}

func processInput(scanner *bufio.Scanner) {

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "/") {
			switch scanner.Text() {
			case "/exit":
				return
			case "/help":
				fmt.Println(HELP)
				continue
			default:
				fmt.Println("Unknown command")
				continue
			}
		}

		if scanner.Text() == "" {
			continue
		}

		if strings.Contains(scanner.Text(), "=") {
			updateVariables(scanner.Text())
			continue
		}

		if line, err := fromInfixToPostfix(scanner.Text()); err == nil {
			if res, err := calculation(line); err == nil {
				fmt.Println(res)
			}
		}
	}
}

func updateVariables(line string) {
	newLine := strings.SplitN(line, "=", 2)
	nameOfVar := strings.TrimSpace(newLine[0])
	value := strings.TrimSpace(newLine[1])

	if !isLatin(nameOfVar) {
		fmt.Println("Invalid identifier")
		return
	}

	if num, err := strconv.Atoi(value); err == nil {
		variables[nameOfVar] = num
		return
	}
	if !isLatin(value) {
		fmt.Println("Invalid assignment")
		return
	}
	if num, ok := variables[value]; ok {
		variables[nameOfVar] = num
	} else {
		fmt.Println("Unknown variable")
		return
	}

}

func isOperator(line string) bool {
	_, err := strconv.Atoi(line)
	if err != nil && !isNameOfVar(line) && !isLatin(line) && line != "(" && line != ")" {
		return true
	} else {
		return false
	}
}

func updateOperation(sign string) string {
	switch {
	case len(sign)%2 == 0 && sign[:1] == "-":
		return "+"
	case len(sign)%2 == 1 && sign[:1] == "-":
		return "-"
	case strings.Contains(sign, "+"):
		return "+"
	case sign == "*":
		return "*"
	case sign == "/":
		return "/"

	default:
		fmt.Println("Invalid expression")
		return ""
	}
}

func isLatin(line string) bool {
	for _, char := range line {
		if !unicode.In(char, unicode.Latin) {
			return false
		}
	}
	return true
}

func isNameOfVar(name string) bool {
	if !isLatin(name) {
		return false
	}
	if _, ok := variables[name]; ok {
		return true
	}
	return false
}

func fromInfixToPostfix(line string) (string, error) {
	operators := Stack[string]{}
	line = strings.ReplaceAll(line, "(", "( ")
	line = strings.ReplaceAll(line, ")", " )")
	chars := strings.Fields(line)
	result := ""
	if strings.Count(line, "(") != strings.Count(line, ")") {
		fmt.Println("Invalid expression")
		return "", errors.New("invalid expression")
	}

	for _, char := range chars {
		if char == "(" {
			operators.Push(char)
		} else if char == ")" {
			for !operators.IsEmpty() && operators.Top() != "(" {
				top := operators.Top()
				result += top + " "
				operators.Pop()
			}
			operators.Pop()
		} else if !isOperator(char) {
			result += char + " "
		} else if isOperator(char) {
			char = updateOperation(char)
			if char == "" {
				return "", errors.New("invalid expression")
			}
			for !operators.IsEmpty() && Precedence(char) <= Precedence(operators.Top()) && operators.Top() != "(" {
				top := operators.Top()
				result += top + " "
				operators.Pop()
			}
			operators.Push(char)
		}
	}

	for !operators.IsEmpty() {
		top := operators.Top()
		result += top + " "
		operators.Pop()
	}

	return result, nil
}

func Precedence(operator string) (precedence int) {
	switch operator {
	case "+":
		precedence = 1
	case "-":
		precedence = 1
	case "*":
		precedence = 2
	case "/":
		precedence = 2
	}
	return precedence
}

func calculation(line string) (int, error) {
	var stack Stack[string]
	chars := strings.Fields(line)

	for _, char := range chars {
		if num, err := strconv.Atoi(char); err == nil {
			stack.Push(strconv.Itoa(num))
		} else if isLatin(char) {
			if isNameOfVar(char) {
				stack.Push(strconv.Itoa(variables[char]))
			} else {
				fmt.Println("Unknown variable")
				return 0, errors.New("un var")
			}

		} else if isOperator(char) {
			i := stack.Pop()
			j := stack.Pop()
			second, _ := strconv.Atoi(i)
			first, _ := strconv.Atoi(j)
			stack.Push(strconv.Itoa(compute(first, second, char)))
		}
	}
	result, _ := strconv.Atoi(stack.Top())
	return result, nil
}

func compute(first int, second int, operation string) (result int) {
	updateOperation(operation)
	switch operation {
	case "+":
		result = first + second
	case "-":
		result = first - second
	case "*":
		result = first * second
	case "/":
		result = first / second
	}
	return result
}
