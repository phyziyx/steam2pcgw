package pkg

import (
	"errors"
	"fmt"
)

type CommandHandle func(App, []string) error

type Command struct {
	commandName string
	function    CommandHandle
}

type App struct {
	name        string
	description string
	version     string
	commands    []Command
}

var app App

func (a *App) runCommand(commandName string, parameters []string) error {
	if len(commandName) == 0 {
		return errors.New("an empty command name was specified")
	}

	for _, v := range a.commands {
		if commandName == v.commandName {
			v.function(*a, parameters)
			return nil
		}
	}

	return errors.New("invalid command provided")
}

func Run(commandLineArgs []string) error {
	setup()

	length := len(commandLineArgs)
	if length == 1 {
		return errors.New("no command provided")
	}

	commandName := commandLineArgs[1]
	parameters := make([]string, 0)
	if length > 1 {
		parameters = commandLineArgs[2:]
	}

	return app.runCommand(commandName, parameters)
}

func setup() {
	app.name = APP_NAME
	app.version = VERSION
	app.description = fmt.Sprintf("%s (%s) - A Steam to PCGW article generator (%s)", APP_NAME, VERSION, REPO_LINK)

	app.addCommand("version", VersionCommand)
	app.addCommand("genart", GenerateArticleCommand)
	app.addCommand("gencover", GenerateCoverCommand)
}

func (a *App) InfoPrint() {
	fmt.Printf("%s - %s (%s)", a.name, a.version, a.description)
}

func (a *App) addCommand(commandName string, function CommandHandle) error {
	for _, v := range a.commands {
		if commandName == v.commandName {
			return errors.New("a command with the name '" + commandName + "' already exists")
		}
	}

	a.commands = append(a.commands, Command{
		commandName: commandName,
		function:    function,
	})

	return nil
}

// TODO:
// - Add aliases for command at some point
// aliases     []string
// aliases ...string
// aliasNames := make([]string, 0, len(aliases))
// for idx := 0; idx < len(aliasNames); idx++ {
// 	aliasNames = append(aliasNames, aliases[idx])
// }
// aliases:     aliasNames,
