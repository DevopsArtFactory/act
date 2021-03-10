package constants

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	//ServiceName means name of service
	ServiceName = "act"

	// DefaultLogLevel is the default global verbosity
	DefaultLogLevel = logrus.WarnLevel

	// DefaultRegion is the default region id
	DefaultRegion = "ap-northeast-2"

	// EmptyString is the empty string
	EmptyString = ""

	// DefaultProfile is the default aws profile
	DefaultProfile = "default"

	// DefaultDuration is the default duration of assume role
	DefaultDuration = 0

	// InfoLogLevel is the info level verbosity
	InfoLogLevel = logrus.InfoLevel

	// ConfigErrorMsg is the default error message when there is no configuration set
	ConfigErrorMsg = "you have no configuration setting"

	// NeedExpiredCheck
	NeedExpiredCheck = true

	// SkipExpiredCheck
	SkipExpiredCheck = false

	// Secret
	Secret = true

	// ExecuteDelimiter is a delimiter for command execution
	ExecuteDelimiter = "--"

	// DefaultExpirationWindow default value of expiration window for assume role
	DefaultExpirationWindow = 5 * time.Minute

	// DefaultKeyType default type of key storage
	DefaultKeyType = "keychain"

	// DefaultMaintenancePriority means the default value of rule priority
	DefaultMaintenancePriority = 10
)

var (
	AWSConfigDirectoryPath = HomeDir() + "/.aws"
	AWSCredentialsPath     = AWSConfigDirectoryPath + "/credentials"
	BaseFilePath           = AWSConfigDirectoryPath + "/config.yaml"
	BaseSerialNumber       = "arn:aws:iam::748177903968:mfa"

	DefaultKeyChainPath    = fmt.Sprintf("%s-vault.keychain", ServiceName)
	DefaultKeyChainAccount = fmt.Sprintf("%s-default", ServiceName)
)

// Get Home Directory
func HomeDir() string {
	if h := os.Getenv("HOME"); h != EmptyString {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
