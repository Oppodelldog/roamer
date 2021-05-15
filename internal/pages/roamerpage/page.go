package roamerpage

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/config"
	"html/template"
	"io"
	"io/fs"
)

type Data struct {
	Title      string
	CssFile    string
	ActionRows ActionRows
}
type ActionRows []Actions
type Actions []Action
type Action struct {
	Icon    string
	Action  string
	Caption string
}

func Render(pageId string, fs fs.FS, writer io.Writer) error {
	var (
		basePath                 = pageId
		roamerPage, configExists = config.RoamerPage(basePath)
	)
	if !configExists {
		return fmt.Errorf("config not found for roamer page '%s'", pageId)
	}

	tpl, err := template.ParseFS(fs, "html/roamer-page.html")
	if err != nil {
		return err
	}

	var actionRows = actionRows(roamerPage)

	return tpl.Execute(writer, Data{
		Title:      roamerPage.Title,
		CssFile:    roamerPage.CSSFile,
		ActionRows: actionRows,
	})

}

func actionRows(roamerPage config.Page) ActionRows {
	var (
		actionRows ActionRows
		actions    Actions
		no         = 0
	)

	for _, action := range roamerPage.Actions {
		actions = append(actions, Action{
			Icon:    action.Icon,
			Action:  action.Action,
			Caption: action.Caption,
		})

		no++
		for _, col := range roamerPage.Columns {
			if no == col {
				actionRows = append(actionRows, actions)
				actions = Actions{}
			}
		}
	}

	if len(actions) > 0 {
		actionRows = append(actionRows, actions)
	}

	return actionRows
}
