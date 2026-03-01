package main

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

const pandocTimeout = 5 * time.Second

func renderMarkdownWithPandoc(markdown string) (string, error) {
	if markdown == "" {
		return "", errors.New("markdown is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), pandocTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pandoc", "-f", "markdown", "-t", "html")
	cmd.Stdin = bytes.NewBufferString(markdown)
	out, err := cmd.Output()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", err
	}

	return string(out), nil
}
