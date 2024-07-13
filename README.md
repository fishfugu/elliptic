# Elliptic Curve Toolbox

Tools for analysis / manipulation of elliptic curves. Provides utilities for experimentation, and educational exploration of elliptic curve theory.

All partly inspired by my work over here: [Recent course output - SCI395, with UNE](https://www.creativearts.com.au/maths/une/sci395). I'll be updating this site and this repo in tandem, as I go, and as I get to different parts of the investigation / analysis.

## Current Features

- **Big Integer Arithmetic**: Perform arithmetic operations on arbitrarily large integers with precise results.

- **Finite Field Calculation**: Prime number finite field operations. Calculating points on an elliptic curve (in the Weierstrass form) defined on a prime integer finite field. Visualise prime integer finite field points, on an elliptic curve.


## Getting Started

### Prerequisites

- Go version 1.21 or higher

May need this:
```
xcode-select --install
```

### Installation

Clone the repository:

```bash
git clone https://github.com/fishfugu/elliptic.git
cd elliptic
```

#### Build the project
```
make build-all
```

#### Run tests
```
make test
```

#### Take it all for a test drive
AKA see example output and kick the tyres:
```
make test-drive
```

### Makefile
```
Usage: make <TARGETS>

BUILD
  build-all                       Build the project - all necessary components
  build-bigmath                   Build bigmath executable
  build-finitefield               Build finitefield executable
  build-ecviz                     Build Elliptic Curve Data Viz Tool

RUN
  run-bigmath                     Run bigmath binary after building it - ensuring latest build is executed - running tests first
  run-finitefield                 Run finitefield binary after building it - ensuring latest build is executed - running tests first
  run-ecviz                       Run Elliptic Curve Data Viz Tool (after doing a build)

TEST and CLEAN
  test                            Run unit tests for all packages under pkg
  clean                           Remove binaries and any temporary files
  testdrive                       Run through all (appropriate) make file commands - just to take it for a test drive (check I haven't done stupidity)

HELPER
  test-quiet                      Run unit tests for all packages under pkg - but quietly - quits at first error
  test-verbose                    Run unit tests for all packages under pkg - in verbose mode
  help                            Show this help
```

## Text Visualisation Example
```
2D Plane Visualization with Cartesian Axes, Reflection Line, and Scale:
|                                                                                        43
|                                                   *                                   
|                                             *                                         
|         *                                                                             
|           * *                                             *                            39
|                                                                                       
|     *                                                                                 
| *                                                                                     
|                                                                                        35
|                                                       *                               
|                                                             *                         
|                       *                                                               
|   *           *                                                 *                      31
|                                                                           *           
|                                                                                       
|                                                                                       
|                                                                                        27
|                                           *                                           
|                     *     *       *                                                   
|                                                                             *         
|                                                                               *        23
| . . . . . . . . . . . . . . * . . . . . . . . . . . . . . . . . . . . . . . . . . . .  43/2
|                             *                                                         
|                                                                               *       
|                                                                             *          19
|                     *     *       *                                                   
|                                           *                                           
|                                                                                       
|                                                                                        15
|                                                                                       
|                                                                           *           
|   *           *                                                 *                     
|                       *                                                                11
|                                                             *                         
|                                                       *                               
|                                                                                       
| *                                                                                      7
|     *                                                                                 
|                                                                                       
|           * *                                             *                           
|         *                                                                              3
|                                             *                                         
|                                                   *                                   
+ - - - / - - - / - - - / - - - / - - - / - - - / - - - / - - - / - - - / - - - / - * -  0
```

## TODO
- [x] Makefile tidy - organise / make headings for help
- [ ] Core algos for elliptic curve operations
    - [x] Arbitrarily large finite filed ops - use BigInt - but make simpler
    - [x] Find Prime
        - [ ] Work out how many checks is "good enough"
    - [ ] Check that use of BigInt versus own reimplement is consistent / sane
    - [ ] Investigate what "multiplying points" and "multiplicative inverse of points" means (extend lit review work: [SCI395 Course Material](https://www.creativearts.com.au/maths/une/sci395)
)
- [ ] Functions for operations within finite fields
    - [x] Calculate points on an EC in prime finite field
    - [x] Make ECs into types to hand around (immutable)
    - [ ] Implement adding points and multplying points in field by real nums
- [ ] Visualization tools for elliptic curves
    - [x] Basic text vis output
    - [x] Debug the text vis output - scale / reflection line is wrong
    - [ ] Improve tick marks / scale numbers under x-axis
    - [ ] Vis adding and multplying numbers (see implementation above)
    - [ ] Turn vis into a window with real drawings thingy
    - [ ] Annimate doubling / adding / rings of points - text and / or proper drawing
- [ ] Additional command line tools for other utilities
    - [ ] Turn this interactive and make current code into `testdrive` option
- [ ] Comprehensive unit tests for elliptic curve functionalities
    - [ ] Do "random" test of more complex finite field calcs - rondomly selected in loop,small numbers without bigint, double checked / compared to function output

## Contributing
Welcome contributions from anyone. Fork repo / submit PR.

Open to questions and input on where to go next - as I'd like to make this whole thing useful to others too.
