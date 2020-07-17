package watcher

import (
	"github.com/cloudogu/confluence-license-checker/license/tester"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"time"
)

var log = logging.MustGetLogger("watcher")

// Watcher detects a change from setup license to a valid production license in a given Confluence configuration.
// A given action must be executed if a license change is detected.
type Watcher interface {
	// Watch watches for license changes.
	Watch() error
}

// ProcessArgs contain necessary arguments
type ProcessArgs struct {
	// CommandArgs is a shell call and necessary arguments that will be executed if a license change is detected.
	// The first argument is the actual command to be executed, and must be present. Any further arguments are optional
	// and depend on the command.
	//
	// Example:
	// 	[]string{ "/bin/echo", "-n", "hello world" }
	CommandArgs []string
	// WatchIntervalInSecs is the interval in seconds in which the config file is inspected for a license change.
	WatchIntervalInSecs int
	// ConfluenceConfigFile is the file which accommodates the license to be watched.
	ConfluenceConfigFile string
	// SetupLicense is the license with which the setup should be executed.
	SetupLicense string
}

// New creates a new Watcher instance.
func New(args *ProcessArgs) Watcher {
	log.Debugf("Found these arguments: %v", args)

	executor := newExecutor()
	licenseChecker := tester.New()

	return &defaultWatcher{
		args:          args,
		cmdExecutor:   executor,
		licenseTester: licenseChecker,
	}
}

type defaultWatcher struct {
	args          *ProcessArgs
	cmdExecutor   executor
	licenseTester tester.Tester
}

// Watch watches in a fixed interval for license changes.
func (dw *defaultWatcher) Watch() error {
	duration := time.Duration(dw.args.WatchIntervalInSecs) * time.Second
	log.Debugf("Start License check using %d seconds", dw.args.WatchIntervalInSecs)

	//lint:ignore SA1015 (don't warn about leaking ticker)
	for range time.Tick(duration) {
		done, err := dw.doWatchWork()
		if err != nil {
			return errors.Wrap(err, "exiting watcher because an error occurred")
		}

		if done {
			return nil
		}
	}

	return nil
}

func (dw *defaultWatcher) doWatchWork() (done bool, err error) {
	log.Debugf("License check time: %s", time.Now().Format(time.RFC3339))

	log.Debug("Checking for license change.")
	changed, err := dw.licenseTester.HasLicenseChanged(dw.args.ConfluenceConfigFile, dw.args.SetupLicense)
	if err != nil {
		return true, err
	}

	if changed {
		log.Debug("Found change.")
		_, err = dw.cmdExecutor.execute(dw.args.CommandArgs)
		return true, err
	}

	log.Debugf("No change found. Checking again in %d seconds.", dw.args.WatchIntervalInSecs)
	return
}
