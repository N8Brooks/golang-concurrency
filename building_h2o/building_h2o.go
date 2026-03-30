//go:build challenge

// Package buildingh2o contains the challenge version of the H2O problem.
//
// In the H2O problem, hydrogen and oxygen threads must proceed in complete
// groups of two hydrogens and one oxygen. All three threads in one molecule
// must invoke bond before any thread from the next molecule does.
package buildingh2o

type H2O struct{}

func NewH2O() *H2O {
	return &H2O{}
}

func (h *H2O) Hydrogen(bond func()) {
	panic("unimplemented")
}

func (h *H2O) Oxygen(bond func()) {
	panic("unimplemented")
}
