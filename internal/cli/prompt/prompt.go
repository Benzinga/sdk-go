package prompt

import (
	"github.com/manifoldco/promptui"
)

const mask = '*'

func ForString(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: existsValidator(),
	}

	return prompt.Run()
}

func ForStringMasked(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: existsValidator(),
		Mask:     mask,
	}

	return prompt.Run()
}
