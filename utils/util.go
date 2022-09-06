package utils

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

// RandomNum: return random number between range min - max
func RandomNum(min, max int) int {
	max++

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// FibonacciSequence: return fibonacci sequence of input index start from 0
func FibonacciSequence(index int) int {
	if index < 2 {
		return index
	}

	return FibonacciSequence(index-1) + FibonacciSequence(index-2)
}

// CheckIsPrime: check if the given number is prime
func CheckIsPrime(num int) bool {
	if num < 2 {
		return true
	}

	sq_root := int(math.Sqrt(float64(num)))
	for i := 2; i <= sq_root; i++ {
		if num%i == 0 {
			return false
		}
	}

	return true
}

// SetSort: function to convert string param to primitive.D
func SetSort(orderBy string) map[string]interface{} {
	var sort int64
	descending := strings.HasPrefix(orderBy, "-")

	if descending {
		sort = -1
		orderBy = orderBy[1:]
	} else {
		sort = 1
	}

	return map[string]interface{}{orderBy: sort}
}
