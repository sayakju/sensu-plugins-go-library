package sensu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
	goCheck := NewGoCheck(&defaultHandlerConfig, nil, func() error {
		return nil
	}, "whoami")

	assert.NotNil(t, goCheck)
	assert.Equal(t, &defaultCheckConfig, goCheck.config)
	assert.NotNil(t, goCheck.validationFunction)
	assert.False(t, goCheck.readEvent)
}

func goCheckExecuteUtil(
	t *testing.T,
	c string,
) (int, string) {

	goCheck := NewGoCheck(&defaultHandlerConfig, nil, func() error {
		return nil
	}, c)

	// Replace stdin reader with file reader and exitFunction with our own so we can know the exit status
	var exitStatus int
	var errorStr = ""
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
	clearEnvironment()
	exitStatus, _ := goCheckExecuteUtil(t, "whoami")
	assert.Equal(t, 0, exitStatus)
}

// Test check Invalid Command
func TestGoCheck_InvalidCommand(t *testing.T) {
	clearEnvironment()
	exitStatus, _ := goCheckExecuteUtil(t, "abc")
	assert.Equal(t, 0, exitStatus)
}
