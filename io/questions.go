package io

import (
	"os"

	"github.com/warpdotdev/denver-survey-client/history"
	"github.com/warpdotdev/denver-survey-client/shell"
)

// QuestionID is a human readable unique question id.
type QuestionID string

const (
	company            QuestionID = "company"
	yearsOfExperience             = "years_of_experience"
	role                          = "role"
	platformDevelopOn             = "platform_develop_on"
	platformDevelopFor            = "platform_develop_for"

	shellType       = "shell_type"
	whyThatShell    = "why_that_shell"
	shellHistory    = "shell_history"
	terminalType    = "terminal_type"
	whyThatTerminal = "why_that_terminal"

	codeEditor         = "code_editor"
	levelOfExpertise   = "level_of_expertise"
	frequencyOfUse     = "frequency_of_use"
	numTerminalWindows = "num_terminal_windows"
	numGithubRepos     = "num_github_repos"
	otherTools         = "other_tools"

	mainReasonForUsingCLI = "main_reason_for_using_cli"
	wantImproved          = "want_improved"
	biggestPainPoint      = "biggest_pain_point"
	payFor                = "pay_for"

	// Email is the email address question
	Email        = "email"
	okToReachOut = "ok_to_reach_out"
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
		m[whyThatTerminal],
		m[levelOfExpertise],
		m[frequencyOfUse],
		m[numTerminalWindows],
		m[numGithubRepos],
		m[shellType],
		m[whyThatShell],
		m[shellHistory],
		m[codeEditor],
		m[otherTools],

		m[mainReasonForUsingCLI],
		m[biggestPainPoint],
		m[payFor],
		m[wantImproved],
		m[Email],
		m[okToReachOut],
	}
}

var questionMap = map[QuestionID]Question{
	company:           {ID: company, Text: "What company do you work at?", Type: FreeForm},
	yearsOfExperience: {ID: yearsOfExperience, Text: "How many years experience using the CLI do you have?", Type: FreeForm},
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
			"Android",
			"Windows",
			"Mac"}},

	shellType: {ID: shellType, Text: "Checking your system...looks your default shell is:", Type: FreeForm,
		SuggestedAnswerFn: func() string {
			return os.ExpandEnv("$SHELL")
		}},
	whyThatShell: {ID: whyThatShell, Text: "Anything in particular that made you pick that shell?",
		Type: FreeForm},
	shellHistory: {ID: shellHistory,
		Text: `Can we take a look at your shell history to get a better sense of how you use the CLI?
> We will strip out all sensitive information first.

** Is this safe? Yes, the data is sanitized and you can see exactly what we will store beforehand.
** But we get that this could be scary, so it's totally up to you if you share (although it would be helpful!)`,
		Type: File,
		Values: []string{
			"Yes (shows a preview before submitting)",
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
	whyThatTerminal: {ID: whyThatTerminal, Text: "Anything in particular that made you pick that terminal?",
		Type: FreeForm},

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
			"One per project I'm working on",
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
	mainReasonForUsingCLI: {ID: mainReasonForUsingCLI, Text: "What's the main reason you use the CLI?",
		Type: FreeForm},
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
	payFor: {ID: payFor, Text: "Is there an improvement to the command-line you would pay $10 a month for?  If so, please tell us about it.",
		Type: FreeForm},
	Email: {ID: Email, Text: "What's your email? Will only be used to send you survey results. [enter blank to skip]",
		Type: FreeForm, Skippable: true},
	okToReachOut: {ID: okToReachOut, Text: "We'd love to reach out and pick your brain on the product - is that OK?",
		Type: YesNo, HasDefault: true},
}
