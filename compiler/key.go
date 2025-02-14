package compiler

import (
	"ok/ast"
	"ok/instruction"
	"strings"
)

func compileKey(compiledFunc *CompiledFunc, n *ast.Key) (string, string, error) {
	arrayOrMapRegister, arrayOrMapKind, err := compileExpr(compiledFunc, n.Expr)
	if err != nil {
		return "", "", err
	}

	// TODO(elliot): Check key is the correct type.
	keyRegister, _, err := compileExpr(compiledFunc, n.Key)
	if err != nil {
		return "", "", err
	}

	resultRegister := compiledFunc.nextRegister()
	if strings.HasPrefix(arrayOrMapKind, "[]") {
		compiledFunc.append(&instruction.ArrayGet{
			Array:  arrayOrMapRegister,
			Index:  keyRegister,
			Result: resultRegister,
		})
	} else {
		compiledFunc.append(&instruction.MapGet{
			Map:    arrayOrMapRegister,
			Key:    keyRegister,
			Result: resultRegister,
		})
	}

	// TODO(elliot): Does not return element type.
	return resultRegister, "", nil
}
