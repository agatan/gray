package types

import (
	"fmt"

	"github.com/agatan/gray/token"
)

type infixOpFunc func(c *Checker, op string, lhs, rhs Type) (Type, error)

func opNumeric(c *Checker, op string, lhs, rhs Type) (Type, error) {
	lbasic, ok := lhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", lhs, op),
		}
	}
	rbasic, ok := rhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", rhs, op),
		}
	}
	if lbasic.Kind() == Int && rbasic.Kind() == Int {
		return BasicTypes[Int], nil
	}
	return nil, &Error{
		Message: fmt.Sprintf("invalid type %s and %s for operator '%s'", lhs, rhs, op),
	}
}

func opNumOrString(c *Checker, op string, lhs, rhs Type) (Type, error) {
	lbasic, ok := lhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", lhs, op),
		}
	}
	rbasic, ok := rhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", rhs, op),
		}
	}
	if lbasic.Kind() == Int && rbasic.Kind() == Int {
		return BasicTypes[Int], nil
	}
	if lbasic.Kind() == String && rbasic.Kind() == String {
		return BasicTypes[String], nil
	}
	return nil, &Error{
		Message: fmt.Sprintf("invalid type %s and %s for operator '%s'", lhs, rhs, op),
	}
}

func opLogical(c *Checker, op string, lhs, rhs Type) (Type, error) {
	lbasic, ok := lhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", lhs, op),
		}
	}
	rbasic, ok := rhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", rhs, op),
		}
	}
	if lbasic.Kind() == Bool && rbasic.Kind() == Bool {
		return BasicTypes[Bool], nil
	}
	return nil, &Error{
		Message: fmt.Sprintf("invalid type %s and %s for operator '%s'", lhs, rhs, op),
	}
}

func opComparison(c *Checker, op string, lhs, rhs Type) (Type, error) {
	lbasic, ok := lhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", lhs, op),
		}
	}
	rbasic, ok := rhs.(*Basic)
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("invalid type %s for operator '%s'", rhs, op),
		}
	}
	if lbasic.Kind() == rbasic.Kind() {
		return BasicTypes[Bool], nil
	}
	return nil, &Error{
		Message: fmt.Sprintf("invalid type %s and %s for operator '%s'", lhs, rhs, op),
	}
}

func opAssign(c *Checker, op string, lhs, rhs Type) (Type, error) {
	ref, ok := lhs.(*InstType)
	if !ok || !c.isCompatibleType(ref.Base(), builtinGenericTypes[refType]) {
		return nil, &Error{
			Message: fmt.Sprintf("type mismatch: expected reference type, but got %s", ref),
		}
	}
	if !c.isCompatibleType(ref.Args()[0], rhs) {
		return nil, &Error{
			Message: fmt.Sprintf("type mismatch: expected %s, but got %s", ref.Args()[0], rhs),
		}
	}
	return BasicTypes[Unit], nil
}

var builtinOperators = map[string]infixOpFunc{
	"-": opNumeric,
	"*": opNumeric,
	"/": opNumeric,
	"%": opNumeric,

	"+": opNumOrString,

	"&&": opLogical,
	"||": opLogical,

	"==": opComparison,
	"!=": opComparison,
	"<":  opComparison,
	"<=": opComparison,
	">":  opComparison,
	">=": opComparison,

	":=": opAssign,
}

func (c *Checker) checkInfixExpr(s *Scope, op string, lhs, rhs Type, pos token.Position) (Type, error) {
	f, ok := builtinOperators[op]
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("unknown operator: '%s'", op),
			Pos:     pos,
		}
	}
	ty, err := f(c, op, lhs, rhs)
	if err != nil {
		err.(*Error).Pos = pos
		return nil, err
	}
	return ty, nil
}
