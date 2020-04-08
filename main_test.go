// This file contains the unit test cases for the greeting service.
// The unit test cases only contains pure functions w/o external dependencies
// We are not going to mock any external dependency here
// With container technology, we just embrace the dependencies in functional testing!

package main

import (
	"testing"
)

var testCases = []struct {
	x        int
	y        int
	expected int
}{
	{1, 2, 3},
	{3, 4, 7},
}

func TestGreetingService(t *testing.T) {
	for _, c := range testCases {
		if r := Add(c.x, c.y); r != c.expected {
			t.Errorf("Test Failed: input: %d, %d; expected: %d; received: %d", c.x, c.y, c.expected, r)
		}
	}
}
