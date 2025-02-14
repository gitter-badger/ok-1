package instruction

import (
	"ok/ast"
)

// Not is a logical NOT of a bool.
type Not struct {
	Left, Result string
}

// Execute implements the Instruction interface for the VM.
func (ins *Not) Execute(registers map[string]*ast.Literal, _ *int) error {
	registers[ins.Result] = ast.NewLiteralBool(
		!(registers[ins.Left].Value == "true"),
	)

	return nil
}
