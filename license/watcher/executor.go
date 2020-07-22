package watcher

import (
	"bytes"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

type executor interface {
	execute(shellCommandArgs []string) (string, error)
}

func newExecutor() executor {
	return &defaultExecutor{}
}

type defaultExecutor struct{}

func (de *defaultExecutor) execute(shellCommandArgs []string) (string, error) {
	argumentRemainder := []string{}
	if len(shellCommandArgs) > 1 {
		argumentRemainder = shellCommandArgs[1:]
	}

	cmd := exec.Command(shellCommandArgs[0], argumentRemainder...)
	cmd.Env = os.Environ()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// Output() runs the actual binary
	output, err := cmd.Output()
	outputStr := string(output)

	if err != nil {
		return err.Error(), errors.Wrapf(err, "Command %s returned error: %s", shellCommandArgs, err.Error())
	}

	log.Infof("Command %v returned successfully: %s", shellCommandArgs, outputStr)
	return outputStr, nil
}
