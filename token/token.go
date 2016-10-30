package token

type Token int

const (
	ERROR Token = iota
	EOF

	literal_beg
	IDENT
	UNIT
	BOOL
	INT
	literal_end
)

func (t Token) IsLiteral() bool {
	return literal_beg < t && t < literal_end
}
