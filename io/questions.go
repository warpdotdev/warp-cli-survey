package io

import (
	"os"

	"github.com/zachlloyd/denver-survey-client/history"
	"github.com/zachlloyd/denver-survey-client/shell"
)

// Questions returns a list of all questions in the survey
func Questions() []Question {
	return []Question{
		{ID: "company", Text: "What company do you work at?", Type: FreeForm},
		{ID: "years_of_experience", Text: "How many years experience do you have?", Type: FreeForm},
		{ID: "role", Text: "Which of these best describes your role?", Type: MultipleChoice,
			Values: []string{
				"Software Engineer",
				"DevOps Engineer / SRE",
				"Engineering Manager",
				"Engineering Leadership (Director / VP / CTO)",
				"Test Engineer",
				"QA"}},
		{ID: "platform_develop_on", Text: "What platform do you (or your team) develop on?", Type: MultipleChoice,
			Values: []string{
				"Mac",
				"Linux",
				"Windows"}},
		{ID: "platform_develop_for", Text: "What platforms do you primarily develop for?",
			Type: MultipleChoice, MultiSelect: true,
			Values: []string{
				"Linux / Unix (server / backend)",
				"Web / Frontend",
				"iOS",
				"AndroId",
				"Windows",
				"Mac"}},
		{ID: "shell_type", Text: "Checking your system...looks your default shell is:", Type: FreeForm,
			SuggestedAnswerFn: func() string {
				return os.ExpandEnv("$SHELL")
			}},
		{ID: "shell_history",
			Text: "Is it OK to upload a redacted version of your shell history file (zsh / bash only)?\nThis will be used only for aggregate analysis of how developers use the CLI.",
			Type: File,
			Values: []string{
				"Yes (we will show you a preview first)",
				"No"},
			GetShellHistoryFn: history.GetRedactedShellHistory,
			ShouldShowFn: func(responsesSoFar map[string]*Answer) bool {
				shellType := shell.GetShellType(responsesSoFar["shell_type"].Text)
				return shellType == shell.Bash || shellType == shell.Zsh
			}},
		{ID: "terminal_type", Text: "What terminal do you typically use?", Type: MultipleChoice,
			Values: []string{
				"Mac Terminal",
				"iTerm",
				"Hyper",
				"Windows Command Line",
				"PowerShell",
				"A linux terminal (e.g. Gnome)"}},
		{ID: "other_tools", Text: "Which of the following applications / platforms / configurations do you use to improve your experience in the command-line?",
			Type: MultipleChoice, MultiSelect: true,
			Values: []string{
				"tmux",
				"screen",
				"ohmyzsh",
				"a dotfiles repo",
				"None"}},
		{ID: "cli_experience", Text: "What would you most like to see improved in the command-line experience?",
			Type: MultipleChoice, MultiSelect: true,
			Values: []string{
				"Better command autocomplete",
				"An easier way of saving and sharing work (e.g. something like a Jupyter notebook for the terminal)",
				"A browser based terminal that attaches to cloud machines",
				"An easier way of setting up and maintaining developer environments",
				"Real-time collaboration (e.g. share terminal input and output with team members)",
				"Collaborative terminal workflows like \"command-reviews\" (similar to a code review but for commands)",
				"Better session and window management (e.g. built in Tmux functionality)"}},
		{ID: "email", Text: "Interested in following Project Denver / helping test the product?  Please let us know your email.", Type: FreeForm},
	}
}
