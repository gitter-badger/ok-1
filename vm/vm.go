package vm

import (
	"ok/ast"
	"ok/compiler"
)

// VM is an instance of a virtual machine to run ok instructions.
type VM struct {
	f *compiler.CompiledFunc
}

// NewVM will create a new VM ready to run the provided instructions.
func NewVM(f *compiler.CompiledFunc) *VM {
	return &VM{
		f: f,
	}
}

// Run will run the program.
func (vm *VM) Run() {
	registers := map[string]*ast.Literal{}
	totalInstructions := len(vm.f.Instructions)
	for i := 0; i < totalInstructions; i++ {
		vm.f.Instructions[i].Execute(registers, &i)
	}
}
