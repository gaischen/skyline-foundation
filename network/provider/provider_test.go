package provider

import (
	"fmt"
	"testing"
)

func TestLabel(t *testing.T) {
LOOP:
	for i := 0; i < 10; i++ {
		if i == 5 {
			fmt.Println("loop:", i)
			break LOOP
		}
	}

	fmt.Println("out..")
}
