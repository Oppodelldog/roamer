package script

import (
	"fmt"
	"time"

	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

type Metadata struct {
	Labels            []string
	HasLoop           bool
	UsesMouse         bool
	HoldsKeys         bool
	EstimatedDuration string
	Valid             bool
	Error             string
}

func Analyze(sequence string) (Metadata, error) {
	elems, err := Parse(sequence)
	if err != nil {
		return Metadata{
			Labels: []string{},
			Valid:  false,
			Error:  err.Error(),
		}, err
	}

	meta := AnalyzeElems(elems)
	meta.Valid = true

	return meta, nil
}

func AnalyzeElems(elems []sequencer.Elem) Metadata {
	meta := Metadata{
		Labels: []string{},
		Valid:  true,
	}
	keyBalance := map[int]int{}
	var duration time.Duration

	analyzeElements(elems, 1, &meta, keyBalance, &duration)

	for _, balance := range keyBalance {
		if balance > 0 {
			meta.HoldsKeys = true
			break
		}
	}

	meta.EstimatedDuration = formatEstimatedDuration(duration)
	meta.Labels = metadataLabels(meta)

	return meta
}

func analyzeElements(elems []sequencer.Elem, multiplier int, meta *Metadata, keyBalance map[int]int, duration *time.Duration) {
	for idx, elem := range elems {
		switch v := elem.(type) {
		case sequencer.Loop:
			if idx == len(elems)-1 {
				meta.HasLoop = true
			}
		case sequencer.Wait:
			*duration += v.Duration * time.Duration(multiplier)
		case sequencer.Repeat:
			analyzeElements(v.Sequence, multiplier*v.Times, meta, keyBalance, duration)
		case general.KeyDown:
			keyBalance[v.Key] += multiplier
		case general.KeyUp:
			keyBalance[v.Key] -= multiplier
		case general.LeftMouseButtonDown,
			general.LeftMouseButtonUp,
			general.RightMouseButtonDown,
			general.RightMouseButtonUp,
			general.LookupMousePos,
			general.MouseMove,
			general.SetMousePos:
			meta.UsesMouse = true
		}
	}
}

func metadataLabels(meta Metadata) []string {
	labels := []string{}

	if meta.HasLoop {
		labels = append(labels, "Loop")
	}

	if meta.UsesMouse {
		labels = append(labels, "Mouse")
	}

	if meta.HoldsKeys {
		labels = append(labels, "Holds keys")
	}

	if meta.EstimatedDuration != "" {
		labels = append(labels, meta.EstimatedDuration)
	}

	return labels
}

func formatEstimatedDuration(duration time.Duration) string {
	if duration <= 0 {
		return ""
	}

	if duration < time.Second {
		return fmt.Sprintf("~%dms", duration.Milliseconds())
	}

	if duration < time.Minute {
		return fmt.Sprintf("~%ds", int(duration.Round(time.Second).Seconds()))
	}

	return fmt.Sprintf("~%dm", int(duration.Round(time.Minute).Minutes()))
}
