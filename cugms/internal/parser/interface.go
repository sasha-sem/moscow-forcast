package parser

type Parser interface {
	Parse() error
	Write() error
}
