package main

import (
	"fmt"
	"github.com/ikcilrep/twister/internal/functionality"
	"github.com/ikcilrep/twister/internal/parser"
	"log"
	"math/big"
	"os"
)

func printUsage() {
	fmt.Printf("Usage:\t%v -g [arguments; required: -o; not required: -s]\t\tGenerate NSE secret key\n", os.Args[0])
	fmt.Printf("or:\t%v -e [arguments; required: -i, -o, -k]\t\t\tEncrypt data with NSE algorithm\n", os.Args[0])
	fmt.Printf("or:\t%v -d [arguments; required: -i, -o, -k]\t\t\tDecrypt data with NSE algorithm\n", os.Args[0])
	fmt.Println("Arguments:")
	fmt.Println("\t-s <size>\tSet desired size of key in bytes to <size>, if not used size will be 32 bytes.")
	fmt.Println("\t-i <path>\tSet data to the content of file placed in <path>")
	fmt.Println("\t-o <path>\tSet output to file placed in <path>.")
	fmt.Println("\t-k <path>\tSet key to key read from file placed in <path>.")
	os.Exit(1)
}

func isNSETransformation() bool {
	return os.Args[1] == "-e" || os.Args[1] == "-d"
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	arguments := parser.NewArguments()
	files, err := parser.ParseArguments(arguments)
	if err != nil {
		log.Fatalln(err)
	}
	defer closeFiles(files)

	var key *big.Int

	if isNSETransformation() {
		err = arguments.VerifyForNSETransformation()
		if err != nil {
			log.Fatalln(err)
		}

		key, err = arguments.GetKey()
		if err != nil {
			log.Fatalln(err)
		}
	}

	switch os.Args[1] {
	case "-g":
		err = arguments.VerifyForKeyGeneration()
		if err != nil {
			log.Fatalln(err)
		}

		err = functionality.GenerateKey(arguments)
	case "-e":
		err = functionality.Encrypt(key, arguments)
	case "-d":
		err = functionality.Decrypt(key, arguments)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
