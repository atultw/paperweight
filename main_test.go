package paperweight

import (
	"fmt"
	"html/template"
	"io"
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
		t.Fatal(err)
	}

	source, err := GlobFileSource("example_site/in/posts/*/main.md")
	if err != nil {
		t.Fatal(err)
	}

	NewArticle := func(htm template.HTML, f Frontmatter, file File) (Article, error) {
		article := Article{
			Title: f["title"].(string),
			Body:  htm,
			Link:  Path("posts/%[4]s/index.html", file.Path),
		}
		allArticles = append(allArticles, article)
		return article, nil
	}

	err = MultiPipeline[File, Article]{
		Source:   &source,
		Parse:    Markup(NewArticle),
		Output:   FormatPath[Article]("example_site/out/posts/%[4]s/index.html"),
		Renderer: HTMLRenderer[Article](templates, "article.gohtml"),
	}.Run()

	if err != nil {
		t.Fatal(err)
	}

	source2, err := NewMultiNetworkSource("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		t.Fatal(err)
	}
	err = MultiPipeline[JSONItem, Article]{
		Source: &source2,
		Parse:  JSONDecode[Article],
		Output: func(_ *JSONItem, a Article) (io.Writer, error) {
			return CreateWithDirs(fmt.Sprintf("example_site/out/posts/%s/index.html", a.Title))
		},
		Renderer: HTMLRenderer[Article](templates, "article.gohtml"),
	}.Run()
	if err != nil {
		t.Fatal(err)
	}

	source, err = GlobFileSource("example_site/in/posts/*/*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	err = MultiPipeline[File, File]{
		Source:   &source,
		Parse:    Passthrough[File],
		Output:   FormatPath[File]("example_site/out/posts/%[4]s/%[5]s"),
		Renderer: SaveToDisk,
	}.Run()
	if err != nil {
		t.Fatal(err)
	}

	err = VariableSource[[]Article]{
		Data:     allArticles,
		Output:   StaticPath("example_site/out/index.html"),
		Renderer: HTMLRenderer[[]Article](templates, "homepage.gohtml"),
	}.Run()
	if err != nil {
		t.Fatal(err)
	}
}
