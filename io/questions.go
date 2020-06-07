package io

import (
	"os"

	"github.com/zachlloyd/denver-survey-client/history"
	"github.com/zachlloyd/denver-survey-client/shell"
)

// QuestionID is a human readable unique question id.
type QuestionID string

const (
	company            QuestionID = "company"
	yearsOfExperience             = "years_of_experience"
	role                          = "role"
	platformDevelopOn             = "platform_develop_on"
	platformDevelopFor            = "platform_develop_for"

	shellType    = "shell_type"
	shellHistory = "shell_history"
	terminalType = "terminal_type"

	codeEditor         = "code_editor"
	levelOfExpertise   = "level_of_expertise"
	frequencyOfUse     = "frequency_of_use"
	numTerminalWindows = "num_terminal_windows"
	numGithubRepos     = "num_github_repos"
	otherTools         = "other_tools"

	wantImproved     = "want_improved"
	biggestPainPoint = "biggest_pain_point"
	payFor           = "pay_for"
	email            = "email"
	okToReachOut     = "ok_to_reach_out"
)

// Questions returns a list of all questions in the survey
func Questions() []Question {
	m := questionMap
	return []Question{
		m[company],
		m[yearsOfExperience],
		m[role],
		m[platformDevelopOn],
		m[platformDevelopFor],

		m[terminalType],
		m[levelOfExpertise],
		m[frequencyOfUse],
		m[numTerminalWindows],
		m[numGithubRepos],
		m[shellType],
		m[shellHistory],
		m[codeEditor],
		m[otherTools],

		m[biggestPainPoint],
		m[payFor],
		m[email],
		m[okToReachOut],
	}
}

var questionMap = map[QuestionID]Question{
	company:           {ID: company, Text: "What company do you work at?", Type: FreeForm},
	yearsOfExperience: {ID: yearsOfExperience, Text: "How many years experience do you have?", Type: FreeForm},
	role: {ID: role, Text: "Which of these best describes your role?", Type: MultipleChoice, ShowOther: true,
		Values: []string{
			"Software Engineer",
			"DevOps Engineer / SRE",
			"Engineering Manager",
			"Engineering Leadership (Director / VP / CTO)",
			"Test Engineer",
			"QA"}},
	platformDevelopOn: {ID: platformDevelopOn, Text: "What type of computer do you write code on?",
		Type: MultipleChoice, ShowOther: true,
		Values: []string{
			"Mac",
			"Linux",
			"Windows"}},
	platformDevelopFor: {ID: platformDevelopFor, Text: "What platforms do you primarily develop for?",
		Type: MultipleChoice, MultiSelect: true, ShowOther: true,
		Values: []string{
			"Linux / Unix (server / backend)",
			"Web / Frontend",
			"iOS",
			"AndroId",
			"Windows",
			"Mac"}},

	shellType: {ID: shellType, Text: "Checking your system...looks your default shell is:", Type: FreeForm,
		SuggestedAnswerFn: func() string {
			return os.ExpandEnv("$SHELL")
		}},
	shellHistory: {ID: shellHistory,
		Text: "Can we take a look at your shell history to get a better sense of how you use the CLI?\n> We will strip out all sensitive information first.",
		Type: File,
		Values: []string{
			"Yes (we will show you a preview first)",
			"No"},
		GetShellHistoryFn: history.GetRedactedShellHistory,
		ShouldShowFn: func(responsesSoFar map[QuestionID]*Answer) bool {
			shellType := shell.GetShellType(responsesSoFar["shell_type"].Text)
			return shellType == shell.Bash || shellType == shell.Zsh
		}},
	terminalType: {ID: terminalType, Text: "What terminal do you typically use?",
		Type: MultipleChoice, MultiSelect: true, ShowOther: true,
		Values: []string{
			"Mac Terminal",
			"The terminal that is embedded in my IDE (e.g. VSCode)",
			"iTerm",
			"Hyper",
			"Windows Command Line",
			"PowerShell",
			"A linux terminal (e.g. Gnome)"}},

	levelOfExpertise: {ID: levelOfExpertise, Text: "How experienced of a command-line user are you?", Type: MultipleChoice,
		Values: []string{
			"Novice (only know the basics like cd, ls, pwd...)",
			"Competent (can use grep, find, chmod)",
			"Advanced (have written scripts, use pipes and xargs)",
			"Expert (use the CLI like a ninja)"}},

	frequencyOfUse: {ID: frequencyOfUse, Text: "How often do you use the command-line?", Type: MultipleChoice,
		Values: []string{
			"Infrequently (not every day)",
			"A few times a day (on average)",
			"A few times an hour (on average)",
			"Constantly (it's always open and I'm using it as one of my main tools)"}},

	numTerminalWindows: {ID: numTerminalWindows, Text: "How many terminal windows or tabs do you usually have open?", Type: MultipleChoice,
		Values: []string{
			"Zero",
			"One total",
			"One per project I'm woking on",
			"Multiple per project",
			"The one embedded in my IDE"}},

	numGithubRepos: {ID: numGithubRepos, Text: "How many different git repos are you typically working with?", Type: MultipleChoice,
		Values: []string{
			"Zero",
			"One",
			"2-4",
			"5+"}},

	codeEditor: {ID: codeEditor, Text: "What code editors or IDEs do you typically use?",
		Type: MultipleChoice, ShowOther: true, MultiSelect: true,
		Values: []string{
			"vim",
			"Emacs",
			"VSCode",
			"JetBrains product (e.g. IntelliJ, WebStorm or PyCharm)",
			"Atom",
			"XCode",
			"Visual Studio",
			"Android Studio"}},
	otherTools: {ID: otherTools, Text: "Which of the following applications / platforms / configurations do you use to improve your experience in the command-line?",
		Type: MultipleChoice, MultiSelect: true, ShowOther: true,
		Values: []string{
			"tmux",
			"screen",
			"ohmyzsh",
			"a dotfiles repo",
			"None"}},
	wantImproved: {ID: wantImproved, Text: "What would you most like to see improved in the command-line experience?",
		Type: MultipleChoice, MultiSelect: true, ShowOther: true,
		Values: []string{
			"Better command autocomplete",
			"An easier way of saving and sharing work (e.g. something like a Jupyter notebook for the terminal)",
			"A browser based terminal that attaches to cloud machines",
			"An easier way of setting up and maintaining developer environments",
			"Real-time collaboration (e.g. share terminal input and output with team members)",
			"Collaborative terminal workflows like \"command-reviews\" (similar to a code review but for commands)",
			"Better session and window management (e.g. built in Tmux functionality)"}},
	biggestPainPoint: {ID: biggestPainPoint, Text: "What's your biggest pain point working in the command-line?",
		Type: FreeForm},
	payFor: {ID: payFor, Text: "Is there an improvement to the command-line you would pay $10 / mo for?  If so, please tell us about it.",
		Type: FreeForm},
	email: {ID: email, Text: "Interested in seeing survey results and following our progress?  Please let us know your email.",
		Type: FreeForm, Skippable: true},
	okToReachOut: {ID: okToReachOut, Text: "We'd love to reach out and pick your brain on the product - is that OK?", Type: YesNo},
}
