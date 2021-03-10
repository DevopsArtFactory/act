package aws

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"

	"github.com/DevopsArtFactory/act/pkg/tools"
)

func GetIAMClientFn(sess client.ConfigProvider, creds *credentials.Credentials) *iam.IAM {
	if creds == nil {
		return iam.New(sess)
	}
	return iam.New(sess, aws.NewConfig().WithCredentials(creds))
}

// CreateNewCredentials creates new ACCESS_KEY, SECRET_ACCESS_KEY
func (c Client) CreateNewCredentials(name string) (*iam.AccessKey, error) {
	input := &iam.CreateAccessKeyInput{
		UserName: aws.String(name),
	}

	result, err := c.IAMClient.CreateAccessKey(input)
	if err != nil {
		return nil, err
	}

	return result.AccessKey, nil
}

// DeleteAccessKey deletes access key of user
func (c Client) DeleteAccessKey(accessKey, userName string) error {
	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accessKey),
		UserName:    aws.String(userName),
	}

	_, err := c.IAMClient.DeleteAccessKey(input)
	if err != nil {
		return err
	}

	return nil
}

// CheckAccessKeyExpired check whether access key is expired or not
func (c Client) CheckAccessKeyExpired(name, accessKeyID string) error {
	keys, err := c.GetAccessKeyList(name)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if *key.AccessKeyId == accessKeyID {
			if tools.IsExpired(*key.CreateDate, 24*180*time.Hour) {
				return errors.New("your access key is expired. please renew by running `act renew-credential`")
			}

			logrus.Debugf("your access key is not expired")
			return nil
		}
	}
	return errors.New("your access key configuration is wrong")
}

// GetAccessKeyList lists all access key of user
func (c Client) GetAccessKeyList(name string) ([]*iam.AccessKeyMetadata, error) {
	if err := tools.ClearOsEnv(); err != nil {
		return nil, nil
	}
	svc := iam.New(GetAwsSession())
	input := &iam.ListAccessKeysInput{
		UserName: aws.String(name),
	}

	result, err := svc.ListAccessKeys(input)
	if err != nil {
		return nil, err
	}

	return result.AccessKeyMetadata, nil
}
