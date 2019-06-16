package sensu

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sensu/sensu-go/command"
	"github.com/sensu/sensu-go/types"
	"os"
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
		return 1, fmt.Errorf("error executing check: %s", err)
	}

	return 0, nil
}

func executeSensuCheck(check *types.Check) error {
	// Inject the dependencies into PATH, LD_LIBRARY_PATH & CPATH so that they
	// are availabe when when the command is executed.
	ex := command.ExecutionRequest{
		Command: check.Command,
		Name:    check.Name,
	}

	// If stdin is true, add JSON event data to command execution.
	if check.Stdin {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		ex.Input = string(input)
	}

	checkExec, _ := command.NewExecutor().Execute(context.Background(), ex)
	output := checkExec.Output

	duration := checkExec.Duration
	status := uint32(checkExec.Status)

	result := string(output)

	attributeMap := make(map[string]interface{})
	attributeMap["status"] = status
	attributeMap["output"] = result
	attributeMap["duration"] = duration

	sensuCheckJson := make(map[string]interface{})
	sensuCheckJson["check"] = attributeMap
	finalOuput, _ := json.Marshal(&sensuCheckJson)
	fmt.Println(string(finalOuput[:]))
	return nil
}
