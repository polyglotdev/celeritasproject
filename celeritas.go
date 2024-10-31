package celeritas

// TestFunc is a test function
func TestFunc(a, b int) int {
	return a + b
}

// Maths takes in an arbitrary number of integers and the math operation to perform
// and returns the result
func Maths(operation string, numbers ...int) int {
	// Add subtract, multiply, divide, modulus operations
	switch operation {
	case "add":
		return Add(numbers...)
	case "subtract":
		return Subtract(numbers...)
	case "multiply":
		return Multiply(numbers...)
	case "divide":
		`return Divide(numbers...)
	case "modulus":
		return Modulus(numbers...)
	}
	return 0
}

// Add takes in an arbitrary number of integers and returns the sum
func Add(numbers ...int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// Subtract takes in an arbitrary number of integers and returns the difference
func Subtract(numbers ...int) int {
	difference := 0
	for _, number := range numbers {
		difference -= number
	}
	return difference
}

// Multiply takes in an arbitrary number of integers and returns the product
func Multiply(numbers ...int) int {
	product := 1
	for _, number := range numbers {
		product *= number
	}
	return product
}

// Divide takes in an arbitrary number of integers and returns the quotient
func Divide(numbers ...int) int {
	quotient := 1
	for _, number := range numbers {
		quotient /= number
	}
	return quotient
}

// Modulus takes in an arbitrary number of integers and returns the remainder
func Modulus(numbers ...int) int {
	remainder := 0
	for _, number := range numbers {
		remainder %= number
	}
	return remainder
}
