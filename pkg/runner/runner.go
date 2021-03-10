package runner

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/DevopsArtFactory/act/pkg/aws"
	"github.com/DevopsArtFactory/act/pkg/builder"
	"github.com/DevopsArtFactory/act/pkg/color"
	"github.com/DevopsArtFactory/act/pkg/config"
	"github.com/DevopsArtFactory/act/pkg/constants"
	"github.com/DevopsArtFactory/act/pkg/schema"
	"github.com/DevopsArtFactory/act/pkg/templates"
	"github.com/DevopsArtFactory/act/pkg/tools"
)

type Runner struct {
	AWSClient aws.Client
	Flag      *builder.Flags
	Config    *schema.Config
}

func New(flags *builder.Flags, config *schema.Config) Runner {
	region := flags.Region
	if len(region) == 0 {
		region = constants.DefaultRegion
	}

	return Runner{
		AWSClient: aws.NewClient(aws.GetAwsSession(), region, nil),
		Flag:      flags,
		Config:    config,
	}
}

//CopyRDSToken copies RDS Token to clipboard
func (r Runner) CopyRDSToken(env, region string) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	if _, ok := config.AssumeRoles[env]; !ok {
		return errors.New("no assume role exists in config file")
	}

	targets := config.Databases[env]

	if len(targets) == 0 {
		return fmt.Errorf("no endpoints exist in configuration file for %s", env)
	}

	target, err := aws.SelectTarget(targets)
	if err != nil {
		return err
	}

	creds := aws.GenerateCreds(config.AssumeRoles[env], config.Name)

	authToken, err := aws.GetDBAuthToken(target, region, strings.Split(config.Name, "@")[0], creds)
	if err != nil {
		return err
	}

	if err := tools.CopyToClipBoard(authToken); err != nil {
		return err
	}
	return nil
}

//ChooseEnv provides interactive terminal to choose the environment
func (r Runner) ChooseEnv() (string, error) {
	config, err := config.GetConfig()
	if err != nil {
		return constants.EmptyString, err
	}

	environs := []string{}
	for key := range config.Databases {
		environs = append(environs, key)
	}

	fmt.Println(environs)

	var env string
	prompt := &survey.Select{
		Message: "Choose the environment: ",
		Options: environs,
	}
	survey.AskOne(prompt, &env)

	if len(env) == 0 {
		return env, errors.New("choosing an environment has been canceled")
	}

	return env, nil
}

// InitConfiguration init new configuration
func (r Runner) InitConfiguration() error {
	if tools.FileExists(constants.BaseFilePath) {
		return fmt.Errorf("you already had configuration file: %s", constants.BaseFilePath)
	}

	// check base AWS directory
	if !tools.FileExists(constants.AWSConfigDirectoryPath) {
		if err := os.MkdirAll(constants.AWSConfigDirectoryPath, 0755); err != nil {
			return err
		}
	}

	// Ask base account name which should be a company mail
	name, err := AskBaseAccountName()
	if err != nil {
		return err
	}

	c := config.GetInitConfig(name)
	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// show configuration file for double-check
	if err := printInitFile(y); err != nil {
		return err
	}

	// ask to continue
	if err := tools.AskContinue("Are you sure to generate configuration file? "); err != nil {
		return errors.New("initialization has been canceled")
	}

	if err := tools.CreateFile(constants.BaseFilePath, string(y)); err != nil {
		return err
	}
	color.Blue.Fprintf(os.Stdout, "New configuration file is successfully generated in %s", constants.BaseFilePath)

	return nil
}

// Setup gets assume role information and copy them to the clipboard
func (r Runner) Setup(out io.Writer, args []string) error {
	var target string
	var err error

	if r.Config == nil {
		return errors.New(constants.ConfigErrorMsg)
	}

	if len(args) == 0 {
		target, err = AskAssumeTarget(r.Config.AssumeRoles)
		if err != nil {
			return err
		}
	} else {
		target = args[0]
	}

	var arn string
	if len(r.Config.Alias) > 0 {
		arn = r.Config.AssumeRoles[r.Config.Alias[target]]
	}

	if len(arn) == 0 {
		arn = r.Config.AssumeRoles[target]
	}

	if err := CheckTarget(arn, target); err != nil {
		return err
	}

	duration := r.Config.Duration
	if r.Flag.Duration > 0 {
		duration = r.Flag.Duration
	}

	assumeCreds, err := config.GetAssumeCreds(arn, r.Config.Name, duration)
	if err != nil {
		return err
	}

	rawOutput := viper.GetBool("raw-output")

	if rawOutput {
		fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", *assumeCreds.AccessKeyId)
		fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", *assumeCreds.SecretAccessKey)
		fmt.Printf("export AWS_SESSION_TOKEN=%s\n", *assumeCreds.SessionToken)
	} else {
		loc, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			return err
		}

		pbcopy := exec.Command("pbcopy")

		in, _ := pbcopy.StdinPipe()

		if err := pbcopy.Start(); err != nil {
			return err
		}

		if _, err := in.Write([]byte(fmt.Sprintf("export AWS_ACCESS_KEY_ID=%s\n", *assumeCreds.AccessKeyId))); err != nil {
			return err
		}

		if _, err := in.Write([]byte(fmt.Sprintf("export AWS_SECRET_ACCESS_KEY=%s\n", *assumeCreds.SecretAccessKey))); err != nil {
			return err
		}

		if _, err := in.Write([]byte(fmt.Sprintf("export AWS_SESSION_TOKEN=%s\n", *assumeCreds.SessionToken))); err != nil {
			return err
		}

		if err := in.Close(); err != nil {
			return err
		}

		err = pbcopy.Wait()
		if err != nil {
			color.Red.Fprintln(out, err.Error())
			return err
		}
		color.Red.Fprintf(out, "Current token expired at: %s", assumeCreds.Expiration.In(loc))
		color.Blue.Fprintln(out, "Assume Credentials copied to clipboard, please paste it.")
	}
	return nil
}

// PrintAssumeList prints all accounts registered for assuming
func (r Runner) PrintAssumeList(out io.Writer) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	color.Blue.Fprintf(out, "[ name: %s, profile: %s ] Account List", config.Name, config.Profile)
	for key := range config.AssumeRoles {
		fmt.Println(key)
	}

	return nil
}

// Who prints the result of `aws sts get-caller-identity`
func (r Runner) Who(out io.Writer) error {
	if err := r.AWSClient.CheckWhoIam(out); err != nil {
		return err
	}

	return nil
}

func (r Runner) HasIP(out io.Writer, args []string) error {
	if len(args) == 0 {
		return errors.New("you have to specify at least 1 ip address")
	}

	targetList, err := ParseTargetList(args)
	if err != nil {
		return err
	}

	acl, err := r.SelectTargetACL([]string{})
	if err != nil {
		return err
	}

	resultSet, err := r.CheckIfTargetExistsInACL(targetList, acl)
	if err != nil {
		return err
	}

	if err := PrintHasIPResult(out, resultSet); err != nil {
		return err
	}

	return nil
}

// DescribeWebACL retrieves waf ip list and show them on the terminal
func (r Runner) DescribeWebACL(out io.Writer, args []string) error {
	target, err := r.SelectTargetACL(args)
	if err != nil {
		return err
	}

	info, err := r.AWSClient.DescribeWebACL(target)
	if err != nil {
		return err
	}

	if err := PrintWebACL(out, info); err != nil {
		return err
	}

	return nil
}

// ResetCredentials creates new credentials
func (r Runner) RenewCredentials(out io.Writer) error {
	// set log level to info
	logrus.SetLevel(constants.InfoLogLevel)

	// get configuration
	c := r.Config

	if c == nil {
		return errors.New(constants.ConfigErrorMsg)
	}

	if err := r.AWSClient.CheckMFAToken(c.Name); err != nil {
		return err
	}

	logrus.Infof("reading current credentials in %s", constants.AWSCredentialsPath)
	cfg, err := config.ReadAWSConfig()
	if err != nil {
		return err
	}

	logrus.Infof("checking current profile: %s", c.Profile)
	section, err := cfg.GetSection(c.Profile)
	if err != nil {
		return err
	}

	// get current key
	currentAccessKey := section.Key("aws_access_key_id").String()
	logrus.Infof("your current access key is: %s", currentAccessKey)

	logrus.Infof("checking current access key list in profile: %s", c.Profile)
	keys, err := r.AWSClient.GetAccessKeyList(c.Name)
	if err != nil {
		return err
	}

	logrus.Infof("you have %d access key(s) in profile: %s", len(keys), c.Profile)
	isDeleted := false
	if len(keys) == 2 {
		logrus.Warnf("delete current access key because you already have two credentials in profile: %s", c.Profile)
		if err := r.AWSClient.DeleteAccessKey(currentAccessKey, c.Name); err != nil {
			return err
		}
		isDeleted = true
	}

	logrus.Infof("creating new credential for %s profile", c.Profile)
	newCredentials, err := r.AWSClient.CreateNewCredentials(c.Name)
	if err != nil {
		return err
	}
	logrus.Infof("New credential is successfully created for %s profile", c.Profile)

	for _, s := range cfg.Sections() {
		if s.Name() == c.Profile {
			cfg.Section(s.Name()).Key("aws_access_key_id").SetValue(*newCredentials.AccessKeyId)
			cfg.Section(s.Name()).Key("aws_secret_access_key").SetValue(*newCredentials.SecretAccessKey)
			break
		}
	}

	logrus.Infof("Saving new credential for %s", c.Profile)
	if err := cfg.SaveTo(constants.AWSCredentialsPath); err != nil {
		return err
	}
	logrus.Infof("New credential for %s is successfully changed", c.Profile)

	if !isDeleted {
		logrus.Infof("Deleting old credential: %s", currentAccessKey)
		if err := r.AWSClient.DeleteAccessKey(currentAccessKey, c.Name); err != nil {
			return err
		}
		logrus.Info("Old credential is successfully deleted")
	}

	color.Green.Fprintln(out, "renew credentials are successfully done")

	return nil
}

// EcrLogin returns authorization data for ecr-login
func (r Runner) EcrLogin(out io.Writer) error {
	data, err := r.AWSClient.GetAuthorizeToken()
	if err != nil {
		return err
	}

	cmd, err := createEcrLoginCommand(*data.AuthorizationToken, *data.ProxyEndpoint)
	if err != nil {
		return err
	}

	if err := tools.CopyToClipBoard(cmd); err != nil {
		return err
	}

	color.Blue.Fprintf(out, "Token is copied to clipboard. Please paste it to terminal.")

	return nil
}

// AskBaseAccountName asks user's base account
func AskBaseAccountName() (string, error) {
	var name string
	prompt := &survey.Input{
		Message: "Your base account(company email): ",
	}
	survey.AskOne(prompt, &name)

	if len(name) == 0 {
		return name, errors.New("input base account has been canceled")
	}

	return name, nil
}

// AskAssumeTarget asks assume target
func AskAssumeTarget(assumeList map[string]string) (string, error) {
	var target string

	keys := tools.GetKeys(assumeList)

	sort.Strings(keys)

	prompt := &survey.Select{
		Message: "Choose the environment: ",
		Options: keys,
	}
	survey.AskOne(prompt, &target)

	if len(target) == 0 {
		return target, errors.New("get assume target has been canceled")
	}

	return target, nil
}

// printInitFile prints expected init file
func printInitFile(b []byte) error {
	_, err := fmt.Println(string(b))
	return err
}

// checkTarget checks if target is in the list
func CheckTarget(arn, target string) error {
	if len(arn) > 0 {
		return nil
	}

	return fmt.Errorf("%s is not registered in the assume list", target)
}

// PrintWebACL prints information
func PrintWebACL(out io.Writer, info *schema.WebACL) error {
	var data = struct {
		Summary *schema.WebACL
	}{
		Summary: info,
	}

	funcMap := template.FuncMap{
		"decorate": color.DecorateAttr,
	}

	w := tabwriter.NewWriter(out, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("Describe Web ACL").Funcs(funcMap).Parse(templates.WafTemplate))

	err := t.Execute(w, data)
	if err != nil {
		return err
	}
	return w.Flush()
}

// ParseTargetList parses arguments to valid ip list
func ParseTargetList(args []string) ([]string, error) {
	var ret []string
	for _, ip := range args {
		if err := IsValidAddress(ip); err != nil {
			return nil, err
		}

		if len(strings.Split(ip, "/")) == 1 {
			ip += "/32"
		}

		ret = append(ret, ip)
	}

	return ret, nil
}

// IsValidAddress checks if address is valid or not
func IsValidAddress(ip string) error {
	if strings.Count(ip, ".") != 3 {
		return fmt.Errorf("wrong IP address: %s", ip)
	}

	split := strings.Split(ip, "/")
	if len(split) > 2 {
		return fmt.Errorf("address cannot has / more than once: %s", ip)
	}

	if len(split) == 2 {
		base, err := strconv.Atoi(strings.TrimSpace(split[1]))
		if err != nil {
			return err
		}

		if base < 0 || base > 32 {
			return fmt.Errorf("cidr base should be between 0 and 32: %s", ip)
		}
	}

	classes := strings.Split(split[0], ".")
	for _, c := range classes {
		n, err := strconv.Atoi(c)
		if err != nil {
			return err
		}

		if n < 0 || n > 255 {
			return fmt.Errorf("each class number should be between 0 and 255: %s", ip)
		}
	}

	return nil
}

// SelectTargetACL makes a user choose ACL from the list
func (r Runner) SelectTargetACL(args []string) (string, error) {
	var target string
	var err error

	if len(args) == 0 {
		target, err = r.AWSClient.SelectACL()
		if err != nil {
			return constants.EmptyString, err
		}
	} else {
		target = args[0]
	}

	return target, err
}

// CheckIfTargetExistsInACL checks if target exists in ACL
func (r Runner) CheckIfTargetExistsInACL(targetList []string, acl string) ([]schema.IPCheckResult, error) {
	info, err := r.AWSClient.DescribeWebACL(acl)
	if err != nil {
		return nil, err
	}

	ret := []schema.IPCheckResult{}
	for _, t := range targetList {
		tmp := schema.IPCheckResult{
			IP:     t,
			Result: false,
		}
		if info.Rules != nil && len(info.Rules) > 0 {
			for _, rule := range info.Rules {
				isFinished := false
				for _, ds := range rule.IPDataSet {
					if tools.IsStringInArray(t, ds.IPList) {
						tmp.IPSetID = ds.ID
						tmp.Result = true
					}

					if tmp.Result {
						isFinished = true
						break
					}
				}
				if isFinished {
					break
				}
			}
		}
		ret = append(ret, tmp)
	}

	return ret, err
}

// PrintHasIPResult prints search result
func PrintHasIPResult(out io.Writer, result []schema.IPCheckResult) error {
	var data = struct {
		Summary []schema.IPCheckResult
	}{
		Summary: result,
	}

	funcMap := template.FuncMap{
		"decorate": color.DecorateAttr,
	}

	w := tabwriter.NewWriter(out, 0, 5, 3, ' ', tabwriter.TabIndent)
	t := template.Must(template.New("Search result").Funcs(funcMap).Parse(templates.IPSearchResultTemplate))

	err := t.Execute(w, data)
	if err != nil {
		return err
	}
	return w.Flush()
}
