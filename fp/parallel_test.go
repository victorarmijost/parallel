package fp

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	const size = 1000

	input := []int{}

	for i := 0; i < size; i++ {
		input = append(input, i)
	}

	out, errs := Parallel(7, input, func(i int) (string, error) {
		if i%2 == 1 {
			return "", fmt.Errorf("error on value %d", i)
		}

		return fmt.Sprintf("%d", i), nil
	})

	if len(errs) != size/2 {
		t.Error("wrong number of errors returned")
	}

	if len(out) != size/2 {
		t.Error("more results were expected")
	}

	fmt.Println(out)

}
