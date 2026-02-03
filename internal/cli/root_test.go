package cli

import "testing"

func TestRun_NoArgs(t *testing.T) {
    if err := Run([]string{}); err == nil {
        t.Fatalf("expected error for missing command")
    }
}

func TestRun_Help(t *testing.T) {
    if err := Run([]string{"help"}); err != nil {
        t.Fatalf("expected nil error for help, got %v", err)
    }
}

func TestRun_UnknownCommand(t *testing.T) {
    if err := Run([]string{"nope"}); err == nil {
        t.Fatalf("expected error for unknown command")
    }
}
