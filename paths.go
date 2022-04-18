package paperweight

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/sync/errgroup"
	template "html/template"
	"io"
	"io/ioutil"
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

// todo: misnomer

func Open(path string) (*os.File, error) {
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

type Frontmatter map[string]interface{}

func MarkupWithFrontmatter(markdown []byte) (template.HTML, Frontmatter, error) {
	gm := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := gm.Convert(markdown, &buf, parser.WithContext(context)); err != nil {
		return "", Frontmatter{}, err
	}
	frontmatter := meta.Get(context)
	return template.HTML(buf.String()), frontmatter, nil
}

type Renderer[T any] func(model T, wr io.Writer) error

type MarkdownPipeline[T any] struct {
	Glob             string
	OutputPathFormat string
	Parse            func(string, template.HTML, Frontmatter) (T, error)
	Renderer         Renderer[T]
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

func (p MarkdownPipeline[T]) Run() error {
	err := GlobRun(p.Glob, func(path string) error {
		md, err := ioutil.ReadFile(filepath.Join(path))
		htm, frontmatter, err := MarkupWithFrontmatter(md)
		if err != nil {
			return err
		}

		outputPath := Path(p.OutputPathFormat, path)
		file, err := Open(outputPath)
		defer file.Close()

		if err != nil {
			return err
		}

		model, err := p.Parse(outputPath, htm, frontmatter)
		if err != nil {
			return err
		}

		err = p.Renderer(model, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

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
		out, err := os.Create(newPath)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}
		in.Close()
		out.Close()
	}
	return nil
}

//func PathFile(inputPath string) string {
//	split := strings.Split(inputPath, "/")
//	return filepath.Join("out", filepath.Join(split[0:len(split)-1]...), "index.html")
//}
