package parser

import (
	"ok/ast"
	"ok/lexer"
)

func consumeExpr(parser *Parser, offset int) (ast.Node, int, error) {
	originalOffset := offset

	// Consume as many subexpressions and operators as possible.
	var parts []interface{}
	var err error
	didJustConsumeLiteral := false
	for offset < len(parser.File.Tokens) {
		tok := parser.File.Tokens[offset]

		// Bail out if we reach the end of the line. However, be careful that we
		// don't look back to the previous token that is from a previous expr
		// read (outside the scope of this expression).
		if offset > originalOffset && parser.File.Tokens[offset-1].IsEndOfLine {
			break
		}

		// Some tokens signal the end of the expression.
		if tok.Kind == lexer.TokenColon ||
			tok.Kind == lexer.TokenComma {
			break
		}

		// Try to consume a literal.
		var literal *ast.Literal
		literal, offset, err = consumeLiteral(parser.File, offset)
		if err == nil {
			err = validateLiteral(literal)
			if err != nil {
				// This kind of error should not stop the parsing.
				parser.Errors = append(parser.Errors, err)
			}

			parts = append(parts, literal)
			didJustConsumeLiteral = true
			continue
		}

		// Array
		var array *ast.Array
		array, offset, err = consumeArray(parser, offset)
		if err == nil {
			parts = append(parts, array)
			didJustConsumeLiteral = false
			continue
		}

		// Map
		var m *ast.Map
		m, offset, err = consumeMap(parser, offset)
		if err == nil {
			parts = append(parts, m)
			didJustConsumeLiteral = false
			continue
		}

		// Call
		var call *ast.Call
		call, offset, err = consumeCall(parser, offset)
		if err == nil {
			parts = append(parts, call)
			didJustConsumeLiteral = false
			continue
		}

		// Grouping "()"
		var group *ast.Group
		group, offset, err = consumeGroup(parser, offset)
		if err == nil {
			parts = append(parts, group)
			didJustConsumeLiteral = false
			continue
		}

		// Unary operator (can not be consumed directly after a literal).
		if !didJustConsumeLiteral {
			var unary *ast.Unary
			unary, offset, err = consumeUnary(parser, offset)
			if err == nil {
				parts = append(parts, unary)
				didJustConsumeLiteral = false
				continue
			}
		}

		// An identifier. It's important this happens last if all else fails.
		var identifier *ast.Identifier
		identifier, offset, err = consumeIdentifier(parser, offset)
		if err == nil {

			// Read ahead for key expression.
			if parser.File.Tokens[offset].Kind == lexer.TokenSquareOpen {
				offset++ // skip "["

				var key ast.Node
				key, offset, err = consumeExpr(parser, offset)
				if err != nil {
					return nil, originalOffset, err
				}

				offset, err = consume(parser.File, offset, []string{lexer.TokenSquareClose})
				if err != nil {
					return nil, originalOffset, err
				}

				parts = append(parts, &ast.Key{
					Expr: identifier,
					Key:  key,
				})
			} else {
				parts = append(parts, identifier)
			}

			continue
		}

		// Otherwise it must be a a valid binary operator.
		switch tok.Kind {
		case
			// Arithmetic
			lexer.TokenPlus, lexer.TokenMinus, lexer.TokenTimes,
			lexer.TokenDivide, lexer.TokenRemainder,

			// Logical
			lexer.TokenAnd, lexer.TokenOr,

			// Comparison
			lexer.TokenEqual, lexer.TokenNotEqual,
			lexer.TokenGreaterThan, lexer.TokenGreaterThanEqual,
			lexer.TokenLessThan, lexer.TokenLessThanEqual,

			// Assignment
			lexer.TokenAssign, lexer.TokenPlusAssign, lexer.TokenMinusAssign,
			lexer.TokenTimesAssign, lexer.TokenDivideAssign,
			lexer.TokenRemainderAssign:

			parts = append(parts, tok)
			offset++
			didJustConsumeLiteral = false
			continue
		}

		break
	}

	if len(parts) == 0 {
		// Nothing was consumed, this is an error.
		return nil, originalOffset, newTokenMismatch("expression",
			parser.File.Tokens[originalOffset-1].Kind,
			parser.File.Tokens[originalOffset].Kind)
	}

	return reduceExpr(parts), offset, nil
}

func consumeExprs(parser *Parser, offset int) ([]ast.Node, int, error) {
	// There must always be one expression.
	var expr ast.Node
	var err error
	expr, offset, err = consumeExpr(parser, offset)
	if err != nil {
		return nil, offset, err
	}

	// Any others are optional, but if there is a comma it must be followed by
	// the next expression.
	exprs := []ast.Node{expr}
	for {
		if parser.File.Tokens[offset].Kind != lexer.TokenComma {
			break
		}

		offset++ // skip comma
		expr, offset, err = consumeExpr(parser, offset)
		if err != nil {
			return nil, offset, err
		}

		exprs = append(exprs, expr)
	}

	return exprs, offset, nil
}

var operatorPrecedence = map[string]int{
	lexer.TokenAssign:          1,
	lexer.TokenPlusAssign:      1,
	lexer.TokenMinusAssign:     1,
	lexer.TokenTimesAssign:     1,
	lexer.TokenDivideAssign:    1,
	lexer.TokenRemainderAssign: 1,

	lexer.TokenOr: 2,

	lexer.TokenAnd: 3,

	lexer.TokenEqual:            4,
	lexer.TokenNotEqual:         4,
	lexer.TokenGreaterThan:      4,
	lexer.TokenGreaterThanEqual: 4,
	lexer.TokenLessThan:         4,
	lexer.TokenLessThanEqual:    4,

	lexer.TokenPlus:  5,
	lexer.TokenMinus: 5,

	lexer.TokenTimes:     6,
	lexer.TokenDivide:    6,
	lexer.TokenRemainder: 6,
}

func reduceExpr(parts []interface{}) ast.Node {
	if len(parts) == 1 {
		return parts[0]
	}

	if len(parts) == 3 {
		return &ast.Binary{
			Left:  parts[0],
			Op:    parts[1].(lexer.Token).Kind,
			Right: parts[2],
		}
	}

	leftPrecedence := operatorPrecedence[parts[1].(lexer.Token).Kind]
	rightPrecedence := operatorPrecedence[parts[len(parts)-2].(lexer.Token).Kind]
	if leftPrecedence > rightPrecedence {
		return &ast.Binary{
			Left:  reduceExpr(parts[:len(parts)-2]),
			Op:    parts[len(parts)-2].(lexer.Token).Kind,
			Right: parts[len(parts)-1].(*ast.Literal),
		}
	}

	return &ast.Binary{
		Left:  parts[0],
		Op:    parts[1].(lexer.Token).Kind,
		Right: reduceExpr(parts[2:]),
	}
}
