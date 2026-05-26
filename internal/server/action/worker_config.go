package action

import (
	"fmt"

	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/script"
)

func StartConfigWorker(actions <-chan Action, _ chan<- []byte) {
	go func() {
		for action := range actions {
			switch v := action.(type) {
			case PageNew:
				createNewPage()
				v.Response <- msgConfig(config.Roamer())
			case PagesSave:
				savePages(v.Pages)
				v.Response <- msgConfig(config.Roamer())
			case PageDelete:
				deletePage(v.PageId)
				v.Response <- msgConfig(config.Roamer())
			case SequenceNew:
				createNewSequence(v.PageId)
				v.Response <- msgConfig(config.Roamer())
			case SequenceSave:
				saveSequence(v)
			case SequenceFormat:
				formatSequence(v)
			case SequenceValidate:
				validateSequence(v)
			case SequenceDelete:
				deleteSequence(v.PageId, v.SequenceIndex)
				v.Response <- msgConfig(config.Roamer())
			case SequenceDuplicate:
				duplicateSequence(v.PageId, v.SequenceIndex)
				v.Response <- msgConfig(config.Roamer())
			case SequenceMove:
				moveSequence(v.PageId, v.SequenceIndex, v.Offset)
				v.Response <- msgConfig(config.Roamer())
			default:
				fmt.Printf("unknown worker action: %T\n", action)
			}
		}
	}()
}

func savePages(pages config.Pages) {
	err := config.SavePages(pages)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteSequence(pageId string, index int) {
	err := config.DeleteSequence(pageId, index)
	if err != nil {
		fmt.Println(err)
	}
}

func duplicateSequence(pageId string, index int) {
	err := config.DuplicateSequence(pageId, index)
	if err != nil {
		fmt.Println(err)
	}
}

func moveSequence(pageId string, index int, offset int) {
	err := config.MoveSequence(pageId, index, offset)
	if err != nil {
		fmt.Println(err)
	}
}

func createNewSequence(pageId string) {
	err := config.NewSequence(pageId)
	if err != nil {
		fmt.Println(err)
	}
}

func deletePage(id string) {
	err := config.DeletePage(id)
	if err != nil {
		fmt.Println(err)
	}
}

func createNewPage() {
	err := config.NewPage()
	if err != nil {
		fmt.Println(err)
	}
}

func saveSequence(v SequenceSave) {
	page, ok := config.RoamerPage(v.PageId)
	if !ok {
		fmt.Printf("roamer-page '%s' not found", v.PageId)

		return
	}

	if len(page.Actions) <= v.SequenceIndex {
		fmt.Printf("roamer-page '%s' not found", v.PageId)

		return
	}

	var action = page.Actions[v.SequenceIndex]

	_, err := script.Parse(v.Sequence)
	if err != nil {
		fmt.Println(err)
		meta, _ := script.Analyze(action.Sequence)
		v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, action.Sequence, meta, false, err.Error())

		return
	}

	err = config.SetSequence(v.PageId, v.SequenceIndex, v.Caption, v.Icon, v.Sequence)
	if err != nil {
		fmt.Println(err)
		meta, _ := script.Analyze(action.Sequence)
		v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, action.Sequence, meta, false, err.Error())

		return
	}

	err = config.Save()
	if err != nil {
		fmt.Println(err)
		meta, _ := script.Analyze(action.Sequence)
		v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, action.Sequence, meta, false, err.Error())

		return
	}

	meta, _ := script.Analyze(v.Sequence)
	v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, v.Sequence, meta, true, "")
}

func formatSequence(v SequenceFormat) {
	elems, err := script.Parse(v.Sequence)
	if err != nil {
		fmt.Println(err)
		meta := script.Metadata{Valid: false, Error: err.Error()}
		v.Response <- msgSequenceFormatResult(v.PageId, v.SequenceIndex, v.Sequence, meta, false, err.Error())

		return
	}

	formatted := script.Write(elems)
	meta := script.AnalyzeElems(elems)
	v.Response <- msgSequenceFormatResult(v.PageId, v.SequenceIndex, formatted, meta, true, "")
}

func validateSequence(v SequenceValidate) {
	meta, err := script.Analyze(v.Sequence)
	if err != nil {
		v.Response <- msgSequenceValidateResult(v.PageId, v.SequenceIndex, v.Sequence, meta, false, err.Error())

		return
	}

	v.Response <- msgSequenceValidateResult(v.PageId, v.SequenceIndex, v.Sequence, meta, true, "")
}
