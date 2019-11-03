package main

import (
	"./parser"
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
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
	fmt.Println("\t-b[i|o][k|d]\tSet input/output key/data format to binary, if not used format will be text\n\t\t\tAbove parameter doesn't matter on input data for encryption, beacuse in that context there isn't any difference.")
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

	arguments := parser.NewArguments()
	err := parser.ParseArguments(&arguments)
	if err != nil {
		fmt.Println(err)
		return
	}

	var data []byte
	key := new(big.Int)
	if arguments.KeyInput.Reader != nil {
		var keyBytes []byte
		keyBytes, err = ioutil.ReadAll(arguments.KeyInput.Reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		if arguments.KeyInput.IsBinary {
			key.SetBytes(keyBytes)
		} else {
			if lastIndex := len(keyBytes) - 1; keyBytes[lastIndex] == '\n' {
				keyBytes = keyBytes[:lastIndex]
			}
			key.SetString(string(keyBytes), parser.KeyBase)
		}
	} else {
		key = nil
	}

	if doesRequireKey() {
		if key == nil {
			fmt.Printf("Key is not set but option %v requires it.\n", os.Args[1])
			return
		}
		data, err = ioutil.ReadAll(arguments.DataInput.Reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(data) == 0 {
			fmt.Println("Data length has to be positive.")
			return
		}
	}

	switch os.Args[1] {
	case "-g":
		key, err = generateKey(arguments.KeySize)
		if err != nil {
			fmt.Println(err)
			return
		}
		if arguments.KeyOutput.IsBinary {
			arguments.KeyOutput.Writer.Write(key.Bytes())
		} else {
			arguments.KeyOutput.Writer.Write([]byte(fmt.Sprintf("%v\n", key.Text(parser.KeyBase))))
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
		bytes, _ := nse.Int64sToBytes(ciphertext)
		bytes = append(bytes, nse.Int8sToBytes(IV)...)
		buffer := make([]byte, 8)
		binary.PutUvarint(buffer, uint64(len(data)))
		bytes = append(bytes, buffer...)
		bytes = append(bytes, salt...)
		// Only for testing
		if arguments.DataOutput.IsBinary {
			arguments.DataOutput.Writer.Write(bytes)
		} else {
			arguments.DataOutput.Writer.Write([]byte(fmt.Sprintf("%v\n", base64.StdEncoding.EncodeToString(bytes))))
		}
	}

}
