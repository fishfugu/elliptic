package main

import (
	"fmt"
	"os"

	pns "elliptic/pkg/pnumsys"
)

// to run:
// make build-pns
// ./bin/pnumsys

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error percolated to main, error: %v\n", err)
	}
}

func run() error {
	if err := pns.GenerateMaxExponentLaTeXTable("doc/latex/max_exponents_table.tex"); err != nil {
		return fmt.Errorf("error generating Max Exponent LaTeX table: %v\n", err)
	}
	fmt.Println("LaTeX table created in 'doc/latex/prime_exponents_table.tex'")

	if err := pns.GenerateOverflowLaTeXTable("doc/latex/overflow_table.tex"); err != nil {
		return fmt.Errorf("error generating Overflow LaTeX table: %v\n", err)
	}
	fmt.Println("LaTeX table created in 'doc/latex/overflow_table.tex'")
	return nil
}
