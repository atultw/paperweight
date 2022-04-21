package paperweight

import "io"

type VariableSource[T any] struct {
	Data     T
	Output   func() (io.Writer, error)
	Renderer Renderer[T]
}

func (v VariableSource[T]) Run() error {
	output, err := v.Output()
	if err != nil {
		return err
	}
	var wr io.Writer
	wr = output
	err = v.Renderer(v.Data, wr)
	if err != nil {
		return err
	}
	return nil
}
