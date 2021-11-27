package index

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"github.com/Oppodelldog/roamer/internal/config"
)

type Data struct {
	Title       string
	WelcomeText string
	Pages       Pages
}

type Pages map[string]Page

type Page struct{ TitleShort string }

func Render(fs fs.FS, writer io.Writer) error {
	tpl, err := template.ParseFS(fs, "html/index.html")
	if err != nil {
		return fmt.Errorf("cannot parse index page template: %w", err)
	}

	var pages = Pages{}
	for name, page := range config.Roamer().Pages {
		pages[name] = Page{TitleShort: page.TitleShort}
	}

	return tpl.Execute(writer, Data{
		Title:       config.Roamer().Title,
		WelcomeText: config.Roamer().WelcomeText,
		Pages:       pages,
	})
}
