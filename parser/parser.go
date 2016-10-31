//line parser.go.y:2
package parser

import __yyfmt__ "fmt"

//line parser.go.y:2
import (
	"github.com/agatan/gray/ast"
	"github.com/agatan/gray/token"
)

//line parser.go.y:13
type parserSymType struct {
	yys  int
	expr ast.Expr
	tok  token.Token
}

const IDENT = 57346
const INT = 57347
const TRUE = 57348
const FALSE = 57349

var parserToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"IDENT",
	"INT",
	"TRUE",
	"FALSE",
}
var parserStatenames = [...]string{}

const parserEofCode = 1
const parserErrCode = 2
const parserInitialStackSize = 16

//line yacctab:1
var parserExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const parserNprod = 6
const parserPrivate = 57344

var parserTokenNames []string
var parserStates []string

const parserLast = 6

var parserAct = [...]int{

	3, 4, 5, 6, 1, 2,
}
var parserPact = [...]int{

	-4, -1000, -1000, -1000, -1000, -1000, -1000,
}
var parserPgo = [...]int{

	0, 5, 4,
}
var parserR1 = [...]int{

	0, 2, 1, 1, 1, 1,
}
var parserR2 = [...]int{

	0, 1, 1, 1, 1, 1,
}
var parserChk = [...]int{

	-1000, -2, -1, 4, 5, 6, 7,
}
var parserDef = [...]int{

	0, -2, 1, 2, 3, 4, 5,
}
var parserTok1 = [...]int{

	1,
}
var parserTok2 = [...]int{

	2, 3, 4, 5, 6, 7,
}
var parserTok3 = [...]int{
	0,
}

var parserErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	parserDebug        = 0
	parserErrorVerbose = false
)

type parserLexer interface {
	Lex(lval *parserSymType) int
	Error(s string)
}

type parserParser interface {
	Parse(parserLexer) int
	Lookahead() int
}

type parserParserImpl struct {
	lval  parserSymType
	stack [parserInitialStackSize]parserSymType
	char  int
}

func (p *parserParserImpl) Lookahead() int {
	return p.char
}

func parserNewParser() parserParser {
	return &parserParserImpl{}
}

const parserFlag = -1000

func parserTokname(c int) string {
	if c >= 1 && c-1 < len(parserToknames) {
		if parserToknames[c-1] != "" {
			return parserToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func parserStatname(s int) string {
	if s >= 0 && s < len(parserStatenames) {
		if parserStatenames[s] != "" {
			return parserStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func parserErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !parserErrorVerbose {
		return "syntax error"
	}

	for _, e := range parserErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + parserTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := parserPact[state]
	for tok := TOKSTART; tok-1 < len(parserToknames); tok++ {
		if n := base + tok; n >= 0 && n < parserLast && parserChk[parserAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if parserDef[state] == -2 {
		i := 0
		for parserExca[i] != -1 || parserExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; parserExca[i] >= 0; i += 2 {
			tok := parserExca[i]
			if tok < TOKSTART || parserExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if parserExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += parserTokname(tok)
	}
	return res
}

func parserlex1(lex parserLexer, lval *parserSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = parserTok1[0]
		goto out
	}
	if char < len(parserTok1) {
		token = parserTok1[char]
		goto out
	}
	if char >= parserPrivate {
		if char < parserPrivate+len(parserTok2) {
			token = parserTok2[char-parserPrivate]
			goto out
		}
	}
	for i := 0; i < len(parserTok3); i += 2 {
		token = parserTok3[i+0]
		if token == char {
			token = parserTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = parserTok2[1] /* unknown char */
	}
	if parserDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", parserTokname(token), uint(char))
	}
	return char, token
}

func parserParse(parserlex parserLexer) int {
	return parserNewParser().Parse(parserlex)
}

func (parserrcvr *parserParserImpl) Parse(parserlex parserLexer) int {
	var parsern int
	var parserVAL parserSymType
	var parserDollar []parserSymType
	_ = parserDollar // silence set and not used
	parserS := parserrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	parserstate := 0
	parserrcvr.char = -1
	parsertoken := -1 // parserrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		parserstate = -1
		parserrcvr.char = -1
		parsertoken = -1
	}()
	parserp := -1
	goto parserstack

ret0:
	return 0

ret1:
	return 1

parserstack:
	/* put a state and value onto the stack */
	if parserDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", parserTokname(parsertoken), parserStatname(parserstate))
	}

	parserp++
	if parserp >= len(parserS) {
		nyys := make([]parserSymType, len(parserS)*2)
		copy(nyys, parserS)
		parserS = nyys
	}
	parserS[parserp] = parserVAL
	parserS[parserp].yys = parserstate

parsernewstate:
	parsern = parserPact[parserstate]
	if parsern <= parserFlag {
		goto parserdefault /* simple state */
	}
	if parserrcvr.char < 0 {
		parserrcvr.char, parsertoken = parserlex1(parserlex, &parserrcvr.lval)
	}
	parsern += parsertoken
	if parsern < 0 || parsern >= parserLast {
		goto parserdefault
	}
	parsern = parserAct[parsern]
	if parserChk[parsern] == parsertoken { /* valid shift */
		parserrcvr.char = -1
		parsertoken = -1
		parserVAL = parserrcvr.lval
		parserstate = parsern
		if Errflag > 0 {
			Errflag--
		}
		goto parserstack
	}

parserdefault:
	/* default state action */
	parsern = parserDef[parserstate]
	if parsern == -2 {
		if parserrcvr.char < 0 {
			parserrcvr.char, parsertoken = parserlex1(parserlex, &parserrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if parserExca[xi+0] == -1 && parserExca[xi+1] == parserstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			parsern = parserExca[xi+0]
			if parsern < 0 || parsern == parsertoken {
				break
			}
		}
		parsern = parserExca[xi+1]
		if parsern < 0 {
			goto ret0
		}
	}
	if parsern == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			parserlex.Error(parserErrorMessage(parserstate, parsertoken))
			Nerrs++
			if parserDebug >= 1 {
				__yyfmt__.Printf("%s", parserStatname(parserstate))
				__yyfmt__.Printf(" saw %s\n", parserTokname(parsertoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for parserp >= 0 {
				parsern = parserPact[parserS[parserp].yys] + parserErrCode
				if parsern >= 0 && parsern < parserLast {
					parserstate = parserAct[parsern] /* simulate a shift of "error" */
					if parserChk[parserstate] == parserErrCode {
						goto parserstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if parserDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", parserS[parserp].yys)
				}
				parserp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if parserDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", parserTokname(parsertoken))
			}
			if parsertoken == parserEofCode {
				goto ret1
			}
			parserrcvr.char = -1
			parsertoken = -1
			goto parsernewstate /* try again in the same state */
		}
	}

	/* reduction by production parsern */
	if parserDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", parsern, parserStatname(parserstate))
	}

	parsernt := parsern
	parserpt := parserp
	_ = parserpt // guard against "declared and not used"

	parserp -= parserR2[parsern]
	// parserp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if parserp+1 >= len(parserS) {
		nyys := make([]parserSymType, len(parserS)*2)
		copy(nyys, parserS)
		parserS = nyys
	}
	parserVAL = parserS[parserp+1]

	/* consult goto table to find next state */
	parsern = parserR1[parsern]
	parserg := parserPgo[parsern]
	parserj := parserg + parserS[parserp].yys + 1

	if parserj >= parserLast {
		parserstate = parserAct[parserg]
	} else {
		parserstate = parserAct[parserj]
		if parserChk[parserstate] != -parsern {
			parserstate = parserAct[parserg]
		}
	}
	// dummy call; replaced with literal code
	switch parsernt {

	case 1:
		parserDollar = parserS[parserpt-1 : parserpt+1]
		//line parser.go.y:24
		{
			parserVAL.expr = parserDollar[1].expr
			if l, ok := parserlex.(*Lexer); ok {
				l.expr = parserVAL.expr
			}
		}
	case 2:
		parserDollar = parserS[parserpt-1 : parserpt+1]
		//line parser.go.y:33
		{
			parserVAL.expr = &ast.Ident{Name: parserDollar[1].tok.Lit}
			parserVAL.expr.SetPosition(parserDollar[1].tok.Position())
		}
	case 3:
		parserDollar = parserS[parserpt-1 : parserpt+1]
		//line parser.go.y:38
		{
			parserVAL.expr = &ast.BasicLit{Kind: token.INT, Lit: parserDollar[1].tok.Lit}
			parserVAL.expr.SetPosition(parserDollar[1].tok.Position())
		}
	case 4:
		parserDollar = parserS[parserpt-1 : parserpt+1]
		//line parser.go.y:43
		{
			parserVAL.expr = &ast.BasicLit{Kind: token.BOOL, Lit: parserDollar[1].tok.Lit}
			parserVAL.expr.SetPosition(parserDollar[1].tok.Position())
		}
	case 5:
		parserDollar = parserS[parserpt-1 : parserpt+1]
		//line parser.go.y:48
		{
			parserVAL.expr = &ast.BasicLit{Kind: token.BOOL, Lit: parserDollar[1].tok.Lit}
			parserVAL.expr.SetPosition(parserDollar[1].tok.Position())
		}
	}
	goto parserstack /* stack new state and value */
}
