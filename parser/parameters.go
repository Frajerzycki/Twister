package parser

import "math/big"

type Parameters struct {
	KeySize          uint
	Key              *big.Int
	Data             []byte
	IsInputDataText  bool
	IsOutputDataText bool
	IsInputKeyText   bool
	IsOutputKeyText  bool
}
