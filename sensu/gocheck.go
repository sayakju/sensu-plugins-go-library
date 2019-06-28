package sensu

import (
	"context"
	"fmt"
	"github.com/sensu/sensu-go/command"
)

type GoCheck struct {
	basePlugin
	validationFunction func() error
	command            string
}

func NewGoCheck(config *PluginConfig, options []*PluginConfigOption,
	validationFunction func() error, command string) *GoCheck {
	goCheck := &GoCheck{
		basePlugin: basePlugin{
			config:          config,
			options:         options,
			readEvent:       false,
			errorExitStatus: 1,
		},
		validationFunction: validationFunction,
		command:            command,
	}

	goCheck.pluginWorkflowFunction = goCheck.goCheckWorkflow
	goCheck.initPlugin()

	return goCheck
}

// Executes the handler's workflow
func (goCheck *GoCheck) goCheckWorkflow(_ []string) (int, error) {
	// Validate input using validateFunction
	err := goCheck.validationFunction()
	if err != nil {
		return 1, fmt.Errorf("error validating input: %s", err)
	}

	ex := command.ExecutionRequest{
		Command: goCheck.command,
	}

	checkExec, err := command.NewExecutor().Execute(context.Background(), ex)
	if err != nil {
		return 2, err
	}

	fmt.Println(checkExec.Output)

	return 0, nil
}
