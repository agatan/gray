package token

type Kind int

const (
	ERROR Kind = iota
	EOF

	literal_beg
	IDENT
	UNIT
	BOOL
	INT
	literal_end
)

func (t Kind) IsLiteral() bool {
	return literal_beg < t && t < literal_end
}
