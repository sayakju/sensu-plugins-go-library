package sensu

import (
	"fmt"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	defaultCheckConfig = PluginConfig{
		Name:     "TestHandler",
		Short:    "Short Description",
		Timeout:  10,
		Keyspace: "sensu.io/plugins/segp/config",
	}
)

func TestNewGoCheck(t *testing.T) {
	goCheck := NewGoCheck(&defaultHandlerConfig, nil, func(check *types.Check) error {
		return nil
	})

	assert.NotNil(t, goCheck)
	assert.Equal(t, &defaultCheckConfig, goCheck.config)
	assert.NotNil(t, goCheck.validationFunction)
	assert.NotNil(t, goCheck.executeFunction)
	assert.Nil(t, goCheck.sensuCheck)
	assert.Equal(t, os.Stdin, goCheck.eventReader)
}

func goCheckExecuteUtil(
	t *testing.T,
	checkConfig *PluginConfig,
	checkFile string,
	cmdLineArgs []string,
	validationFunction func(check *types.Check) error,
	expectedValue1 interface{},
	expectedValue2 interface{},
	expectedValue3 interface{},
) (int, string) {
	values := handlerValues{}
	options := getHandlerOptions(&values)

	goCheck := NewGoCheck(checkConfig, options, validationFunction)

	// Simulate the command line arguments if necessary
	if len(cmdLineArgs) > 0 {
		goCheck.cmdArgs.SetArgs(cmdLineArgs)
	} else {
		goCheck.cmdArgs.SetArgs([]string{})
	}

	// Replace stdin reader with file reader and exitFunction with our own so we can know the exit status
	var exitStatus int
	var errorStr = ""
	goCheck.eventReader = getFileReader(checkFile)
	goCheck.exitFunction = func(i int) {
		exitStatus = i
	}
	goCheck.errorLogFunction = func(format string, a ...interface{}) {
		errorStr = fmt.Sprintf(format, a...)
	}
	goCheck.Execute()

	assert.Equal(t, expectedValue1, values.arg1)
	assert.Equal(t, expectedValue2, values.arg2)
	assert.Equal(t, expectedValue3, values.arg3)

	return exitStatus, errorStr
}

// Test check override
func TestGoCheck_Execute(t *testing.T) {
	var validateCalled, executeCalled bool
	clearEnvironment()
	exitStatus, _ := goCheckExecuteUtil(t, &defaultHandlerConfig, "test/sensu-check.json", nil,
		func(event *types.Check) error {
			validateCalled = true
			assert.NotNil(t, event)
			return nil
		},
		"value-check1", uint64(1357), false)
	assert.Equal(t, 0, exitStatus)
	assert.True(t, validateCalled)
	assert.True(t, executeCalled)
}
