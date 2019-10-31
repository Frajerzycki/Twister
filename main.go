package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/Frajerzycki/GONSE"
	"io/ioutil"
	"math/big"
	"os"
)

func printUsage() {
	fmt.Printf("Usage:\t%v -g [arguments]\tGenerate NSE secret key\n", os.Args[0])
	fmt.Printf("or:\t%v -e [arguments]\tEncrypt data with NSE algorithm\n", os.Args[0])
	fmt.Println("Arguments:")
	fmt.Println("\t-s <size>\tSet desired size of key in bits to <size>, if not used size will be 256 bits")
	//fmt.Println("\t-f <path>\tSet data to content of file placed in <path>")
	fmt.Println("\t-i\t\tSet data to input from stdin, this data have to end with EOF")
	fmt.Println("\t-k <key>\tSet key to <key>")
	//fmt.Println("\t-kf <path>\tSet key to integer parsed from content of file placed in <path>")
	os.Exit(1)
}

func generateKey(keySize uint) (*big.Int, error) {
	max := big.NewInt(1)
	max.Lsh(max, keySize)
	return rand.Int(rand.Reader, max)
}

func randomBytes(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func getDataFromStd() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	return ioutil.ReadAll(reader)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	var keySize uint = 256
	var err error
	var data []byte
	var key *big.Int
	for index := 2; index < len(os.Args); index++ {
		err = nil
		switch os.Args[index] {
		case "-s":
			keySize, err = parseKeySize(&index)
		case "-i":
			data, err = getDataFromStd()
		case "-k":
			key, err = parseKey(&index)
		case "-kf":
			key, err = parseKeyFromFile(&index)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	switch os.Args[1] {
	case "-g":
		key, err := generateKey(keySize)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(key.Text(keyBase))
	case "-e":
		var salt []byte
		var IV []int8
		salt, err = randomBytes(16)
		IV, err := nse.GenerateIV(len(data))
		if err != nil {
			fmt.Println(err)
			return
		}

		ciphertext, err := nse.Encrypt(data, salt, IV, key)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Only for testing
		fmt.Println(ciphertext)
	}

}
