package sensu

import (
	"fmt"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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
	assert.Equal(t, os.Stdin, goCheck.checkReader)
}

func goCheckExecuteUtil(
	t *testing.T,
	checkConfig *PluginConfig,
	checkFile string,
	cmdLineArgs []string,
	validationFunction func(check *types.Check) error,
) (int, string) {

	goCheck := NewGoCheck(checkConfig, nil, validationFunction)

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
	var validateCalled bool
	clearEnvironment()
	exitStatus, _ := goCheckExecuteUtil(t, &defaultHandlerConfig, "test/sensu-check.json", nil,
		func(check *types.Check) error {
			validateCalled = true
			assert.NotNil(t, check)
			return nil
		})
	assert.Equal(t, 0, exitStatus)
	assert.True(t, validateCalled)
}

// Test check
func TestGoCheck_Execute1(t *testing.T) {
	var validateCalled bool
	clearEnvironment()
	file, _ := ioutil.TempFile(os.TempDir(), "test")
	file.WriteString(string("hello world"))
	os.Stdin = file
	exitStatus, _ := goCheckExecuteUtil(t, &defaultHandlerConfig, "test/sensu-check1.json", nil,
		func(check *types.Check) error {
			validateCalled = true
			assert.NotNil(t, check)
			return nil
		})
	assert.Equal(t, 0, exitStatus)
	assert.True(t, validateCalled)
}
