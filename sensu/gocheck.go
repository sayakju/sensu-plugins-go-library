package sensu

import (
	"fmt"
	"github.com/sensu/sensu-go/types"
	"os"
	"os/exec"
)

type GoCheck struct {
	basePlugin
	validationFunction func(check *types.Check) error
	executeFunction    func(check *types.Check) error
}

func NewGoCheck(config *PluginConfig, options []*PluginConfigOption,
	validationFunction func(check *types.Check) error) *GoCheck {
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
		executeFunction:    executeSensuCheck,
	}

	goCheck.pluginWorkflowFunction = goCheck.goCheckWorkflow
	goCheck.initPlugin()

	return goCheck
}

// Executes the handler's workflow
func (goCheck *GoCheck) goCheckWorkflow(_ []string) (int, error) {
	// Validate input using validateFunction
	err := goCheck.validationFunction(goCheck.sensuCheck)
	if err != nil {
		return 1, fmt.Errorf("error validating input: %s", err)
	}

	// Execute check logic
	err = goCheck.executeFunction(goCheck.sensuCheck)
	if err != nil {
		return 1, fmt.Errorf("error executing handler: %s", err)
	}

	return 0, nil
}

func executeSensuCheck(check *types.Check) error {
	output, _ := exec.Command(check.Command).Output()
	result := string(output)
	fmt.Println(result)
	return nil
}
