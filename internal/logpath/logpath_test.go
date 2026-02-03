package logpath_test

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/BarrettBr/RWND/internal/logpath"
)

func TestResolveRecordPath_DirectoryCreatesNumberedFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "001_dummy.jsonl"), []byte("x"), 0644); err != nil {
		t.Fatalf("write existing: %v", err)
	}

	got, err := logpath.ResolveRecordPath(dir, ":8080", nil)
	if err != nil {
		t.Fatalf("ResolveRecordPath: %v", err)
	}

	pattern := regexp.MustCompile(`[/\\]002_\d{8}T\d{6}Z_listen-8080\.jsonl$`)
	if !pattern.MatchString(got) {
		t.Fatalf("unexpected path: %s", got)
	}
}

func TestResolveRecordPath_FilePassthrough(t *testing.T) {
	dir := t.TempDir()
	want := filepath.Join(dir, "custom.jsonl")

	got, err := logpath.ResolveRecordPath(want, ":8080", nil)
	if err != nil {
		t.Fatalf("ResolveRecordPath: %v", err)
	}
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolveRecordPath_UsesTargetHostInFilename(t *testing.T) {
	dir := t.TempDir()
	target, err := url.Parse("http://example.com:3000")
	if err != nil {
		t.Fatalf("parse target: %v", err)
	}

	got, err := logpath.ResolveRecordPath(dir, ":8080", target)
	if err != nil {
		t.Fatalf("ResolveRecordPath: %v", err)
	}

	pattern := regexp.MustCompile(`[/\\]\d{3}_\d{8}T\d{6}Z_listen-8080_target-example-com-3000\.jsonl$`)
	if !pattern.MatchString(got) {
		t.Fatalf("unexpected path: %s", got)
	}
}

func TestResolveReplayPath_PicksLatestNumberedFile(t *testing.T) {
	dir := t.TempDir()
	files := []string{
		"002_first.jsonl",
		"010_second.jsonl",
		"007_third.jsonl",
	}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	got, err := logpath.ResolveReplayPath(dir)
	if err != nil {
		t.Fatalf("ResolveReplayPath: %v", err)
	}
	want := filepath.Join(dir, "010_second.jsonl")
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolveReplayPath_FilePassthrough(t *testing.T) {
	dir := t.TempDir()
	want := filepath.Join(dir, "explicit.jsonl")

	got, err := logpath.ResolveReplayPath(want)
	if err != nil {
		t.Fatalf("ResolveReplayPath: %v", err)
	}
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestResolveReplayPath_NoLogsReturnsError(t *testing.T) {
	dir := t.TempDir()
	if _, err := logpath.ResolveReplayPath(dir); err == nil {
		t.Fatalf("expected error for empty log directory")
	}
}
