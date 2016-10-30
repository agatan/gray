package token

type Token struct {
	PosImpl
	Kind Kind
	Lit  string
}
