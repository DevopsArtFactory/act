package schema

type Config struct {
	Profile     string              `yaml:"profile"`
	Name        string              `yaml:"name"`
	Duration    int                 `yaml:"duration"`
	Alias       map[string]string   `yaml:"alias"`
	AssumeRoles map[string]string   `yaml:"assume_roles"`
	Databases   map[string][]string `yaml:"databases"`
	Maintenance struct {
		Message string `yaml:"message"`
		Arns    []struct {
			LoadbalancerArn string `yaml:"loadbalancer_arn"`
		}
	} `yaml:"maintenance"`
	Loadtest struct {
		RDS []string `yaml:"rds"`
		ASG []string `yaml:"asg"`
	} `yaml:"loadtest"`
}

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	CreateDate      string
}

type WebACL struct {
	ID    string
	Name  string
	Rules []ACLRule
}

type ACLRule struct {
	Type       string
	ActionType string
	Priority   int64
	RuleID     string
	IPDataSet  []IPDataSet
}

type IPDataSet struct {
	ID     string
	IPList []string
}

type IPCheckResult struct {
	IP      string
	IPSetID string
	Result  bool
}

type MaintenanceMessage struct {
	Timestamp int64  `json:"timestamp"`
	Code      int64  `json:"code"`
	Message   string `json:"message"`
}
