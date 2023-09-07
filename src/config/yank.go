package config

type YankType int

const (
	EMPTY_YANK YankType = iota
	LineWise
	CharWise
	BlockWise
)

type Yank struct {
	PartialBuffer Buffer
	Type          YankType
}
