package main

import (
	"fmt"
	"github.com/Frajerzycki/Twister/parser"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

func printUsage() {
	fmt.Printf("Usage:\t%v -g [arguments]\tGenerate NSE secret key\n", os.Args[0])
	fmt.Printf("or:\t%v -e [arguments]\tEncrypt data with NSE algorithm\n", os.Args[0])
	fmt.Println("Arguments:")
	fmt.Println("\t-s <size>\tSet desired size of key in bytes to <size>, if not used size will be 32 bytes")
	fmt.Println("\t-i <path>\tSet data to content of file placed in <path>")
	fmt.Println("\t-k <key>\tSet key to <key>")
	fmt.Println("\t-b[i|o][k|d]\tSet input/output key/data format to binary, if not used format will be text\n\t\t\tAbove parameter doesn't matter on input data for encryption (-bid), beacuse in that context there isn't any difference.")
	fmt.Println("\t-o <path>\tRedirect output to file placed in <path>, if not used output will be STDOUT.")
	fmt.Println("\t-kf <path>\tSet key to integer parsed from content of file placed in <path>")
	os.Exit(1)
}

func doesRequireKey() bool {
	return os.Args[1] == "-e"
}

func getKey(arguments *parser.Arguments) (*big.Int, error) {
	key := new(big.Int)
	if arguments.KeyInput.Reader != nil {
		keyBytes, err := ioutil.ReadAll(arguments.KeyInput.Reader)
		if err != nil {
			return nil, err
		}
		if arguments.KeyInput.IsBinary {
			return key.SetBytes(keyBytes), nil
		} else {
			if lastIndex := len(keyBytes) - 1; keyBytes[lastIndex] == '\n' {
				keyBytes = keyBytes[:lastIndex]
			}
			key.SetString(string(keyBytes), parser.KeyBase)
			return key, nil
		}

	}
	return nil, nil
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

	var data []byte
	var key *big.Int
	key, err = getKey(arguments)
	if err != nil {
		log.Fatalln(err)
	}

	if doesRequireKey() {
		if key == nil {
			log.Fatalf("Key is not set but option %v requires it.\n", os.Args[1])
		}
		data, err = ioutil.ReadAll(arguments.DataInput.Reader)
		if err != nil {
			log.Fatalln(err)
		}
		if len(data) == 0 {
			log.Fatalln(notPositiveLengthError)
		}
	}

	switch os.Args[1] {
	case "-g":
		err = generateKey(arguments.KeySize, arguments)
	case "-e":
		err = encrypt(data, key, arguments)
	}

	if err != nil {
		log.Fatalln(err)
	}
}
