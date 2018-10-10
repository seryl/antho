package validation

import (
	"github.com/google/go-jsonnet"
)

// RunnerInputs look like the following
type RunnerInput struct {
	filename string
}

type Runner struct {
	VM *jsonnet.VM
}

func NewRunner(file string, jpaths string[]) *Runner {
	return &Validator{
		VM: jsonnet.MakeVM(&jsonnet.FileImporter{
			JPaths: jpaths,
		})
	}
}

func (v *Validator) Evaluate() {

}
