package aws

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/elbv2"

	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/schema"
)

func GetELBClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *elbv2.ELBV2 {
	if creds == nil {
		return elbv2.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return elbv2.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

func (c Client) DescribeListeners(loadbalancerArn string) ([]*elbv2.Listener, error) {
	input := &elbv2.DescribeListenersInput{
		LoadBalancerArn: aws.String(loadbalancerArn),
	}

	result, err := c.ELBClient.DescribeListeners(input)
	if err != nil {
		return nil, err
	}

	return result.Listeners, nil
}

// CreateMaintenanceRule creates maintenance rule
func (c Client) CreateMaintenanceRule(listenerArn *string, message string) error {
	maintenanceMessage, _ := json.Marshal(schema.MaintenanceMessage{
		Timestamp: 1565934483777,
		Code:      -99999,
		Message:   message,
	})

	fixedResponseActionConfig := &elbv2.FixedResponseActionConfig{
		ContentType: aws.String("application/json"),
		MessageBody: aws.String(string(maintenanceMessage)),
		StatusCode:  aws.String("503"),
	}

	action := &elbv2.Action{
		FixedResponseConfig: fixedResponseActionConfig,
		Type:                aws.String("fixed-response"),
	}

	input := &elbv2.CreateRuleInput{
		Actions: []*elbv2.Action{action},
		Conditions: []*elbv2.RuleCondition{
			{
				Field: aws.String("path-pattern"),
				Values: []*string{
					aws.String("/*"),
				},
			},
		},
		ListenerArn: listenerArn,
		Priority:    aws.Int64(constants.DefaultMaintenancePriority),
	}

	_, err := c.ELBClient.CreateRule(input)
	if err != nil {
		return err
	}

	return nil
}

// DescribeRules describes list of rules in listener
func (c Client) DescribeRules(listenerArn *string) ([]*elbv2.Rule, error) {
	input := &elbv2.DescribeRulesInput{
		ListenerArn: listenerArn,
	}

	result, err := c.ELBClient.DescribeRules(input)
	if err != nil {
		return nil, err
	}

	return result.Rules, nil
}

// DeleteMaintenanceRule deletes maintenance rule
func (c Client) DeleteMaintenanceRule(ruleArn *string) error {
	input := &elbv2.DeleteRuleInput{
		RuleArn: ruleArn,
	}

	_, err := c.ELBClient.DeleteRule(input)
	if err != nil {
		return err
	}

	return nil
}
