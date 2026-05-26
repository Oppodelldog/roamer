package script_test

import (
	"reflect"
	"testing"

	"github.com/Oppodelldog/roamer/internal/script"
)

func TestAnalyzeMetadata(t *testing.T) {
	meta, err := script.Analyze("KD A;LD;W 60ms;LU;W 1s;L")
	if err != nil {
		t.Fatalf("did not expect an error, got %v", err)
	}

	wantLabels := []string{"Loop", "Mouse", "Holds keys", "~1s"}
	if !reflect.DeepEqual(wantLabels, meta.Labels) {
		t.Fatalf("labels did not match:\ngot : %#v\nwant: %#v", meta.Labels, wantLabels)
	}

	if !meta.HasLoop || !meta.UsesMouse || !meta.HoldsKeys || !meta.Valid {
		t.Fatalf("metadata flags were not set as expected: %#v", meta)
	}
}

func TestAnalyzeInvalidSequence(t *testing.T) {
	meta, err := script.Analyze("KU")
	if err == nil {
		t.Fatal("expected an error")
	}

	if meta.Valid {
		t.Fatalf("expected invalid metadata, got %#v", meta)
	}

	if meta.Error == "" {
		t.Fatal("expected metadata error text")
	}
}

func TestAnalyzeBalancedKeysDoNotHold(t *testing.T) {
	meta, err := script.Analyze("KD A;W 50ms;KU A")
	if err != nil {
		t.Fatalf("did not expect an error, got %v", err)
	}

	if meta.HoldsKeys {
		t.Fatalf("expected balanced key down/up not to hold keys: %#v", meta)
	}

	wantLabels := []string{"~50ms"}
	if !reflect.DeepEqual(wantLabels, meta.Labels) {
		t.Fatalf("labels did not match:\ngot : %#v\nwant: %#v", meta.Labels, wantLabels)
	}
}

func TestAnalyzeRepeatMetadata(t *testing.T) {
	meta, err := script.Analyze("R 3 [LD;W 400ms;LU]")
	if err != nil {
		t.Fatalf("did not expect an error, got %v", err)
	}

	if !meta.UsesMouse {
		t.Fatalf("expected repeat metadata to include mouse usage: %#v", meta)
	}

	if meta.EstimatedDuration != "~1s" {
		t.Fatalf("expected repeated duration to be ~1s, got %q", meta.EstimatedDuration)
	}

	wantLabels := []string{"Mouse", "~1s"}
	if !reflect.DeepEqual(wantLabels, meta.Labels) {
		t.Fatalf("labels did not match:\ngot : %#v\nwant: %#v", meta.Labels, wantLabels)
	}
}
