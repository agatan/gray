package types

import (
	"fmt"

	"github.com/agatan/gray/token"
)

type infixOpFunc func(op string, lhs, rhs Type) (Type, error)

func opNumeric(op string, lhs, rhs Type) (Type, error) {
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

func opNumOrString(op string, lhs, rhs Type) (Type, error) {
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

func opLogical(op string, lhs, rhs Type) (Type, error) {
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

func opComparison(op string, lhs, rhs Type) (Type, error) {
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
}

func (c *Checker) checkInfixExpr(s *Scope, op string, lhs, rhs Type, pos token.Position) (Type, error) {
	f, ok := builtinOperators[op]
	if !ok {
		return nil, &Error{
			Message: fmt.Sprintf("unknown operator: '%s'", op),
			Pos:     pos,
		}
	}
	ty, err := f(op, lhs, rhs)
	if err != nil {
		err.(*Error).Pos = pos
		return nil, err
	}
	return ty, nil
}
