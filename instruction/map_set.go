package instruction

import (
	"ok/ast"
)

// MapSet sets a number value to an index.
type MapSet struct {
	Map, Key, Value string
}

// Execute implements the Instruction interface for the VM.
func (ins *MapSet) Execute(registers map[string]*ast.Literal, _ *int) error {
	key := registers[ins.Key].Value
	registers[ins.Map].Map[key] = registers[ins.Value]

	return nil
}
