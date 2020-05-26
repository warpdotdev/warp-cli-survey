package question

import (
	"os"

	"github.com/zachlloyd/denver-survey-client/history"
	"github.com/zachlloyd/denver-survey-client/shell"
)

func Questions() []Question {
	return []Question{
		Question{Id: "company", Question: "What company do you work at?", QuestionType: FreeForm},
		Question{Id: "years_of_experience", Question: "How many years experience do you have?", QuestionType: FreeForm},
		Question{Id: "role", Question: "Which of these best describes your role?", QuestionType: MultipleChoice,
			Values: []string{
				"Software Engineer",
				"DevOps Engineer / SRE",
				"Engineering Manager",
				"Engineering Leadership (Director / VP / CTO)",
				"Test Engineer",
				"QA"}},
		Question{Id: "platform_develop_on", Question: "What platform do you (or your team) develop on?", QuestionType: MultipleChoice,
			Values: []string{
				"Mac",
				"Linux",
				"Windows"}},
		Question{Id: "platform_develop_for", Question: "What platforms do you primarily develop for?",
			QuestionType: MultipleChoice, MultiSelect: true,
			Values: []string{
				"Linux / Unix (server / backend)",
				"Web / Frontend",
				"iOS",
				"AndroId",
				"Windows",
				"Mac"}},
		Question{Id: "shell_type", Question: "Checking your system...looks your default shell is:", QuestionType: FreeForm,
			SuggestedAnswerFn: func() string {
				return os.ExpandEnv("$SHELL")
			}},
		Question{Id: "shell_history",
			Question:     "Is it OK to upload a redacted version of your shell history file (zsh / bash only)?\nThis will be used only for aggregate analysis of how developers use the CLI.",
			QuestionType: File,
			Values: []string{
				"Yes (we will show you a preview first)",
				"No"},
			GetShellHistoryFn: history.GetRedactedShellHistory,
			ShouldShowFn: func(responsesSoFar map[string]*Response) bool {
				shellType := shell.GetShellType(responsesSoFar["shell_type"].Answer)
				return shellType == shell.Bash || shellType == shell.Zsh
			}},
		Question{Id: "terminal_type", Question: "What terminal do you typically use?", QuestionType: MultipleChoice,
			Values: []string{
				"Mac Terminal",
				"iTerm",
				"Hyper",
				"Windows Command Line",
				"PowerShell",
				"A linux terminal (e.g. Gnome)"}},
		Question{Id: "other_tools", Question: "Which of the following applications / platforms / configurations do you use to improve your experience in the command-line?",
			QuestionType: MultipleChoice, MultiSelect: true,
			Values: []string{
				"tmux",
				"screen",
				"ohmyzsh",
				"a dotfiles repo",
				"None"}},
		Question{Id: "cli_experience", Question: "What would you most like to see improved in the command-line experience?",
			QuestionType: MultipleChoice, MultiSelect: true,
			Values: []string{
				"Better command autocomplete",
				"An easier way of saving and sharing work (e.g. something like a Jupyter notebook for the terminal)",
				"A browser based terminal that attaches to cloud machines",
				"An easier way of setting up and maintaining developer environments",
				"Real-time collaboration (e.g. share terminal input and output with team members)",
				"Collaborative terminal workflows like \"command-reviews\" (similar to a code review but for commands)",
				"Better session and window management (e.g. built in Tmux functionality)"}},
		Question{Id: "email", Question: "Interested in following Project Denver / helping test the product?  Please let us know your email.", QuestionType: FreeForm},
	}
}
