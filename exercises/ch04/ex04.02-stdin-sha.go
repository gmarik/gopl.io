package main

import (
	"crypto/sha256"
	"crypto/sha512"

	"fmt"
	"io/ioutil"
	"os"
)

// Exercise 4.2: Write a program that prints the SHA256 hash of its standard
// input by default but supports a command-line flag to print the SHA384 or
// SHA512 hash instead.

func main() {
	alg := "sha256"

	if len(os.Args) == 2 {
		alg = os.Args[1]
	}

	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch alg {
	case "sha256":
		fmt.Printf("%x\n", sha256.Sum256(in))
	case "sha384":
		fmt.Printf("%x\n", sha512.Sum384(in))
	case "sha512":
		fmt.Printf("%x\n", sha512.Sum512(in))
	default:
		fmt.Println("Unknown hash:", alg)
		os.Exit(1)
	}
}
