package secrets

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const service = "zirc-client"

func Lookup(account string) (string, error) {
	if strings.TrimSpace(account) == "" {
		return "", fmt.Errorf("keyring account must not be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bin, err := exec.LookPath("secret-tool")
	if err != nil {
		return "", fmt.Errorf("Linux Secret Service client secret-tool is unavailable; install libsecret-tools and unlock your desktop keyring")
	}
	cmd := exec.CommandContext(ctx, bin, "lookup", "service", service, "account", account)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not read the requested secret from the desktop keyring")
	}
	value := strings.TrimRight(string(out), "\r\n")
	if value == "" {
		return "", fmt.Errorf("no non-empty secret is stored for this ZIRC keyring account")
	}
	return value, nil
}

func Clear(account string) error {
	if strings.TrimSpace(account) == "" {
		return fmt.Errorf("keyring account must not be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bin, err := exec.LookPath("secret-tool")
	if err != nil {
		return fmt.Errorf("Linux Secret Service client secret-tool is unavailable; install libsecret-tools")
	}
	cmd := exec.CommandContext(ctx, bin, "clear", "service", service, "account", account)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not remove the requested secret from the desktop keyring")
	}
	return nil
}
