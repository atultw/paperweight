package paperweight

import "io"

func StaticPath(path string) func() (io.Writer, error) {
	return func() (io.Writer, error) {
		return CreateWithDirs(path)
	}
}

func FormatPath[T any](format string) func(*File, T) (io.Writer, error) {
	return func(file *File, _ T) (io.Writer, error) {
		p := Path(format, file.Path)
		return CreateWithDirs(p)
	}
}
