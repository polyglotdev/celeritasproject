// Package initial contains tests for the initial package
package initial

import (
	"math"
	"testing"
)

func TestSayHello(t *testing.T) {
	type args struct {
		name     string
		language string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "english",
			args: args{
				name:     "John",
				language: English,
			},
			want: "Hello, John!",
		},
		{
			name: "spanish",
			args: args{
				name:     "Ezra",
				language: Spanish,
			},
			want: "Hola, Ezra!",
		},
		{
			name: "french",
			args: args{
				name:     "Ezra",
				language: French,
			},
			want: "Bonjour, Ezra!",
		},
		{
			name: "german",
			args: args{
				name:     "Stephanie",
				language: German,
			},
			want: "Hallo, Stephanie!",
		},
		{
			name: "italian",
			args: args{
				name:     "Ezra",
				language: Italian,
			},
			want: "Ciao, Ezra!",
		},
		{
			name: "no language passed in",
			args: args{
				name:     "Elijah",
				language: "",
			},
			want: "Hello, Elijah!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := SayHello(tt.args.name, tt.args.language); got != tt.want {
				ts.Errorf("SayHello() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModulus(t *testing.T) {
	type args struct {
		numbers []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "two numbers",
			args: args{
				numbers: []int{10, 3},
			},
			want: 1,
		},
		{
			name: "three numbers",
			args: args{
				numbers: []int{10, 3, 2},
			},
			want: 1,
		},
		{
			name: "no numbers",
			args: args{
				numbers: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := Modulus(tt.args.numbers...); got != tt.want {
				ts.Errorf("Modulus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	type args struct {
		numbers []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "division by zero",
			args: args{
				numbers: []int{10, 0},
			},
			want: 0,
		},
		{
			name: "basic division",
			args: args{
				numbers: []int{10, 2},
			},
			want: 5,
		},
		{
			name: "multiple numbers",
			args: args{
				numbers: []int{100, 2, 2, 5},
			},
			want: 5,
		},
		{
			name: "division with negative numbers",
			args: args{
				numbers: []int{-10, 2},
			},
			want: -5,
		},
		{
			name: "empty slice",
			args: args{
				numbers: []int{},
			},
			want: 0,
		},
		{
			name: "single number",
			args: args{
				numbers: []int{5},
			},
			want: 5,
		},
		{
			name: "integer overflow case",
			args: args{
				numbers: []int{math.MinInt, -1},
			},
			want: 0,
		},
		{
			name: "zero divided by number",
			args: args{
				numbers: []int{0, 5},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := Divide(tt.args.numbers...); got != tt.want {
				ts.Errorf("Divide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		numbers []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "basic addition",
			args: args{
				numbers: []int{1, 2, 3},
			},
			want: 6,
		},
		{
			name: "empty slice",
			args: args{
				numbers: []int{},
			},
			want: 0,
		},
		{
			name: "single number",
			args: args{
				numbers: []int{5},
			},
			want: 5,
		},
		{
			name: "negative numbers",
			args: args{
				numbers: []int{-1, -2, -3},
			},
			want: -6,
		},
		{
			name: "integer overflow positive",
			args: args{
				numbers: []int{math.MaxInt, 1},
			},
			want: 0,
		},
		{
			name: "integer overflow negative",
			args: args{
				numbers: []int{math.MinInt, -1},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := Add(tt.args.numbers...); got != tt.want {
				ts.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaths(t *testing.T) {
	tests := []struct {
		name      string
		operation Operation
		numbers   []int
		want      int
	}{
		{
			name:      "add operation",
			operation: OpAdd,
			numbers:   []int{1, 2, 3},
			want:      6,
		},
		{
			name:      "subtract operation",
			operation: OpSubtract,
			numbers:   []int{10, 3, 2},
			want:      -15,
		},
		{
			name:      "multiply operation",
			operation: OpMultiply,
			numbers:   []int{2, 3, 4},
			want:      24,
		},
		{
			name:      "divide operation",
			operation: OpDivide,
			numbers:   []int{12, 3, 2},
			want:      2,
		},
		{
			name:      "modulus operation",
			operation: OpModulus,
			numbers:   []int{10, 3},
			want:      1,
		},
		{
			name:      "invalid operation returns 0",
			operation: Operation(99), // Some arbitrary invalid operation
			numbers:   []int{1, 2, 3},
			want:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := Maths(tt.operation, tt.numbers...); got != tt.want {
				ts.Errorf("Maths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestFunc(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "basic addition",
			args: args{
				a: 1,
				b: 2,
			},
			want: 3,
		},
		{
			name: "negative numbers",
			args: args{
				a: -1,
				b: -2,
			},
			want: -3,
		},
		{
			name: "integer overflow",
			args: args{
				a: math.MaxInt,
				b: 1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := TestFunc(tt.args.a, tt.args.b); got != tt.want {
				ts.Errorf("TestFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_String(t *testing.T) {
	tests := []struct {
		name string
		op   Operation
		want string
	}{
		{
			name: "add operation",
			op:   OpAdd,
			want: "add",
		},
		{
			name: "invalid operation",
			op:   Operation(99),
			want: "unknown",
		},
		{
			name: "divide operation",
			op:   OpDivide,
			want: "divide",
		},
		{
			name: "modulus operation",
			op:   OpModulus,
			want: "modulus",
		},
		{
			name: "subtract operation",
			op:   OpSubtract,
			want: "subtract",
		},
		{
			name: "multiply operation",
			op:   OpMultiply,
			want: "multiply",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := tt.op.String(); got != tt.want {
				ts.Errorf("Operation.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wouldOverflow(t *testing.T) {
	type args struct {
		op Operation
		a  int
		b  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "add operation",
			args: args{
				op: OpAdd,
				a:  1,
				b:  2,
			},
			want: false,
		},
		{
			name: "integer overflow",
			args: args{
				op: OpAdd,
				a:  math.MaxInt,
				b:  1,
			},
			want: true,
		},
		{
			name: "integer overflow negative",
			args: args{
				op: OpAdd,
				a:  math.MinInt,
				b:  -1,
			},
			want: true,
		},
		{
			name: "divide operation",
			args: args{
				op: OpDivide,
				a:  math.MinInt,
				b:  -1,
			},
			want: true,
		},
		{
			name: "multiply operation",
			args: args{
				op: OpMultiply,
				a:  math.MaxInt,
				b:  2,
			},
			want: true,
		},
		{
			name: "multiply operation negative",
			args: args{
				op: OpMultiply,
				a:  math.MinInt,
				b:  -1,
			},
			want: true,
		},
		{
			name: "subtract operation",
			args: args{
				op: OpSubtract,
				a:  math.MinInt,
				b:  1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			if got := wouldOverflow(tt.args.op, tt.args.a, tt.args.b); got != tt.want {
				ts.Errorf("wouldOverflow() = %v, want %v", got, tt.want)
			}
		})
	}
}
