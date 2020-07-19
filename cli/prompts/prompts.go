package prompts

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"sopr/types"
	"sort"
	"strings"
)

// BranchNamePrompt is a cli component for getting the branch name from a user
func BranchNamePrompt() string {
	prompt := promptui.Prompt{
		Label: "Branch Name",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println(fmt.Sprintf("Prompt failed %v", err))
		os.Exit(1)
	}

	return result
}

// SelectOption struct
type SelectOption struct {
	Name         string
	Selected     bool
	IsSelectable bool
}

var doneOption SelectOption = SelectOption{
	Name:         "Done",
	Selected:     false,
	IsSelectable: false,
}

// RepoSelectPrompt determines which repositories to action
func RepoSelectPrompt(projectList []types.ProjectConfig) []types.ProjectConfig {
	var options []SelectOption
	var results []types.ProjectConfig
	lastPosition := 0

	for _, project := range projectList {
		option := SelectOption{
			Name:         project.Name,
			Selected:     false,
			IsSelectable: true,
		}

		options = append(options, option)
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].Name < options[j].Name
	})

	options = append(options, doneOption)

	index := -1
	for index < len(options)-1 {
		var err error
		var result string
		var newOptions []SelectOption

		// if index >= 0 {
		// 	options = append(options[:index], options[index+1:]...)
		// }

		projectSelectPrompt := promptui.Select{
			Items:        options,
			CursorPos:    lastPosition,
			Size:         20,
			HideSelected: true,
			IsVimMode:    false,
			Stdout:       &bellSkipper{},
			Templates: &promptui.SelectTemplates{
				Label:    "Select Projects",
				Active:   "{{if .IsSelectable}} {{if .Selected }} [x] {{else}} [ ] {{end}} {{end}} {{.Name | cyan}}",
				Inactive: "{{if .IsSelectable}} {{if .Selected}} [x] {{else}} [ ] {{end}} {{end}} {{.Name}}",
			},
		}

		index, result, err = projectSelectPrompt.Run()
		if err != nil {
			fmt.Println(fmt.Sprintf("Error: failed to open repo select list: %s", err))
			os.Exit(1)
		}

		lastPosition = index + 1

		name := strings.TrimLeft(strings.Split(result, " ")[0], "{")

		if name == doneOption.Name {
			break
		}

		for _, option := range options {
			if option.Name == name {
				option.Selected = !option.Selected
			}

			newOptions = append(newOptions, option)
		}

		options = newOptions
		projectSelectPrompt.Items = options
	}

	for _, option := range options {
		if option.Selected {
			for _, project := range projectList {
				if project.Name == option.Name {
					results = append(results, project)
					break
				}
			}
		}
	}

	return results
}

// bellSkipper implements an io.WriteCloser that skips the terminal bell
// character (ASCII code 7), and writes the rest to os.Stderr. It is used to
// replace readline.Stdout, that is the package used by promptui to display the
// prompts.
//
// This is a workaround for the bell issue documented in
// https://github.com/manifoldco/promptui/issues/49.
// thanks to https://https://github.com/mroth
type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}
