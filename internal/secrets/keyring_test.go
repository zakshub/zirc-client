package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLookupUsesSecretToolWithoutPrintingSecret(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "secret-tool")
	script := "#!/bin/sh\nprintf 'test-secret\\n'\n"
	if err := os.WriteFile(bin, []byte(script), 0700); err != nil { t.Fatal(err) }
	t.Setenv("PATH", dir)
	got, err := Lookup("account-name")
	if err != nil { t.Fatal(err) }
	if got != "test-secret" { t.Fatalf("unexpected keyring value %q", got) }
}

func TestLookupRejectsEmptyAccount(t *testing.T) {
	if _, err := Lookup(" "); err == nil { t.Fatal("expected empty account to be rejected") }
}
