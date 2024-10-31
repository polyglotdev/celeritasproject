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
		return add(numbers...)
	case "subtract":
		return subtract(numbers...)
	case "multiply":
		return multiply(numbers...)
	case "divide":
		return divide(numbers...)
	case "modulus":
		return modulus(numbers...)
	}
	return 0
}

// add takes in an arbitrary number of integers and returns the sum
func add(numbers ...int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// subtract takes in an arbitrary number of integers and returns the difference
func subtract(numbers ...int) int {
	difference := 0
	for _, number := range numbers {
		difference -= number
	}
	return difference
}

// multiply takes in an arbitrary number of integers and returns the product
func multiply(numbers ...int) int {
	product := 1
	for _, number := range numbers {
		product *= number
	}
	return product
}

// divide takes in an arbitrary number of integers and returns the quotient
func divide(numbers ...int) int {
	quotient := 1
	for _, number := range numbers {
		quotient /= number
	}
	return quotient
}

// modulus takes in an arbitrary number of integers and returns the remainder
func modulus(numbers ...int) int {
	remainder := 0
	for _, number := range numbers {
		remainder %= number
	}
	return remainder
}
