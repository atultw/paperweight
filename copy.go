package paperweight

import (
	"io"
	"os"
	"path/filepath"
)

type CopyPipeline struct {
	Glob         string
	OutputFormat string
}

func (p CopyPipeline) Run() error {
	paths, err := filepath.Glob(p.Glob)
	if err != nil {
		return err
	}

	for _, path := range paths {
		newPath := Path(p.OutputFormat, path)
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		out, err := CreateWithDirs(newPath)
		if err != nil {
			in.Close()
			return err
		}
		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
	}
	return nil
}
