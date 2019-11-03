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

func doesRequireKey() bool {
	return os.Args[1] == "-e"
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	parameters := parser.NewArguments()
	err := parser.ParseArguments(&parameters)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data []byte
	key := new(big.Int)
	if doesRequireKey() {
		data, err = ioutil.ReadAll(parameters.DataInput.Reader)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if parameters.KeyInput.Reader != nil {
		var keyBytes []byte
		keyBytes, err = ioutil.ReadAll(parameters.KeyInput.Reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		if parameters.KeyInput.IsText {
			if lastIndex := len(keyBytes) - 1; keyBytes[lastIndex] == '\n' {
				keyBytes = keyBytes[:lastIndex]
			}
			key.SetString(string(keyBytes), parser.KeyBase)
		} else {
			key.SetBytes(keyBytes)
		}
	}

	switch os.Args[1] {
	case "-g":
		key, err = generateKey(parameters.KeySize)
		if err != nil {
			fmt.Println(err)
			return
		}
		if parameters.KeyOutput.IsText {
			parameters.KeyOutput.Writer.Write([]byte(key.Text(parser.KeyBase) + "\n"))
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

		ciphertext, err := nse.Encrypt(data, salt, IV, key)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Only for testing
		if parameters.DataOutput.IsText {
			bytes, _ := nse.Int64sToBytes(ciphertext)
			fmt.Printf("\n%v\n", base64.StdEncoding.EncodeToString(bytes))
		}
	}

}
