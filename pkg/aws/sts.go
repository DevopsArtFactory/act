package aws

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/DevopsArtFactory/act/pkg/color"
)

func GetSTSClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *sts.STS {
	if creds == nil {
		return sts.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return sts.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// CheckWhoIam calls get-caller-identity and print the result
func (c Client) CheckWhoIam(out io.Writer) error {
	result, err := c.STSClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}

	color.Green.Fprintf(out, result.String())

	return nil
}

// CheckMFAToken checks MFA authentication for specific functions
func (c Client) CheckMFAToken(name string) error {
	roleSessionName := getRoleSessionName(name)
	mfaSerialNumber := getMFASerialNumber(roleSessionName)
	mfaToken, err := AskMFAToken()
	if err != nil {
		return err
	}

	return c.GetSessionToken(mfaSerialNumber, mfaToken)
}

// getSessionToken retrieves session token with MFA authentication
func (c Client) GetSessionToken(mfaSerialNumber, mfaToken string) error {
	input := &sts.GetSessionTokenInput{
		SerialNumber: aws.String(mfaSerialNumber),
		TokenCode:    aws.String(mfaToken),
	}

	_, err := c.STSClient.GetSessionToken(input)
	if err != nil {
		return err
	}

	return nil
}

// AskMFAToken gets MFA token from command line interface
func AskMFAToken() (string, error) {
	var v string
	fmt.Fprintf(os.Stderr, "MFA token code: ")
	_, err := fmt.Scanln(&v)

	return v, err
}
