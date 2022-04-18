package paperweight

import (
	"html/template"
	"testing"
)

func TestOne(t *testing.T) {
	type Article struct {
		Title string
		Body  template.HTML
		Link  string
	}

	allArticles := make([]Article, 0)

	templates, err := LoadTemplates("example_site/in/templates/*.gohtml")
	if err != nil {
		return
	}

	err = MarkdownPipeline[Article]{
		Glob:             "example_site/in/posts/**/*.md",
		OutputPathFormat: "example_site/out/posts/%[4]s/index.html",
		Parse: func(path string, body template.HTML, f Frontmatter) (Article, error) {
			article := Article{
				Title: f["title"].(string),
				Body:  body,
				Link:  Path("posts/%[4]s/index.html", path),
			}
			allArticles = append(allArticles, article)
			return article, nil
		},
		Renderer: HTMLRenderer[Article](templates, "article.gohtml"),
	}.Run()

	err = CopyPipeline{
		Glob:         "example_site/in/posts/*/*.jpg",
		OutputFormat: "example_site/out/posts/%[4]s/%[5]s",
	}.Run()

	index, err := Open("example_site/out/index.html")
	err = HTMLRenderer[[]Article](templates, "homepage.gohtml")(allArticles, index)
	if err != nil {
		t.Fatal(err)
	}
}
