package instruction

import (
	"fmt"
	"ok/ast"
	"strings"
)

// Len is used to determine the size of an array or map.
type Len struct {
	Argument, Result string
}

// Execute implements the Instruction interface for the VM.
func (ins *Len) Execute(registers map[string]*ast.Literal, _ *int) error {
	r := registers[ins.Argument]
	var result int
	if strings.HasPrefix(r.Kind, "[]") {
		result = len(r.Array)
	} else {
		result = len(r.Map)
	}

	registers[ins.Result] = ast.NewLiteralNumber(fmt.Sprintf("%d", result))

	return nil
}
