package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func GetASGClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *autoscaling.AutoScaling {
	if creds == nil {
		return autoscaling.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return autoscaling.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

func (c Client) GetExactASGNames(asgName string) ([]*string, error) {
	var response []*string
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	result, err := c.ASGClient.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}

	for _, v := range result.AutoScalingGroups {
		if strings.Contains(*v.AutoScalingGroupName, asgName) {
			response = append(response, v.AutoScalingGroupName)
		}
	}

	return response, nil
}

func (c Client) GetLoadtestASGStatus(asgName *string) (*int64, error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			asgName,
		},
		MaxRecords: aws.Int64(100),
	}

	result, err := c.ASGClient.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}

	return result.AutoScalingGroups[0].DesiredCapacity, nil
}

func (c Client) updateLoadtestASG(asgName *string, capacity int64) error {
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: asgName,
		MinSize:              aws.Int64(capacity),
		DesiredCapacity:      aws.Int64(capacity),
		MaxSize:              aws.Int64(capacity),
	}

	_, err := c.ASGClient.UpdateAutoScalingGroup(input)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) StartLoadtestASG(asgName *string) error {
	return c.updateLoadtestASG(asgName, 1)
}

func (c Client) StopLoadtestASG(asgName *string) error {
	return c.updateLoadtestASG(asgName, 0)
}
