package config

type YankType int

const (
	LineWise YankType = iota
	CharWise
	BlockWise
)

type Yank struct {
	PartialBuffer Buffer
	Type          YankType
}
