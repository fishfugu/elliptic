package pnumsys

import (
	"fmt"
	"math"
	"os"
	"strconv"

	ba "elliptic/pkg/bigarith"
)

// In Go, the maximum value for an int64 type is 2^63 - 1,
// which equals 9,223,372,036,854,775,807.
// This is because int64 is a 64-bit signed integer,
// where one bit is used for the sign (positive or negative),
// leaving 63 bits for the value.
// This is a 19-digit number - so our 172 digit numbers will need be represented using bit.Int
// Or Int from the bigarith package

// PrimeMultiplicativeNumber is a prime multiplicative system number representation
type primeMultiplicativeNumber struct {
	intList []int   // the int list representing the prime multiplicative system number - only up to 9
	intVal  *ba.Int // an integer which when printed in base 10 represents the prime multiplicative system number
}

func NewPrimeMultiplicativeNumber(ints []int) primeMultiplicativeNumber {
	tempIntVal := ba.NewInt("0")
	for pos, thisInt := range ints {
		if thisInt > 0 {
			// tempIntVal += thisInt * int64(math.Pow(10, 170-float64(pos)))
			numeralPosString := strconv.Itoa(171 - pos)
			multiplier := ba.NewInt("10").ToThePowerOf(numeralPosString, "")  // 10^{170-pos}
			adder := ba.NewInt(strconv.Itoa(thisInt)).Times(multiplier.Val()) // thisInt * 10^{170-pos}
			tempIntVal = tempIntVal.Plus(adder.Val())                         // tempIntVal += thisInt * 10^{170-pos}
		}
	}
	return primeMultiplicativeNumber{
		intList: ints,
		intVal:  &tempIntVal,
	}
}

// List of primes less than or equal to 1021
var primes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 607, 613, 617, 619, 631, 641, 643, 647, 653, 659, 661, 673, 677, 683, 691, 701, 709, 719, 727, 733, 739, 743, 751, 757, 761, 769, 773, 787, 797, 809, 811, 821, 823, 827, 829, 839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929, 937, 941, 947, 953, 967, 971, 977, 983, 991, 997, 1009, 1013, 1019, 1021}

// primeIndexMap maps each prime number up to 1021 to its position in the list of primes up to 1021, adjusted to start from 1.
var primeIndexMap = map[int]int{
	2: 1, 3: 2, 5: 3, 7: 4, 11: 5, 13: 6, 17: 7, 19: 8, 23: 9, 29: 10, 31: 11, 37: 12,
	41: 13, 43: 14, 47: 15, 53: 16, 59: 17, 61: 18, 67: 19, 71: 20, 73: 21, 79: 22,
	83: 23, 89: 24, 97: 25, 101: 26, 103: 27, 107: 28, 109: 29, 113: 30, 127: 31,
	131: 32, 137: 33, 139: 34, 149: 35, 151: 36, 157: 37, 163: 38, 167: 39, 173: 40,
	179: 41, 181: 42, 191: 43, 193: 44, 197: 45, 199: 46, 211: 47, 223: 48, 227: 49,
	229: 50, 233: 51, 239: 52, 241: 53, 251: 54, 257: 55, 263: 56, 269: 57, 271: 58,
	277: 59, 281: 60, 283: 61, 293: 62, 307: 63, 311: 64, 313: 65, 317: 66, 331: 67,
	337: 68, 347: 69, 349: 70, 353: 71, 359: 72, 367: 73, 373: 74, 379: 75, 383: 76,
	389: 77, 397: 78, 401: 79, 409: 80, 419: 81, 421: 82, 431: 83, 433: 84, 439: 85,
	443: 86, 449: 87, 457: 88, 461: 89, 463: 90, 467: 91, 479: 92, 487: 93, 491: 94,
	499: 95, 503: 96, 509: 97, 521: 98, 523: 99, 541: 100, 547: 101, 557: 102, 563: 103,
	569: 104, 571: 105, 577: 106, 587: 107, 593: 108, 599: 109, 601: 110, 607: 111,
	613: 112, 617: 113, 619: 114, 631: 115, 641: 116, 643: 117, 647: 118, 653: 119,
	659: 120, 661: 121, 673: 122, 677: 123, 683: 124, 691: 125, 701: 126, 709: 127,
	719: 128, 727: 129, 733: 130, 739: 131, 743: 132, 751: 133, 757: 134, 761: 135,
	769: 136, 773: 137, 787: 138, 797: 139, 809: 140, 811: 141, 821: 142, 823: 143,
	827: 144, 829: 145, 839: 146, 853: 147, 857: 148, 859: 149, 863: 150, 877: 151,
	881: 152, 883: 153, 887: 154, 907: 155, 911: 156, 919: 157, 929: 158, 937: 159,
	941: 160, 947: 161, 953: 162, 967: 163, 971: 164, 977: 165, 983: 166, 991: 167,
	997: 168, 1009: 169, 1013: 170, 1019: 171, 1021: 172,
}

// Function to generate the LaTeX table
func GenerateMaxExponentLaTeXTable(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the beginning of the table
	fmt.Fprintln(file, "% paste this into your latex file")
	fmt.Fprintf(file, "\\begin{longtable}{|c|c|c|}\n")
	fmt.Fprintf(file, "\\caption{Primes and Their Maximum Exponents Under 1021} \\label{tab:primes_exponents} \\\\\n")
	fmt.Fprintf(file, "\\hline\n")
	fmt.Fprintf(file, "\\textbf{Prime (p)} & \\textbf{Max Exponent \\((k_p)\\)} & \\textbf{Max Value} \\\\\n")
	fmt.Fprintf(file, "\\hline\n")
	fmt.Fprintf(file, "\\endfirsthead\n")
	fmt.Fprintf(file, "\\hline\n")
	fmt.Fprintf(file, "\\textbf{Prime (p)} & \\textbf{Max Exponent \\((k_p)\\)} & \\textbf{Max Value} \\\\\n")
	fmt.Fprintf(file, "\\hline\n")
	fmt.Fprintf(file, "\\endhead\n")
	fmt.Fprintf(file, "\\hline\n")
	fmt.Fprintf(file, "\\endfoot\n")

	// Calculate max exponents for each prime and write rows
	for _, prime := range primes {
		kp := maxExponent(prime, 1021)
		fmt.Fprintf(file, "%d & %d & %.0f \\\\\n", prime, kp, math.Pow(float64(prime), float64(kp)))
		if kp == 1 { // if you hit 1, stop, you can't get any lower
			break
		}
	}
	// fill in last lines
	fmt.Fprintln(file, "... & 1 & ... \\\\")
	fmt.Fprintf(file, "%d & %d & %d \\\\\n", 1021, 1, 1021)

	fmt.Fprintf(file, "\\hline\n")
	fmt.Fprintf(file, "\\end{longtable}\n")
	return nil
}

// Factorise returns the prime factors of a number.
func Factorise(n int) []int {
	factors := []int{}
	d := 2
	for n > 1 {
		for n%d == 0 {
			factors = append(factors, d)
			n /= d
		}
		d++
		if d*d > n {
			if n > 1 {
				factors = append(factors, n)
				break
			}
		}
	}
	return factors
}

// ConvertToNewSystem converts a base-10 number to the new system representation.
func ConvertToNewSystem(n int) ([]int, primeMultiplicativeNumber) {
	exponentCounts := make([]int, 172)
	factors := Factorise(n)
	for _, factor := range factors {
		// make the counts in reverse orderd
		primePos := p(factor)
		exponentCounts[172-primePos]++
	}
	newNum := NewPrimeMultiplicativeNumber(exponentCounts)
	return factors, newNum // Trim the last '*'
}

func GenerateOverflowLaTeXTable(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	modulus := 1021

	fmt.Fprintln(file, "% paste this into your latex file")
	fmt.Fprintln(file, "\\begin{landscape}")                    // Start landscape mode
	fmt.Fprintln(file, "\\begin{longtable}{|c|c|c|c|p{10cm}|}") // Adjust the width in the p{width} as needed
	fmt.Fprintln(file, "\\caption{Extended Prime Exponents Under Modulus 1021} \\\\")
	fmt.Fprintln(file, "\\hline")
	fmt.Fprintln(file, "Prime^Exponent & Value & Modulus 1021 & Factors & New System Rep \\\\")
	fmt.Fprintln(file, "\\hline")
	fmt.Fprintln(file, "\\endfirsthead")
	fmt.Fprintln(file, "\\hline")
	fmt.Fprintln(file, "Prime^Exponent & Value & Modulus 1021 & Factors & New System Rep \\\\")
	fmt.Fprintln(file, "\\hline")
	fmt.Fprintln(file, "\\endhead")
	fmt.Fprintln(file, "\\hline")
	fmt.Fprintln(file, "\\endfoot")

	for _, prime := range primes {
		maxExp := maxExponent(prime, modulus)
		for exp := maxExp + 1; exp <= 2*maxExp; exp++ {
			value := math.Pow(float64(prime), float64(exp))
			modValue := int(value) % modulus
			primeFactors, newNumberSystemVal := ConvertToNewSystem(modValue)
			fmt.Fprintf(file, "\\(%d^{%d}\\) & %.1f & %d & %v & %s \\\\\n", prime, exp, value, modValue, primeFactors, newNumberSystemVal.intVal.Val())
		}
	}
	fmt.Fprintln(file, "\\end{longtable}")
	fmt.Fprintln(file, "\\end{landscape}") // End landscape mode
	return nil
}

func p(p int) int {
	// Lookup in the map
	if index, found := primeIndexMap[p]; found {
		return index
	} else {
		return -1
	}
}

func maxExponent(base, limit int) int {
	count := 0
	product := 1
	for product*base <= limit {
		product *= base
		count++
	}
	return count
}

func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
