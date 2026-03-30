// Package testsuite contains reusable behavioral tests for H2O implementations.
package testsuite

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
)

type H2O interface {
	Hydrogen(bond func())
	Oxygen(bond func())
}

type worker struct {
	done chan struct{}
}

func startHydrogen(h2o H2O, bond func()) worker {
	w := worker{done: make(chan struct{})}
	go func() {
		h2o.Hydrogen(bond)
		close(w.done)
	}()
	return w
}

func startOxygen(h2o H2O, bond func()) worker {
	w := worker{done: make(chan struct{})}
	go func() {
		h2o.Oxygen(bond)
		close(w.done)
	}()
	return w
}

func requireClosed(t *testing.T, ch <-chan struct{}, message string) {
	t.Helper()

	select {
	case <-ch:
	default:
		t.Fatal(message)
	}
}

func runSequence(t *testing.T, h2o H2O, sequence string) string {
	t.Helper()

	var mu sync.Mutex
	actual := make([]byte, 0, len(sequence))
	workers := make([]worker, 0, len(sequence))

	for _, atom := range sequence {
		switch atom {
		case 'H':
			workers = append(workers, startHydrogen(h2o, func() {
				mu.Lock()
				actual = append(actual, 'H')
				mu.Unlock()
			}))
		case 'O':
			workers = append(workers, startOxygen(h2o, func() {
				mu.Lock()
				actual = append(actual, 'O')
				mu.Unlock()
			}))
		default:
			t.Fatalf("unexpected atom %q", atom)
		}
		synctest.Wait()
	}

	synctest.Wait()

	for i, w := range workers {
		requireClosed(t, w.done, fmt.Sprintf("worker %d did not complete", i))
	}

	return string(actual)
}

func validateMolecules(t *testing.T, actual string) {
	t.Helper()

	if len(actual)%3 != 0 {
		t.Fatalf("got %d bond calls, want a multiple of 3", len(actual))
	}

	for i := 0; i < len(actual); i += 3 {
		molecule := actual[i : i+3]
		h := strings.Count(molecule, "H")
		o := strings.Count(molecule, "O")
		if h != 2 || o != 1 {
			t.Fatalf("molecule %q has %d hydrogens and %d oxygens, want 2 H and 1 O", molecule, h, o)
		}
	}
}

func Run(t *testing.T, newImpl func() H2O) {
	t.Helper()

	t.Run("OneShot", func(t *testing.T) {
		for _, sequence := range []string{
			"HHO",
			"HOH",
			"OHH",
			"OOHHHH",
			"HHHHOO",
			"HOHHOH",
		} {
			t.Run(sequence, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					actual := runSequence(t, newImpl(), sequence)
					validateMolecules(t, actual)
				})
			})
		}
	})

	t.Run("Reusable", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			h2o := newImpl()
			for _, sequence := range []string{"HHO", "OHH", "OOHHHH", "HOHHOH", "HHHHOO"} {
				actual := runSequence(t, h2o, sequence)
				validateMolecules(t, actual)
			}
		})
	})

	t.Run("EarlierMoleculeBondsFirst", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			h2o := newImpl()
			type call struct {
				group int
				atom  byte
			}

			var mu sync.Mutex
			calls := make([]call, 0, 6)

			launch := func(group int, atom byte) worker {
				bond := func() {
					mu.Lock()
					calls = append(calls, call{group: group, atom: atom})
					mu.Unlock()
				}
				if atom == 'H' {
					return startHydrogen(h2o, bond)
				}
				return startOxygen(h2o, bond)
			}

			workers := make([]worker, 0, 6)
			for _, atom := range []struct {
				group int
				atom  byte
			}{
				{group: 1, atom: 'H'},
				{group: 1, atom: 'H'},
				{group: 1, atom: 'O'},
				{group: 2, atom: 'O'},
				{group: 2, atom: 'H'},
				{group: 2, atom: 'H'},
			} {
				workers = append(workers, launch(atom.group, atom.atom))
				synctest.Wait()
			}
			synctest.Wait()

			for i, w := range workers {
				requireClosed(t, w.done, fmt.Sprintf("worker %d did not complete", i))
			}

			if len(calls) != 6 {
				t.Fatalf("got %d bond calls, want 6", len(calls))
			}

			for i := range 3 {
				if calls[i].group != 1 {
					t.Fatalf("call %d came from group %d, want first molecule to bond first", i, calls[i].group)
				}
			}
			for i := 3; i < 6; i++ {
				if calls[i].group != 2 {
					t.Fatalf("call %d came from group %d, want second molecule to bond second", i, calls[i].group)
				}
			}

			actual := make([]byte, 0, len(calls))
			for _, call := range calls {
				actual = append(actual, call.atom)
			}
			validateMolecules(t, string(actual))
		})
	})
}

func Benchmark(b *testing.B, newImpl func() H2O) {
	b.Helper()
	b.ReportAllocs()

	for _, molecules := range []int{1, 8} {
		b.Run(fmt.Sprintf("Molecules%d", molecules), func(b *testing.B) {
			for b.Loop() {
				h2o := newImpl()
				var wg sync.WaitGroup
				wg.Add(molecules * 3)

				for range molecules * 2 {
					go func() {
						h2o.Hydrogen(func() {})
						wg.Done()
					}()
				}

				for range molecules {
					go func() {
						h2o.Oxygen(func() {})
						wg.Done()
					}()
				}

				wg.Wait()
			}
		})
	}
}
