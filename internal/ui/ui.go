// Package ui provides minimal terminal output helpers (stdlib only).
package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Header(title string) {
	fmt.Printf("\n== %s ==\n", title)
}

func Ok(msg string)      { fmt.Printf("[ok]      %s\n", msg) }
func Found(msg string)   { fmt.Printf("[found]   %s\n", msg) }
func Missing(msg string) { fmt.Printf("[missing] %s\n", msg) }
func Warn(msg string)    { fmt.Printf("[warn]    %s\n", msg) }
func Info(msg string)    { fmt.Printf("[info]    %s\n", msg) }
func Step(n int, msg string) {
	fmt.Printf("\n  Step %d: %s\n", n, msg)
}

// Confirm prompts for a yes/no answer. Default is no.
func Confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}

// Ask prompts for a free-text answer.
func Ask(prompt string) string {
	fmt.Printf("%s: ", prompt)
	r := bufio.NewReader(os.Stdin)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func Fatal(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}
