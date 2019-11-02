package main

import (
	"./parser"
	"bufio"
	"crypto/rand"
	"encoding/base64"
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
	fmt.Println("\t-t[i|o][k|d]\tSet input/output key/data format to text, if not used format will be binary\n\t\t\tAbove parameter doesn't matter on input data for encryption. Also, if key is given in terminal it is automatically set to text format.")
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

	parameters := parser.Arguments{KeySize: uint(256)}
	err := parser.ParseArguments(&parameters)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data []byte
	if parameters.Source.IsStdin {
		data, err = getDataFromStd()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	switch os.Args[1] {
	case "-g":
		parameters.Key, err = generateKey(parameters.KeySize)
		if err != nil {
			fmt.Println(err)
			return
		}
		if parameters.IsOutputDataText {
			fmt.Print(parameters.Key.Text(parser.KeyBase))
		}
	case "-e":
		var salt []byte
		var IV []int8
		salt, err = randomBytes(16)
		IV, err := nse.GenerateIV(len(data))
		if err != nil {
			fmt.Println(err)
			return
		}

		ciphertext, err := nse.Encrypt(data, salt, IV, parameters.Key)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Only for testing
		if parameters.IsOutputDataText {
			bytes, _ := nse.Int64sToBytes(ciphertext)
			fmt.Printf("\n%v\n", base64.StdEncoding.EncodeToString(bytes))
		}
	}

}
