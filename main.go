package main

import (
	"crypto/rand"
	"fmt"
	_ "github.com/Frajerzycki/GONSE"
	"math/big"
	"os"
	"strconv"
)

func printUsage() {
	fmt.Printf("Usage: %v -g [arguments]\tGenerate nse secret key\n", os.Args[0])
	fmt.Println("Arguments:")
	fmt.Println("\t-s <size>\tSet desired size of key in bits to <size>, if not used size will be 256 bits")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	var keySize uint = 256
	lenArgs := len(os.Args)
	for index := 2; index < lenArgs; index++ {
		switch os.Args[index] {
		case "-s":
			index++
			if index < lenArgs {
				if size, err := strconv.Atoi(os.Args[index]); err == nil {
					keySize = uint(size)
					continue
				}
			}
			print("Size of key is not an integer!")
		}
	}
	switch os.Args[1] {
	case "-g":
		max := big.NewInt(1)
		max.Lsh(max, keySize)
		key, err := rand.Int(rand.Reader, max)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(key.Text(16))
	}
}
