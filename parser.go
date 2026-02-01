package main

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	prefixes = map[[2]byte]string{
		{'N', 'o'}: "Notification",
		{'P', 'r'}: "PreToolUse",
		{'P', 'o'}: "PostToolUse",
		{'B', 'a'}: "Batch",
		{'E', 'd'}: "Edit",
		{'R', 'e'}: "Read",
		{'W', 'r'}: "Write",
		{'G', 'l'}: "Glob",
		{'G', 'r'}: "Grep",
		{'T', 'a'}: "Task",
		{'W', 'e'}: "Web",
	}
)

func parse(buf []byte) (bool, string, error) {
	hook, err := extract(buf, []byte(`"hook_event_name":"`))
	if err != nil {
		return false, "", err
	}

	// Only the "idle_prompt" notification is supported. For hook type
	// "Notification", we return "Idle" immediately without inspecting the
	// payload.
	if hook == "Notification" {
		return false, "Idle", nil
	}

	word, err := extract(buf, []byte(`"tool_name":"`))
	if err != nil {
		return false, "", err
	}

	// "PostToolUse" runs immediately after a tool completes successfully.
	// A hook value of "PostToolUse" indicates the tool has finished
	// executing.
	return hook == "PostToolUse", word, nil
}

func extract(buf, target []byte) (string, error) {
	i := bytes.Index(buf, target)
	if i < 0 || len(buf)-(i+2) < len(target) {
		return "", errors.New("index out of bounds")
	}

	var k [2]byte
	copy(k[:], buf[i+len(target):i+len(target)+2])

	word, ok := prefixes[k]
	if !ok {
		return "", fmt.Errorf("word not found for prefix: %q", k)
	}

	return word, nil
}
