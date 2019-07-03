package sensu

import (
	"fmt"
	"github.com/sensu/sensu-go/types"
	"os"
)

type GoCheck struct {
	basePlugin
	validationFunction func(check *types.Check, entity *types.Entity) error
	executeFunction    func(event *types.Check, entity *types.Entity) (int, error)
}

func NewGoCheck(config *PluginConfig, options []*PluginConfigOption,
	validationFunction func(check *types.Check, entity *types.Entity) error, executeFunction func(event *types.Check, entity *types.Entity) (int, error)) *GoCheck {
	goCheck := &GoCheck{
		basePlugin: basePlugin{
			config:          config,
			options:         options,
			sensuCheck:      nil,
			sensuEntity:     nil,
			checkReader:     os.Stdin,
			entityReader:    os.Stdin,
			readCheck:       true,
			readEntity:      true,
			checkMandatory:  true,
			entityMandatory: false,
			errorExitStatus: -1,
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
	err := goCheck.validationFunction(goCheck.sensuCheck, goCheck.sensuEntity)
	if err != nil {
		return -1, fmt.Errorf("error validating input: %s", err)
	}

	// Execute check logic using executeFunction
	return goCheck.executeFunction(goCheck.sensuCheck, goCheck.sensuEntity)
}
