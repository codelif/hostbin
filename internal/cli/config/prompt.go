package config

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

type Prompter interface {
	Init(defaultServerURL string) (string, string, error)
	Value(key, current string) (string, error)
}

type HuhPrompter struct{}

func (HuhPrompter) Init(defaultServerURL string) (string, string, error) {
	serverURL := defaultServerURL
	authKey := ""

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Admin server URL").
				Description("Press Enter to keep the default.").
				Value(&serverURL).
				Validate(validateServerURLInput),
			huh.NewInput().
				Title("Auth key").
				EchoMode(huh.EchoModePassword).
				Value(&authKey).
				Validate(validateAuthKeyInput),
		),
	)

	if err := form.Run(); err != nil {
		return "", "", err
	}

	return strings.TrimSpace(serverURL), strings.TrimSpace(authKey), nil
}

func (HuhPrompter) Value(key, current string) (string, error) {
	value := current
	input := huh.NewInput().
		Title(promptTitle(key)).
		Value(&value)

	switch key {
	case "server_url":
		input = input.Description("Press Enter to keep the current value.").Validate(validateServerURLInput)
	case "auth_key":
		input = input.EchoMode(huh.EchoModePassword).Validate(validateAuthKeyInput)
	}

	if err := huh.NewForm(huh.NewGroup(input)).Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(value), nil
}

func promptTitle(key string) string {
	switch key {
	case "server_url":
		return "Admin server URL"
	case "auth_key":
		return "Auth key"
	case "timeout":
		return "Timeout"
	case "editor":
		return "Editor"
	case "color":
		return "Color mode"
	default:
		return fmt.Sprintf("Value for %s", key)
	}
}

func validateServerURLInput(value string) error {
	return File{ServerURL: value, AuthKey: strings.Repeat("x", 32)}.Validate()
}

func validateAuthKeyInput(value string) error {
	return File{ServerURL: DefaultServerURL, AuthKey: value}.Validate()
}
