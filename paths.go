package paperweight

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	template "html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Path(format string, original string) string {
	erased := make([]interface{}, 0)
	for _, s := range strings.Split(original, string(os.PathSeparator)) {
		erased = append(erased, s)
	}
	//println(erased)
	return fmt.Sprintf(format, erased...)
	//fmt.Println(o)
	//return o
}

func CreateWithDirs(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func LoadTemplates(source string) (*template.Template, error) {
	templates, err := template.ParseGlob(source)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func GlobRun(glob string, action func(string) error) error {

	paths, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return nil
	}
	//sem := make(chan int, 1_000)

	var g errgroup.Group

	for _, path := range paths {
		path := path
		g.Go(func() error {
			return action(path)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func HTMLRenderer[T any](template *template.Template, name string) Renderer[T] {
	return func(model T, wr io.Writer) error {
		err := template.ExecuteTemplate(wr, name, model)
		if err != nil {
			return err
		}
		return nil
	}
}

func SaveToDisk(file File, wr io.Writer) error {
	_, err := wr.Write(file.Contents)
	return err
}

//func PathFile(inputPath string) string {
//	split := strings.Split(inputPath, "/")
//	return filepath.Join("out", filepath.Join(split[0:len(split)-1]...), "index.html")
//}
