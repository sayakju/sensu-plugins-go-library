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
		Name:     "TestCheck",
		Short:    "Short Description",
		Timeout:  10,
		Keyspace: "sensu.io/plugins/segp/config",
	}
)

func TestNewGoCheck(t *testing.T) {
	goCheck := NewGoCheck(&defaultCheckConfig, nil, func(check *types.Check, entity *types.Entity) error {
		return nil
	}, func(check *types.Check, entity *types.Entity) (int, error) {
		return 0, nil
	})

	assert.NotNil(t, goCheck)
	assert.Equal(t, &defaultCheckConfig, goCheck.config)
	assert.NotNil(t, goCheck.validationFunction)
	assert.NotNil(t, goCheck.executeFunction)
	assert.False(t, goCheck.readEvent)
	assert.True(t, goCheck.readCheck)
	assert.Nil(t, goCheck.sensuCheck)
	assert.Equal(t, os.Stdin, goCheck.checkReader)
}

func goCheckExecuteUtil(t *testing.T, checkConfig *PluginConfig, checkFile string, cmdLineArgs []string,
	validationFunction func(check *types.Check, entity *types.Entity) error, executeFunction func(*types.Check, *types.Entity) (int, error)) (int, string) {

	goCheck := NewGoCheck(checkConfig, nil, validationFunction, executeFunction)

	// Simulate the command line arguments if necessary
	if len(cmdLineArgs) > 0 {
		goCheck.cmdArgs.SetArgs(cmdLineArgs)
	} else {
		goCheck.cmdArgs.SetArgs([]string{})
	}

	// Replace stdin reader with file reader and exitFunction with our own so we can know the exit status
	var exitStatus int
	var errorStr = ""
	goCheck.checkReader = getFileReader(checkFile)
	goCheck.exitFunction = func(i int) {
		exitStatus = i
	}
	goCheck.errorLogFunction = func(format string, a ...interface{}) {
		errorStr = fmt.Sprintf(format, a...)
	}
	goCheck.Execute()

	return exitStatus, errorStr
}

// Test check
func TestGoCheck_Execute(t *testing.T) {
	var validateCalled, executeCalled bool
	clearEnvironment()
	exitStatus, _ := goCheckExecuteUtil(t, &defaultCheckConfig, "test/sensu-check.json", nil,
		func(check *types.Check, entity *types.Entity) error {
			validateCalled = true
			assert.NotNil(t, check)
			return nil
		}, func(check *types.Check, entity *types.Entity) (int, error) {
			executeCalled = true
			assert.NotNil(t, check)
			return 0, nil
		})
	assert.Equal(t, 0, exitStatus)
	assert.True(t, validateCalled)
	assert.True(t, executeCalled)
}
