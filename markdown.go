package paperweight

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"io"
)

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

type MultiPipeline[I Source, T any] struct {
	Source   MultiSource[I]
	Parse    func(I) (T, error)
	Output   func(*I, T) (io.Writer, error)
	Renderer Renderer[T]
}

func (p MultiPipeline[I, T]) Run() error {
	for p.Source.Next() {
		mdfile, err := p.Source.Get()
		if err != nil {
			return err
		}

		//outputPath := Path(p.OutputPathFormat, path)
		//file, err := CreateWithDirs(outputPath)

		if err != nil {
			return err
		}

		model, err := p.Parse(mdfile)
		if err != nil {
			return err
		}

		wr, err := p.Output(&mdfile, model)
		if err != nil {
			return err
		}

		err = p.Renderer(model, wr)
		if err != nil {
			return err
		}

	}

	return nil
}

type Transformer[I Source, T any] func(I) (T, error)

func Markup[I Source, T any](transform func(htm template.HTML, f Frontmatter, source I) (T, error)) Transformer[I, T] {
	return func(source I) (T, error) {
		var res T
		htm, f, err := MarkupWithFrontmatter(source.Data())
		if err != nil {
			return res, err
		}
		return transform(htm, f, source)
	}
}
