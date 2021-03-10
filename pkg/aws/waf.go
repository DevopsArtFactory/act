package aws

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/waf"

	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/schema"
)

func GetWafClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *waf.WAF {
	if creds == nil {
		return waf.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return waf.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// SelectAcl selects target acl name from all acl list
func (c Client) SelectACL() (string, error) {
	var target string

	ACLs, err := c.GetAllWebACLs()
	if err != nil {
		return constants.EmptyString, err
	}

	var options []string
	for _, acl := range ACLs {
		options = append(options, fmt.Sprintf("%s / %s", *acl.Name, *acl.WebACLId))
	}

	prompt := &survey.Select{
		Message: "Choose the ACL: ",
		Options: options,
	}
	survey.AskOne(prompt, &target)

	if len(target) == 0 {
		return constants.EmptyString, errors.New("you canceled ACL selection")
	}

	return ParseWebACLID(target), nil
}

// GetAllWebACLs retrieves all ACLs in AWS WAF
func (c Client) GetAllWebACLs() ([]*waf.WebACLSummary, error) {
	input := &waf.ListWebACLsInput{
		Limit: aws.Int64(100),
	}

	result, err := c.WafClient.ListWebACLs(input)
	if err != nil {
		return nil, err
	}

	return result.WebACLs, nil
}

// DescribeWebACL describes web acl
func (c Client) DescribeWebACL(target string) (*schema.WebACL, error) {
	var ret schema.WebACL
	info, err := c.GetWebACLInfo(target)
	if err != nil {
		return nil, nil
	}

	// basic information
	ret.ID = *info.WebACLId
	ret.Name = *info.Name

	if len(info.Rules) > 0 {
		rules := []schema.ACLRule{}
		for _, rule := range info.Rules {
			rules = append(rules, schema.ACLRule{
				Type:       *rule.Type,
				ActionType: *rule.Action.Type,
				Priority:   *rule.Priority,
				RuleID:     *rule.RuleId,
			})
		}
		ret.Rules = rules
	}

	if len(ret.Rules) > 0 {
		for i := range ret.Rules {
			ruleInfo, err := c.DescribeRule(ret.Rules[i].RuleID)
			if err != nil {
				return nil, err
			}

			var dataSet []schema.IPDataSet
			if len(ruleInfo.Predicates) > 0 {
				for _, p := range ruleInfo.Predicates {
					data, err := c.GetIPSet(*p.DataId)
					if err != nil {
						return nil, err
					}

					tempIPSet := schema.IPDataSet{
						ID: *p.DataId,
					}

					l := []string{}
					if len(data.IPSetDescriptors) > 0 {
						for _, descriptor := range data.IPSetDescriptors {
							l = append(l, *descriptor.Value)
						}
					}

					tempIPSet.IPList = l
					dataSet = append(dataSet, tempIPSet)
				}

				ret.Rules[i].IPDataSet = dataSet
			}
		}
	}

	return &ret, nil
}

// GetWebACLInfo retrieves web acl information
func (c Client) GetWebACLInfo(target string) (*waf.WebACL, error) {
	input := &waf.GetWebACLInput{
		WebACLId: aws.String(target),
	}

	result, err := c.WafClient.GetWebACL(input)
	if err != nil {
		return nil, err
	}

	return result.WebACL, nil
}

// DescribeRule describes web ACL rule
func (c Client) DescribeRule(ruleID string) (*waf.Rule, error) {
	input := &waf.GetRuleInput{
		RuleId: aws.String(ruleID),
	}

	result, err := c.WafClient.GetRule(input)
	if err != nil {
		return nil, err
	}

	return result.Rule, nil
}

// GetIPSet retrieves information about IP Set
func (c Client) GetIPSet(dataID string) (*waf.IPSet, error) {
	input := &waf.GetIPSetInput{
		IPSetId: aws.String(dataID),
	}

	result, err := c.WafClient.GetIPSet(input)
	if err != nil {
		return nil, err
	}

	return result.IPSet, nil
}

// ParseWebACLID parses web ACL ID from option string
func ParseWebACLID(str string) string {
	return strings.TrimSpace(strings.Split(str, "/")[1])
}
