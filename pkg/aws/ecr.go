package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func GetEcrClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *ecr.ECR {
	if creds == nil {
		return ecr.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return ecr.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// GetAuthorizeToken retrieves authorize token via API
func (c Client) GetAuthorizeToken() (*ecr.AuthorizationData, error) {
	input := &ecr.GetAuthorizationTokenInput{}

	result, err := c.ECRClient.GetAuthorizationToken(input)
	if err != nil {
		return nil, err
	}

	if result.AuthorizationData == nil || len(result.AuthorizationData) == 0 {
		return nil, fmt.Errorf("there is no authorization data")
	}

	return result.AuthorizationData[0], nil
}
