package tester

import (
	"bufio"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

var log = logging.MustGetLogger("tester")

type Tester interface {
	HasLicenseChanged(configFile string, knownLicense string) (changed bool, err error)
	HasSetupLicense(configFile string, setupLicense string) (unchanged bool, err error)
}

func New() Tester {
	opener := newFileOpener()
	return &defaultLicenseTester{opener: opener}
}

type defaultLicenseTester struct {
	opener fileOpener
}

func (lc *defaultLicenseTester) HasSetupLicense(configFile string, setupLicense string) (changed bool, err error) {
	licenseChanged, err := lc.HasLicenseChanged(configFile, setupLicense)
	return !licenseChanged, err
}

func (lc *defaultLicenseTester) HasLicenseChanged(configFile string, knownLicense string) (changed bool, err error) {
	log.Debugf("Checking configuration file '%s'", configFile)
	log.Debugf("Comparing to license '%s'", knownLicense)

	licenseLine, err := readLicenseLineFrom(configFile, lc.opener)
	if err != nil {
		return false, errors.Wrap(err, "failed to check license")
	}

	if strings.Contains(licenseLine, knownLicense) {
		log.Debug("Found old license in license file")
		return false, nil
	}

	log.Debugf("Detected a different license in configuration file")
	return true, nil
}

func readLicenseLineFrom(fileToCheck string, opener fileOpener) (string, error) {
	file, err := opener.Open(fileToCheck)
	if err != nil {
		return "", errors.Wrapf(err, "error while opening config file '%s'", fileToCheck)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	licenseProperty := "atlassian.license.message"
	fullLicenseLine := ""
	for scanner.Scan() {
		currentLine := scanner.Text()
		if strings.Contains(currentLine, licenseProperty) {
			fullLicenseLine = currentLine
		}
	}

	if fullLicenseLine == "" {
		return "", errors.Errorf("failed to find property '%s' in file '%s'", licenseProperty, fileToCheck)
	}

	return fullLicenseLine, nil
}

type fileOpener interface {
	Open(filePath string) (io.ReadCloser, error)
}

func newFileOpener() fileOpener {
	return &defaultFileOpener{}
}

type defaultFileOpener struct{}

func (d defaultFileOpener) Open(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}
