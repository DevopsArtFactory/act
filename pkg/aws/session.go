package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/DevopsArtFactory/act/pkg/constants"
)

// GetAwsSession creates new session for AWS
func GetAwsSession() *session.Session {
	mySession := session.Must(session.NewSession())
	return mySession
}

// GetAwsSessionWithConfig creates new session for AWS
func GetAwsSessionWithConfig(region string) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		SharedConfigState: session.SharedConfigDisable,
	})
}

// GenerateCreds creates new credentials with MFA
func GenerateCreds(assumeRole, name string) *credentials.Credentials {
	awsSession := GetAwsSession()
	roleSessionName := getRoleSessionName(name)
	mfaSerialNumber := getMFASerialNumber(roleSessionName)

	var creds *credentials.Credentials

	mfaNumber, err := stscreds.StdinTokenProvider()
	if err != nil {
		return nil
	}

	if len(assumeRole) != 0 {
		creds = stscreds.NewCredentials(awsSession, assumeRole, func(p *stscreds.AssumeRoleProvider) {
			p.SerialNumber = aws.String(mfaSerialNumber)
			p.TokenCode = aws.String(mfaNumber)
			p.RoleSessionName = roleSessionName
		})
	}
	return creds
}

// GetRoleSessionName create username which will be used to session name
func getRoleSessionName(name string) string {
	return name
}

// GetMFASerialNumber generates MFA serial number
func getMFASerialNumber(email string) string {
	return fmt.Sprintf("%s/%s", constants.BaseSerialNumber, email)
}
