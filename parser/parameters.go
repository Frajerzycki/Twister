package parser

import "math/big"

type Arguments struct {
	KeySize          uint
	Key              *big.Int
	Data             []byte
	IsInputDataText  bool
	IsOutputDataText bool
	IsInputKeyText   bool
	IsOutputKeyText  bool
}
