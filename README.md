# Elliptic Curve Toolbox

Welcome to the Elliptic Curve Toolbox, a comprehensive suite of tools designed for the analysis and manipulation of elliptic curves. This toolbox aims to provide robust utilities for academic research, cryptographic applications, and educational exploration of elliptic curve theory.

## Current Features

- **Big Integer Arithmetic**: Perform arithmetic operations on arbitrarily large integers with precise results.

## Getting Started

### Prerequisites

- Go version 1.16 or higher

### Installation

Clone the repository:

```bash
git clone https://github.com/fishfugu/elliptic.git
cd elliptic
```

#### Build the project

```
go build ./cmd/bigmath
```

#### Usage

To perform big integer arithmetic, run:

```
./bigmath
```

This will execute the arithmetic operations defined in `cmd/bigmath/main.go`

### Examples

```
// Add two large numbers
result, err := bigarith.Add("12345678901234567890", "98765432109876543210")
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Addition Result:", result)
}
```

## TODO
- [ ] Core algorithms for elliptic curve operations
- [ ] Functions for operations within finite fields
- [ ] Visualization tools for elliptic curves
- [ ] Additional command line tools for other utilities
- [ ] Comprehensive unit tests for elliptic curve functionalities
- [ ] Enhanced documentation and examples for all tools

## Contributing

Welcome contributions from anyone. To contribute, fork the repository and submit a pull request.

