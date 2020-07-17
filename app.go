package main

import (
	"fmt"
	"github.com/cloudogu/confluence-license-checker/license/tester"
	"github.com/cloudogu/confluence-license-checker/license/watcher"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	watchIntervalFlagName  = "watch-interval"
	setupLicenseFlagName   = "setup-license"
	setupLicenseEnvVarName = "SETUP_LICENSE"
	confluenceConfigFile   = "/var/atlassian/confluence/confluence.cfg.xml"
)

var (
	// Version of the application
	Version string
)

// logging format
var format = logging.MustStringFormatter(
	`{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
)

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "define log level",
			Value: "warning",
		},
		&cli.BoolFlag{
			Name:  "show-stack",
			Usage: "show stacktrace on errors",
		},
		&cli.BoolFlag{
			Name:  "skip-root",
			Usage: "skip root check",
		},
	}
}

func configureLogging(c *cli.Context) error {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logLevel, err := logging.LogLevel(c.String("log-level"))
	if err != nil {
		fmt.Println("invalid log level specified, please use critical, error, warning, notice, info or debug")
		return errors.Wrap(err, "failed to configure logging")
	}
	logging.SetLevel(logLevel, "")
	return nil
}

func isPrintStack() bool {
	for _, arg := range os.Args {
		if arg == "--show-stack" {
			return true
		}
	}
	return false
}

func checkMainError(err error, ex exiter) {
	if err != nil {
		if isPrintStack() {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Printf("%+s\n", err)
		}
		ex.exit(1)
	}
}

// projects main function
func main() {
	app := cli.NewApp()
	app.Name = "license-checker"
	app.Usage = "a tool that checks for a Confluence license"
	app.Version = Version
	app.Commands = []*cli.Command{WatchCommand(), TestLicenseCommand()}

	app.Flags = createGlobalFlags()
	app.Before = configureLogging

	err := app.Run(os.Args)
	exiter := newExiter()
	checkMainError(err, exiter)
}

func WatchCommand() *cli.Command {
	return &cli.Command{
		Name:  "watch",
		Usage: "watch for a Confluence license change and execute a command",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    watchIntervalFlagName,
				Aliases: []string{"w"},
				Usage:   "the watch interval in seconds",
				Value:   30,
			},
			&cli.StringFlag{
				Name:    setupLicenseFlagName,
				Aliases: []string{"l"},
				Usage:   "provides a license to watch instead from a environment variable",
				EnvVars: []string{setupLicenseEnvVarName},
			},
		},
		Action: watchExecuteAction,
	}
}

func TestLicenseCommand() *cli.Command {
	return &cli.Command{
		Name:  "test-setup",
		Usage: "check if a setup-specific Confluence license is currently configured",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    setupLicenseFlagName,
				Aliases: []string{"l"},
				Usage:   "provides a license to watch instead from a environment variable",
				EnvVars: []string{setupLicenseEnvVarName},
			},
		},
		Action: TestLicenseAction,
	}
}

func watchExecuteAction(c *cli.Context) error {
	watchInterval := c.Int(watchIntervalFlagName)
	if watchInterval < 1 {
		return errors.Errorf("cannot start license watcher: value for flag '--%s' must be greater than zero", watchIntervalFlagName)
	}

	if c.NArg() == 0 {
		err := cli.ShowAppHelp(c)
		return errors.Wrap(err, "cannot start license watcher: a shell command must be provided")
	}

	license := c.String(setupLicenseFlagName)
	if license == "" {
		return errors.Errorf("cannot start license watcher: a start license must be provided either by flag '--%s' or by environment variable '${%s}'",
			setupLicenseFlagName, setupLicenseEnvVarName)
	}

	args := &watcher.ProcessArgs{
		CommandArgs:          c.Args().Slice(),
		WatchIntervalInSecs:  watchInterval,
		ConfluenceConfigFile: confluenceConfigFile,
		SetupLicense:         license,
	}

	ex := watcher.New(args)
	err := ex.Watch()
	if err != nil {
		return errors.Wrap(err, "license watcher failed with an error")
	}

	fmt.Println("Confluence license watcher quits.")
	return nil
}

func TestLicenseAction(c *cli.Context) error {
	license := c.String(setupLicenseFlagName)
	if license == "" {
		return errors.Errorf("cannot start license watcher: a start license must be provided either by flag '--%s' or by environment variable '${%s}'",
			setupLicenseFlagName, setupLicenseEnvVarName)
	}

	licTester := tester.New()
	hasSetupLic, err := licTester.HasSetupLicense(confluenceConfigFile, license)
	if err != nil {
		return errors.Wrap(err, "license watcher failed with an error")
	}

	if !hasSetupLic {
		return errors.New("Found a non-setup license. License check must not be started.")
	}

	fmt.Println("Confluence license watcher quits.")
	return nil
}

type exiter interface {
	exit(exitCode int)
}

func newExiter() exiter {
	return &defaultExiter{}
}

type defaultExiter struct{}

func (*defaultExiter) exit(exitCode int) {
	os.Exit(exitCode)
}
