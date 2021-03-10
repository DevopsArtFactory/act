package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"

	awsact "github.com/DevopsArtFactory/act/pkg/aws"
)

// GetAssumeCreds creates a credentials for assuming role.
func GetAssumeCreds(arn string, sessionName string, duration int) (*sts.Credentials, error) {
	sess := awsact.GetAwsSession()
	svc := awsact.GetSTSClientFn(sess, "ap-northeast-2", nil)
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int64(int64(duration)),
	}

	result, err := svc.AssumeRole(input)
	if err != nil {
		return nil, err
	}

	return result.Credentials, nil
}
