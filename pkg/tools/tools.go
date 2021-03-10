package tools

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
)

// ClearOsEnv removes all environment variables about AWS
func ClearOsEnv() error {
	logrus.Debugf("remove environment variable")
	if err := os.Unsetenv("AWS_ACCESS_KEY_ID"); err != nil {
		return err
	}
	if err := os.Unsetenv("AWS_SECRET_ACCESS_KEY"); err != nil {
		return err
	}

	if err := os.Unsetenv("AWS_SESSION_TOKEN"); err != nil {
		return err
	}

	return nil
}

// Check if file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CopyToClipBoard copies token to clipboard
func CopyToClipBoard(token string) error {
	pbcopy := exec.Command("pbcopy")
	in, _ := pbcopy.StdinPipe()

	if err := pbcopy.Start(); err != nil {
		return err
	}

	if _, err := in.Write([]byte(token)); err != nil {
		return err
	}

	if err := in.Close(); err != nil {
		return err
	}

	err := pbcopy.Wait()
	if err != nil {
		return err
	}

	logrus.Info("Token is copied to clipboard.")

	return nil
}

//Figure out if string is in array
func IsStringInArray(s string, arr []string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}

	return false
}

//CreateFile creates/overrides files
func CreateFile(filePath string, writeData string) error {
	if err := ioutil.WriteFile(filePath, []byte(writeData), 0644); err != nil {
		return err
	}
	return nil
}

// AskContinue provides interactive terminal for users to answer if they continue process or not
func AskContinue(msg string) error {
	var answer string
	prompt := &survey.Input{
		Message: msg,
	}
	survey.AskOne(prompt, &answer)

	if IsStringInArray(strings.ToLower(answer), []string{"yes", "y"}) {
		return nil
	}

	return errors.New("stop process")
}

// GetKeys returns key list from map structure
func GetKeys(m map[string]string) []string {
	ret := []string{}
	for k := range m {
		ret = append(ret, k)
	}

	return ret
}

// IsExpired compares current time with (targetDate + timeAdded)
func IsExpired(targetDate time.Time, timeAdded time.Duration) bool {
	return time.Since(targetDate.Add(timeAdded)) > 0
}

// SetUpLogs set logrus log format
func SetUpLogs(stdErr io.Writer, level string) error {
	logrus.SetOutput(stdErr)
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}
	logrus.SetLevel(lvl)
	return nil
}

// Ask provides interactive terminal for users to answer and return string
func Ask(msg string, isSecret bool) (string, error) {
	var answer string
	var prompt survey.Prompt
	if isSecret {
		prompt = &survey.Password{
			Message: msg,
		}
	} else {
		prompt = &survey.Input{
			Message: msg,
		}
	}
	survey.AskOne(prompt, &answer)

	if len(answer) == 0 {
		return answer, errors.New("answer is required")
	}

	return answer, nil
}

// FormatKeyForDisplay displays key with masking
func FormatKeyForDisplay(k string) string {
	return fmt.Sprintf("****************%s", k[len(k)-4:])
}
