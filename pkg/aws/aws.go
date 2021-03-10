package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/waf"
)

const PORT = 3310

type Client struct {
	RDSClient *rds.RDS
	STSClient *sts.STS
	S3Client  *s3.S3
	WafClient *waf.WAF
	ECRClient *ecr.ECR
	IAMClient *iam.IAM
	ELBClient *elbv2.ELBV2
	ASGClient *autoscaling.AutoScaling
	Region    string
}

func NewClient(sess client.ConfigProvider, region string, creds *credentials.Credentials) Client {
	return Client{
		RDSClient: GetRDSClientFn(sess, region, creds),
		STSClient: GetSTSClientFn(sess, region, creds),
		S3Client:  GetS3ClientFn(sess, region, creds),
		WafClient: GetWafClientFn(sess, region, creds),
		IAMClient: GetIAMClientFn(sess, creds),
		ELBClient: GetELBClientFn(sess, region, creds),
		ECRClient: GetEcrClientFn(sess, region, creds),
		ASGClient: GetASGClientFn(sess, region, creds),
	}
}

//Get DB Auth token
func GetDBAuthToken(target, region, user string, creds *credentials.Credentials) (string, error) {
	endpoint := fmt.Sprintf("%s:%d", target, PORT)
	return rdsutils.BuildAuthToken(endpoint, region, user, creds)
}

// HeadS3Bucket checks if s3 bucket exists or not
func (c Client) HeadS3Bucket(bucket string) error {
	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err := c.S3Client.HeadBucket(input)
	if err != nil {
		return err
	}

	return nil
}
