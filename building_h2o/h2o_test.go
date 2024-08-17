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
		if err := validateInput(input); err != nil {
			t.Fatalf("invalid input: %s", input)
		}

		h2o := NewH2O()
		actual := h2o.run(input)
		expected := makeExpected(input, actual)

		t.Logf("input:    %s", input)
		t.Logf("actual:   %s", actual)
		t.Logf("expected: %s", expected)

		if expected != actual {
			t.Errorf("actual does not equal expected")
		}
	}
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

	h := strings.Count(water, "H")
	o := strings.Count(water, "O")
	if h+o != len(water) {
		errs = append(errs, errors.New("water molecule must only contain H and O atoms"))
	}
	if h != 2*n {
		errs = append(errs, errors.New("number of H atoms must be 2n"))
	}
	if o != n {
		errs = append(errs, errors.New("number of O atoms must be n"))
	}

	return errors.Join(errs...)
}

func (h2o *H2O) run(water string) string {
	var wg sync.WaitGroup
	wg.Add(len(water))
	ch := make(chan rune, len(water))

	for _, c := range water {
		switch c {
		case 'H':
			go func() {
				defer wg.Done()
				h2o.Hydrogen(func() {
					ch <- 'H'
				})
			}()
		case 'O':
			go func() {
				defer wg.Done()
				h2o.Oxygen(func() {
					ch <- 'O'
				})
			}()
		}
	}

	wg.Wait()
	close(ch)

	var actual strings.Builder
	actual.Grow(len(water))
	for i := 0; i < len(water); i++ {
		actual.WriteRune(<-ch)
	}
	return actual.String()
}

func makeExpected(input, actual string) string {
	// Map input to actual length
	if len(input) > len(actual) {
		actual += strings.Repeat(" ", len(input)-len(actual))
	} else if len(input) < len(actual) {
		actual = actual[:len(input)]
	}

	splits := strings.Split(actual, "")
	chunks := slices.Chunk(splits, 3)
	var expected strings.Builder
	for chunk := range chunks {
		molecule := strings.Join(chunk, "")
		h := strings.Count(molecule, "H")
		o := strings.Count(molecule, "O")
		if h == 2 && o == 1 {
			expected.WriteString(molecule)
		} else {
			expected.WriteString("HHO")
		}
	}
	return expected.String()
}
