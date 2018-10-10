package evaluators

import (
	"github.com/google/go-jsonnet"
)

// Equal will true if the target
type Equal struct {
	VM *jsonnet.VM
}

func (e *Equal) SetVM(vm *jsonnet.VM) {
	e.VM = vm
}

func (e *Equal) Run(input ValidatorInput) int {
	return 0
}

// NotEqual validates whether or not a condition does not equal
type NotEqual struct {
	VM *jsonnet.VM
}

func (ne *NotEqual) SetVM(vm *jsonnet.VM) {
	ne.VM = vm
}

func (ne *NotEqual) Run(input ValidatorInput) int {
	return 0
}
