package parser

import "math/big"

type DataSource struct {
	IsStdin bool
	path    string
}

type Arguments struct {
	KeySize          uint
	Key              *big.Int
	Source           DataSource
	IsInputDataText  bool
	IsOutputDataText bool
	IsInputKeyText   bool
	IsOutputKeyText  bool
}
