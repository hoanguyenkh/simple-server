package main

import (
	"testing"
)

func TestGenUniqueNumber(t *testing.T) {
	// Generate a set of unique numbers
	numSet := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		num := genUniqueNumber()
		if _, exists := numSet[num]; exists {
			t.Errorf("Duplicate number found: %s", num)
		}
		numSet[num] = true
	}
}

func BenchmarkGenUniqueNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		genUniqueNumber()
	}
}
