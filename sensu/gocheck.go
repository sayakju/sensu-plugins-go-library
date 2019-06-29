package sensu

import (
	"fmt"
	"github.com/sensu/sensu-go/types"
	"os"
)

type GoCheck struct {
	basePlugin
	validationFunction func(check *types.Check) error
	executeFunction    func(event *types.Check) error
}

func NewGoCheck(config *PluginConfig, options []*PluginConfigOption,
	validationFunction func(check *types.Check) error, executeFunction func(event *types.Check) error) *GoCheck {
	goCheck := &GoCheck{
		basePlugin: basePlugin{
			config:          config,
			options:         options,
			sensuCheck:      nil,
			checkReader:     os.Stdin,
			readCheck:       true,
			checkMandatory:  true,
			errorExitStatus: 1,
		},
		validationFunction: validationFunction,
		executeFunction:    executeFunction,
	}

	goCheck.pluginWorkflowFunction = goCheck.goCheckWorkflow
	goCheck.initPlugin()

	return goCheck
}

// Executes the check's workflow
func (goCheck *GoCheck) goCheckWorkflow(_ []string) (int, error) {
	// Validate input using validateFunction
	err := goCheck.validationFunction(goCheck.sensuCheck)
	if err != nil {
		return 1, fmt.Errorf("error validating input: %s", err)
	}

	// Execute check logic using executeFunction
	err = goCheck.executeFunction(goCheck.sensuCheck)
	if err != nil {
		return 1, fmt.Errorf("error executing check: %s", err)
	}

	return 0, nil
}
