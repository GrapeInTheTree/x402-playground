package quiz

import (
	"testing"
)

func TestParseGoTestJSON_PassAndFail(t *testing.T) {
	output := `{"Action":"run","Test":"TestFoo"}
{"Action":"pass","Test":"TestFoo"}
{"Action":"run","Test":"TestBar"}
{"Action":"fail","Test":"TestBar"}
{"Action":"fail","Package":"example"}
`
	passed, failed := parseGoTestJSON(output)
	if passed != 1 {
		t.Errorf("expected 1 passed, got %d", passed)
	}
	if failed != 1 {
		t.Errorf("expected 1 failed, got %d", failed)
	}
}

func TestParseGoTestJSON_AllPass(t *testing.T) {
	output := `{"Action":"run","Test":"TestA"}
{"Action":"pass","Test":"TestA"}
{"Action":"run","Test":"TestB"}
{"Action":"pass","Test":"TestB"}
{"Action":"pass","Package":"example"}
`
	passed, failed := parseGoTestJSON(output)
	if passed != 2 {
		t.Errorf("expected 2 passed, got %d", passed)
	}
	if failed != 0 {
		t.Errorf("expected 0 failed, got %d", failed)
	}
}

func TestParseGoTestJSON_Empty(t *testing.T) {
	passed, failed := parseGoTestJSON("")
	if passed != 0 || failed != 0 {
		t.Errorf("expected 0/0, got %d/%d", passed, failed)
	}
}

func TestParseGoTestJSON_InvalidJSON(t *testing.T) {
	output := "not json\nstill not json\n"
	passed, failed := parseGoTestJSON(output)
	if passed != 0 || failed != 0 {
		t.Errorf("expected 0/0 for invalid JSON, got %d/%d", passed, failed)
	}
}

func TestParseGoTestJSON_PackageLevelEventsIgnored(t *testing.T) {
	// Package-level events (no Test field) should not be counted
	output := `{"Action":"pass","Package":"example"}
{"Action":"fail","Package":"example2"}
`
	passed, failed := parseGoTestJSON(output)
	if passed != 0 || failed != 0 {
		t.Errorf("package events should be ignored, got passed=%d failed=%d", passed, failed)
	}
}
