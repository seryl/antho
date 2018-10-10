package evaluators

import (
	"github.com/google/go-jsonnet"
)

type Evaluator interface {
	SetVM(vm *jsonnet.VM)
}

func CreateEvaluator(myString string) (ev Evaluator) {
	switch myString {
	case "equal", "eql":
	case "not_equal", "dne":
		ev =
	default:
	}

	return ev
}
