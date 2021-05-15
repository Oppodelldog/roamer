package index

import (
	"github.com/Oppodelldog/roamer/internal/config"
	"html/template"
	"io"
	"io/fs"
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
		return err
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
