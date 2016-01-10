package main

// Write a function that counts the number of bits that are different in two
// SHA256 hashes. (See PopCount from Section 2.6.2.)
//
//
import sha "crypto/sha256"
import "fmt"

func main() {

	sha1 := sha.Sum256([]byte("hello"))
	sha2 := sha.Sum256([]byte("world"))

	fmt.Println("Different bits:", CountDifferent(sha1, sha2))
	fmt.Printf("sha1:%x\n", sha1)
	fmt.Printf("sha2:%x\n", sha2)
}

// PopCount returns the population count (number of set bits) of x.
func PopCountClearLastBit(x byte) (sum int8) {
	for ; x > 0; x = x & (x - 1) {
		sum += 1
	}
	return sum
}

func CountDifferent(sha1, sha2 [32]byte) (sum int8) {
	for i := 0; i < 32; i += 1 {
		sum += PopCountClearLastBit(sha1[i] ^ sha2[i])
	}

	return sum
}
