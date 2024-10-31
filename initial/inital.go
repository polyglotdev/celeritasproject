// Package inital contains initial code for the celeritas project
package initial

import (
	"fmt"
	"math"
	"strings"
)

// TestFunc is a test function that adds two numbers with overflow protection
func TestFunc(a, b int) int {
	if wouldOverflow(OpAdd, a, b) {
		return 0
	}
	return a + b
}

// Maths takes in an arbitrary number of integers and the math operation to perform
// and returns the result
func Maths(operation Operation, numbers ...int) int {
	switch operation {
	case OpAdd:
		return Add(numbers...)
	case OpSubtract:
		return Subtract(numbers...)
	case OpMultiply:
		return Multiply(numbers...)
	case OpDivide:
		return Divide(numbers...)
	case OpModulus:
		return Modulus(numbers...)
	}
	return 0
}

// Add takes in an arbitrary number of integers and returns the sum
func Add(numbers ...int) int {
	sum := 0
	for _, number := range numbers {
		if wouldOverflow(OpAdd, sum, number) {
			return 0 // or handle overflow as needed
		}
		sum += number
	}
	return sum
}

// Subtract takes in an arbitrary number of integers and returns the difference
func Subtract(numbers ...int) int {
	difference := 0
	for _, number := range numbers {
		if wouldOverflow(OpSubtract, difference, number) {
			return 0 // or handle overflow as needed
		}
		difference -= number
	}
	return difference
}

// Multiply takes in an arbitrary number of integers and returns the product
func Multiply(numbers ...int) int {
	product := 1
	for _, number := range numbers {
		if number == 0 {
			return 0
		}
		if wouldOverflow(OpMultiply, product, number) {
			return 0 // or handle overflow as needed
		}
		product *= number
	}
	return product
}

// Divide takes in an arbitrary number of integers and returns the quotient
func Divide(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}

	quotient := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if numbers[i] == 0 { // Check for division by zero
			return 0
		}
		// Check for MinInt / -1 overflow case
		if quotient == math.MinInt && numbers[i] == -1 {
			return 0 // or math.MaxInt, depending on your requirements
		}
		quotient /= numbers[i]
	}
	return quotient
}

// Modulus takes in an arbitrary number of integers and returns the remainder
func Modulus(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}

	remainder := numbers[0] // Start with first number
	for i := 1; i < len(numbers); i++ {
		remainder %= numbers[i] // Take modulus with subsequent numbers
	}
	return remainder
}

const (
	English = "english"
	Spanish = "spanish"
	French  = "french"
	German  = "german"
	Italian = "italian"
)

// SayHello takes in a name and a language and returns a greeting
func SayHello(name, language string) string {
	switch strings.ToLower(language) {
	case English:
		return fmt.Sprintf("Hello, %s!", name)
	case Spanish:
		return fmt.Sprintf("Hola, %s!", name)
	case French:
		return fmt.Sprintf("Bonjour, %s!", name)
	case German:
		return fmt.Sprintf("Hallo, %s!", name)
	case Italian:
		return fmt.Sprintf("Ciao, %s!", name)
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Operation represents a mathematical operation
type Operation int

const (
	OpAdd Operation = iota
	OpSubtract
	OpMultiply
	OpDivide
	OpModulus
	OpInvalid Operation = -1 // Explicit invalid operation
)

// wouldOverflow checks if performing the operation on a and b would result in overflow
func wouldOverflow(op Operation, a, b int) bool {
	switch op {
	case OpAdd:
		if a > 0 && b > 0 {
			return a > math.MaxInt-b
		}
		if a < 0 && b < 0 {
			return a < math.MinInt-b
		}

	case OpSubtract:
		if b > 0 {
			return a < math.MinInt+b
		}
		if b < 0 {
			return a > math.MaxInt+b
		}

	case OpMultiply:
		if a > 0 && b > 0 {
			return a > math.MaxInt/b
		}
		if a < 0 && b < 0 {
			return a < math.MaxInt/b
		}
		if a > 0 && b < 0 {
			return b < math.MinInt/a
		}
		if a < 0 && b > 0 {
			return a < math.MinInt/b
		}

	case OpDivide:
		return a == math.MinInt && b == -1
	}

	return false
}

func (op Operation) String() string {
	switch op {
	case OpAdd:
		return "add"
	case OpSubtract:
		return "subtract"
	case OpMultiply:
		return "multiply"
	case OpDivide:
		return "divide"
	case OpModulus:
		return "modulus"
	default:
		return "unknown"
	}
}
