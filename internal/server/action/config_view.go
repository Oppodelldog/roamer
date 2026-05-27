package action

import (
	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/script"
)

type (
	ConfigView struct {
		Title       string              `json:"Title"`
		WelcomeText string              `json:"WelcomeText"`
		Pages       map[string]PageView `json:"Pages"`
	}
	PageView struct {
		TitleShort string       `json:"TitleShort"`
		Title      string       `json:"Title"`
		CSSFile    string       `json:"CssFile"`
		Theme      config.Theme `json:"Theme"`
		ThemeClass string       `json:"ThemeClass"`
		Columns    []int        `json:"Columns"`
		Actions    []ActionView `json:"Actions"`
	}
	ActionView struct {
		Icon     string          `json:"Icon"`
		Caption  string          `json:"Caption"`
		Sequence string          `json:"Sequence"`
		Meta     script.Metadata `json:"Meta"`
	}
)

func configView(cfg config.Config) ConfigView {
	view := ConfigView{
		Title:       cfg.Title,
		WelcomeText: cfg.WelcomeText,
		Pages:       map[string]PageView{},
	}

	for pageId, page := range cfg.Pages {
		pageView := PageView{
			TitleShort: page.TitleShort,
			Title:      page.Title,
			CSSFile:    page.CSSFile,
			Theme:      pageTheme(page),
			ThemeClass: pageThemeClass(pageId, page),
			Columns:    page.Columns,
			Actions:    make([]ActionView, 0, len(page.Actions)),
		}

		for _, action := range page.Actions {
			meta, _ := script.Analyze(action.Sequence)
			pageView.Actions = append(pageView.Actions, ActionView{
				Icon:     action.Icon,
				Caption:  action.Caption,
				Sequence: action.Sequence,
				Meta:     meta,
			})
		}

		view.Pages[pageId] = pageView
	}

	return view
}

func pageTheme(page config.Page) config.Theme {
	if page.Theme != (config.Theme{}) {
		return page.Theme
	}

	switch page.CSSFile {
	case "7d2d.css":
		return config.Theme{
			BackgroundImage: "/img/background/forest.jpg",
			BackgroundColor: "#131311",
			AccentColor:     "#9f3025",
			CardColor:       "#ecf1e8",
			CardTextColor:   "#141714",
		}
	case "altf4.css":
		return config.Theme{
			BackgroundImage: "/img/background/knight.jpg",
			BackgroundColor: "#131311",
			AccentColor:     "#a9a59b",
			CardColor:       "#ecf1e8",
			CardTextColor:   "#141714",
		}
	case "remnant.css":
		return config.Theme{
			BackgroundImage: "/img/background/remnant.jpg",
			BackgroundColor: "#111315",
			AccentColor:     "#b9d7ff",
			CardColor:       "#f8f9f6",
			CardTextColor:   "#141714",
		}
	case "rust.css":
		return config.Theme{
			BackgroundImage: "/img/background/rust.jpg",
			BackgroundColor: "#171814",
			AccentColor:     "#b7d36a",
			CardColor:       "#ecf1e8",
			CardTextColor:   "#141714",
		}
	case "valheim.css":
		return config.Theme{
			BackgroundImage: "/img/background/wood.jpg",
			BackgroundColor: "#171814",
			AccentColor:     "#ffd36a",
			CardColor:       "#efe6d6",
			CardTextColor:   "#141714",
		}
	default:
		return config.Theme{
			BackgroundImage: "/img/background/index.jpg",
			BackgroundColor: "#141f35",
			AccentColor:     "#ffe36d",
			CardColor:       "#ecf1e8",
			CardTextColor:   "#141714",
		}
	}
}

func pageThemeClass(pageId string, page config.Page) string {
	switch page.CSSFile {
	case "7d2d.css":
		return "theme-7d2d"
	case "altf4.css":
		return "theme-altf4"
	case "remnant.css":
		return "theme-remnant"
	case "rust.css":
		return "theme-rust"
	case "valheim.css":
		return "theme-valheim"
	}

	switch pageId {
	case "7d2d", "altf4", "remnant", "rust", "valheim":
		return "theme-" + pageId
	default:
		return "theme-default"
	}
}
