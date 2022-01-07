package action

import (
	"fmt"

	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

func Worker(actions <-chan Action, broadcast chan<- []byte) {
	var (
		sequence = ""
		seq      = sequencer.New(1)
	)

	seq.BeforeSequence(func() {
		key.ResetPressed()
	})

	go func() {
		for action := range actions {
			switch v := action.(type) {
			case SequenceClearSequence:
				sequence = ""
			case SequenceSetConfigSequence:
				page, ok := config.RoamerPage(v.PageId)
				if !ok {
					fmt.Printf("roamer-page '%s' not found", v.PageId)
					continue
				}

				if len(page.Actions) <= v.SequenceIndex {
					fmt.Printf("roamer-page '%s' not found", v.PageId)
					continue
				}

				var action = page.Actions[v.SequenceIndex]

				elems, err := script.Parse(action.Sequence)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println(script.Write(elems))
				seq.EnqueueSequence(func() []sequencer.Elem {
					return elems
				})

				sequence = action.Sequence
			case SequencePause:
				err := seq.Pause()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case SequenceAbort:
				seq.Abort()

				sequence = ""
			case SequenceSave:
				saveSequence(v)
				continue
			case LoadSoundSettings:
				v.Response <- msgSoundSettings(getSoundSettings())
				continue
			case SetSoundVolume:
				setSoundVol(v.Id, v.Value)
			case SetMainSoundVolume:
				setMainSoundVol(v.Value)
				continue
			case PageNew:
				createNewPage()
				v.Response <- msgConfig(config.Roamer())
			case PageDelete:
				deletePage(v.PageId)
				v.Response <- msgConfig(config.Roamer())
			case SequenceNew:
				createNewSequence(v.PageId)
				v.Response <- msgConfig(config.Roamer())
			case SequenceDelete:
				deleteSequence(v.PageId, v.SequenceIndex)
				v.Response <- msgConfig(config.Roamer())
			case PagesSave:
				savePages(v.Pages)
				v.Response <- msgConfig(config.Roamer())
			default:
				fmt.Printf("unknown action: %T\n", action)
			}

			broadcast <- msgState(SequenceState{
				Sequence:    sequence,
				IsPaused:    seq.IsPaused(),
				HasSequence: seq.HasSequence()})
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
		v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, action.Sequence, false)

		return
	}

	err = config.SetSequence(v.PageId, v.SequenceIndex, v.Caption, v.Sequence)
	if err != nil {
		fmt.Println(err)
		v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, action.Sequence, false)

		return
	}

	err = config.Save()
	if err != nil {
		fmt.Println(err)
		v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, action.Sequence, false)

		return
	}

	v.Response <- msgSequenceSaveResult(v.PageId, v.SequenceIndex, v.Sequence, true)
}
