package bigarith

import (
	"fmt"
	"math/big"
)

// Int is a bigarith type for Integer Arithmetic
type Int struct {
	stringValue string // the string representing the integer value
}

func NewInt(a string) Int {
	return new(Int).set(a)
}

// Discovery functions
// the functions that surface info about the Int - don't return an Int itself

// Val returns the string representation of the bigarith.Int value
func (i Int) Val() string {
	return i.stringValue
}

// Compare takes a string representation of an integer and compares
// it with the current bigarith.Int,
// returns:
// -1 if x <  y, 0 if x == y, +1 if x >  y
func (i Int) Compare(a string) int {
	return bigInt(i.stringValue).Cmp(bigInt(a))
}

// ProbablyPrime reports whether the current integer is probably prime,
// applying the Miller-Rabin test with n pseudorandomly chosen bases,
// as well as a Baillie-PSW test
func (i Int) ProbablyPrime() bool {
	return bigInt(i.stringValue).ProbablyPrime(256)
}

// Manipulation functions
// functions that manipulate the value of an Int
// should return the new value without modifying the original Int

// set sets the string value that represents the bigarith.Int value
// and returns a new Int with the updated value
func (i Int) set(a string) Int {
	i.stringValue = a
	return i
}

// SetBigInt sets the string value that represents the bigarith.Int value
// by getting a string representation from a big.Int object
// and returns a new Int with the updated value
func (i Int) SetBigInt(a big.Int) Int {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from error: %s", r)
		}
	}()
	return i.set(a.String())
}

// Float returns the value of the current integer as a new bigarith.Float
func (i Int) Float() Float {
	return NewFloat(i.Val())
}

// Neg returns the negated value of the current integer as a new bigarith.Int
func (i Int) Neg() Int {
	return i.SetBigInt(*new(big.Int).Neg(bigInt(i.stringValue)))
}

// Plus adds a string representation of an integer to the current integer,
// and returns the result as a new bigarith.Int
func (i Int) Plus(a string) Int {
	return i.SetBigInt(*new(big.Int).Add(bigInt(i.stringValue), bigInt(a)))
}

// Minus subtracts a string representation of an integer from the current integer,
// and returns the result as a new bigarith.Int
func (i Int) Minus(a string) Int {
	return i.SetBigInt(*new(big.Int).Sub(bigInt(i.stringValue), bigInt(a)))
}

// Times multiplies the current integer with a string representation of another integer,
// and returns the result as a new bigarith.Int
func (i Int) Times(a string) Int {
	return i.SetBigInt(*new(big.Int).Mul(bigInt(i.stringValue), bigInt(a)))
}

// For "DividedBy" see "Other Functions" below

// ModularInverse calculates the modular inverse of the current integer mod p
// using the Extended Euclidean Algorithm.
// Returns a new bigarith.Int if an inverse exists, or panics if it doesn't.
func (i Int) ModularInverse(p string) Int {
	// Convert string to big.Int
	a := bigInt(i.Val())
	mod := bigInt(p)

	gcd := new(big.Int).GCD(nil, nil, a, mod)
	if gcd.Cmp(big.NewInt(1)) != 0 {
		panic(fmt.Sprintf("No modular inverse exists: GCD of %s and %s is not 1", a.String(), mod.String()))
	}

	// Variables for the extended Euclidean algorithm
	t := new(big.Int)
	newT := big.NewInt(1)
	r := new(big.Int).Set(mod)
	newR := new(big.Int).Set(a)

	for newR.Cmp(big.NewInt(0)) != 0 {
		quotient := new(big.Int).Div(r, newR)

		// Update r and newR
		r, newR = newR, new(big.Int).Sub(r, new(big.Int).Mul(quotient, newR))

		// Update t and newT
		t, newT = newT, new(big.Int).Sub(t, new(big.Int).Mul(quotient, newT))
	}

	// If r != 1, then there is no modular inverse
	if r.Cmp(big.NewInt(1)) != 0 {
		panic(fmt.Sprintf("ModularInverse returned nil. a: '%s', p: '%s'", i.stringValue, p))
	}

	// Ensure the result is positive by adding mod if necessary
	if t.Cmp(big.NewInt(0)) < 0 {
		t.Add(t, mod)
	}

	return i.SetBigInt(*t)
}

// Mod calculates the current integer mod p,
// and returns the result as a new bigarith.Int
func (i Int) Mod(p string) Int {
	return i.SetBigInt(*new(big.Int).Mod(bigInt(i.stringValue), bigInt(p)))
}

// GCD calculates the GCD of the current integer and p,
// and returns the result as a new bigarith.Int
func (i Int) GCD(p string) Int {
	return i.SetBigInt(*new(big.Int).GCD(nil, nil, bigInt(i.stringValue), bigInt(p)))
}

// DivideInField divides the current integer by a string representation of an integer d (mod p),
// by finding the ModularInverse of d mod p, and multiplying
// and returns the result as a new bigarith.Int
func (i Int) DivideInField(d, p string) Int {
	// p has to be greater than i and d
	if NewInt(p).Compare(d) < 0 || NewInt(p).Compare(i.Val()) < 0 {
		panic(fmt.Sprintf("p should be greater than both the Int and d, the things it's being divided by: i Int = %s, d = %s, p = %s", i.Val(), d, p))
	}
	return i.Times(NewInt(d).ModularInverse(p).Val()).Mod(p)
}

// ToThePowerOf raises the current integer to the power of a mod m,
// and returns the result as a new bigarith.Int
func (i Int) ToThePowerOf(a, m string) Int {
	var mInt *big.Int
	if m != "" {
		mInt = bigInt(m)
	}
	return i.SetBigInt(*new(big.Int).Exp(bigInt(i.stringValue), bigInt(a), mInt))
}

// Factorial calculates the factorial of the current bigarith.Int,
// and returns the result as a new bigarith.Int
func (i Int) Factorial() Int {
	result := NewInt("1")
	for i.Compare("1") > 0 { // for (i Int) > 1
		result = result.Times(i.Val()) // Update result with each multiplication
		i = i.Minus("1")               // Use the returned value for further calculations
	}
	return result
}

// Other Functions
// functions that don't behave in the standard way and need more consideration
// DividedBy divides the current integer by a string representation of an integer,
// and returns the result as a bigarith.Float
func (i Int) DividedBy(a string) Float {
	return NewFloat(i.Val()).DividedBy(a)
}

// FindPrime finds a prime number within the range specified by the strings i and high.
// Returns the first prime found as a string, or an error if no prime is found or inputs are invalid.
func (i Int) FindPrime(high string) Int {
	// TODO: is there any reason to implement different way / direction of search?
	// Start searching for a prime at the low end of the range
	for p := i; p.Compare(high) < 0; p = p.Plus("1") {
		if p.ProbablyPrime() {
			return p
		}
	}
	return NewInt("0")
}
