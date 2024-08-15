package h2o

import (
	"errors"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestH2O(t *testing.T) {
	inputs := []string{
		"HOH",
		"OOHHHH",
		"HHHHHHHHHHOHHOHHHHOOHHHOOOOHHOOHOHHHHHOOHOHHHOOOOOOHHHHHHHHH",
	}

	for _, input := range inputs {
		h2o := NewH2O()
		actual := h2o.run(input)
		expected, err := makeExpected(input, actual)
		if err != nil {
			t.Errorf("invalid input: %s", input)
		}

		t.Logf("input: %s", input)
		t.Logf("actual: %s", actual)
		t.Logf("expected: %s", expected)

		if expected != actual {
			t.Errorf("actual does not equal expected")
		}
	}
}

func (h2o *H2O) run(water string) string {
	var wg sync.WaitGroup
	wg.Add(len(water))

	ch := make(chan string, len(water))

	releaseHydrogen := func() {
		ch <- "H"
	}

	releaseOxygen := func() {
		ch <- "O"
	}

	for _, c := range water {
		switch c {
		case 'H':
			go func() {
				defer wg.Done()
				h2o.Hydrogen(releaseHydrogen)
			}()
		case 'O':
			go func() {
				defer wg.Done()
				h2o.Oxygen(releaseOxygen)
			}()
		}
	}

	wg.Wait()
	result := ""
	for i := 0; i < len(water); i++ {
		result += <-ch
	}
	close(ch)

	return result
}

func makeExpected(input, actual string) (string, error) {
	if err := validateInput(input); err != nil {
		return "", err
	}

	// Map input to actual length
	if len(input) > len(actual) {
		actual += strings.Repeat(" ", len(input)-len(actual))
	} else if len(input) < len(actual) {
		actual = actual[:len(input)]
	}

	splits := strings.Split(actual, "")
	chunks := slices.Chunk(splits, 3)
	expected := ""
	for chunk := range chunks {
		expected += minDiffWater(chunk)
	}

	return expected, nil
}

func validateInput(water string) error {
	errs := []error{}
	if len(water)%3 != 0 {
		errs = append(errs, errors.New("water molecule length must be a multiple of 3"))
	}
	n := len(water) / 3
	if n <= 0 {
		errs = append(errs, errors.New("n must be greater than 0"))
	}
	if n > 20 {
		errs = append(errs, errors.New("n must be less than or equal to 20"))
	}
	if strings.Count(water, "O")+strings.Count(water, "H") != len(water) {
		errs = append(errs, errors.New("water molecule must only contain O and H atoms"))
	}
	if strings.Count(water, "H") != 2*n {
		errs = append(errs, errors.New("number of H atoms must be 2n"))
	}
	if strings.Count(water, "O") != n {
		errs = append(errs, errors.New("number of O atoms must be n"))
	}
	return errors.Join(errs...)
}

// Returns the water molecule with the least difference in order to the actual
func minDiffWater(chunk []string) string {
	if len(chunk) != 3 {
		panic("chunk must be of length 3")
	}

	result := ""
	var hCount, oCount int
	for _, c := range chunk {
		switch c {
		case "H":
			if hCount < 2 {
				result += "H"
				hCount++
			} else {
				result += "O"
				oCount++
			}
		case "O":
			if oCount < 1 {
				result += "O"
				oCount++
			} else {
				result += "H"
				hCount++
			}
		default:
			if hCount < 2 {
				result += "H"
				hCount++
			} else {
				result += "O"
				oCount++
			}
		}
	}
	return result
}
