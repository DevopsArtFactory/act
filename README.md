# act
act is a command line tool for developers who are using AWS Cloud. This command line interface will help you do work more easily and fast.
You can assume to the different AWS account or generate rds authentication token with one command.

## Important Notice
- For almost of commands, act will check whether the access key is expired or not (180 days)
- If your access key is older than `180 days`, run `act renew-credential`.
  - It will change your `aws_access_key_id` and `aws_secret_access_key` in $HOME/.aws/credentials.

## Install
* macOS user
```bash
$ brew tap DevopsArtFactory/devopsart
$ brew install act
$ act version -v info
```

* Linux user
```bash
$ curl -Lo act https://devopsartfactory.s3.ap-northeast-2.amazonaws.com/act/releases/latest/act-linux-amd64
$ sudo install act /usr/local/bin/
$ act version -v info
```

* Windows user
  - file : https://devopsartfactory.s3.ap-northeast-2.amazonaws.com/act/releases/latest/act-windows-amd64.exe
  - Simply download it and place it in your PATH as act.exe.
  
  
 ## Auto completion
- zsh 
  - This is recommended.
 ```bash
$ echo "source <(act completion zsh)" >> ~/.zshrc
$ source  ~/.zshrc
```

- bash 
 ```bash
$ echo "source <(act completion bash)" >> ~/.bash_rc or ~/.bash_profile
$ source  ~/.bashrc
```

## Setting Configuration
- Configuration file should be in `$HOME/.aws/config.yaml`.
- You can easily create your configuration file with `act init`
- `act init` will create a configuration for default profile.
```bash
$ act init
? Your base account:  gslee@gmail.com
- profile: default
  name: gslee@gmail.com
  duration: 3600
  assume_roles:
    preprod: ""
    prod: ""
  databases:
    preprod:
    - ""
    prod:
    - ""

? Are you sure to generate configuration file?  yes
New configuration file is successfully generated in $HOME//Users/gslee/.aws/config.yaml
```

## Alias for assume role
- You can set alias with alias list.
- **You cannot use `-` prefix for alias because golang will detect it as flag.**
```bash
$ vim ~/.aws/config.yaml
- profile: default
  name: gslee@gmail.com
  alias:
    d: preprod
    p: prod
  assume_roles:
    preprod: arn:aws:iam::xxxxxxxxxxxx:role/userassume-benx-preprod-admin
    prod: arn:aws:iam::xxxxxxxxxxx:role/userassume-benx-prod-admin

# This is equal to `act setup preprod`
$ act setup d
Assume Credentials copied to clipboard, please paste it.
```

## RDS IAM Authentication
- You can get RDS auth token in order to log in to the database.
- If you follow these steps, then you will get your token in the clipboard

```bash
$ act get rds-token
[preprod-aurora prod]
? Choose the environment:   [Use arrows to move, type to filter]
> preprod-aurora
  prod

# If you select environment
[preprod-aurora prod]
? Choose the environment:  preprod-aurora
? Choose an instance: wet  [Use arrows to move, type to filter]
> xxxxxxxxxxxxxxx.cluster-xxxxxxx.ap-northeast-2.rds.amazonaws.com
  yyyyyyyyyyyyyyy.cluster-yyyyyyy.ap-northeast-2.rds.amazonaws.com
  zzzzzzzzzzzzzzz.cluster-zzzzzzz.ap-northeast-2.rds.amazonaws.com

# If you choose instance
$ act get rds-token
[preprod-aurora prod]
? Choose the environment:  preprod-aurora
? Choose an instance: xxxxxxxxxxxxxx.cluster-xxxxxxx.ap-northeast-2.rds.amazonaws.com
Assume Role MFA token code: 712352
INFO[0084] Token is copied to clipboard.
```


## Commands 
```bash
AWS command line helper tool

managing configuration of act
  init             initialize act command line tool

commands related to aws IAM credentials
  renew-credential recreates aws credential of profile

commands for controlling assume role
  setup            create assume credentials for multi-account
  who              check the account information of current shell

commands for retrieving information related to AWS WAF.
  describe-web-acl retrieve detailed list of web acl
  has-ip           check if ip is registered in the web acl

Other Commands:
  assume           do work about assume role
  completion       Output shell completion for the given shell (bash or zsh)
  ecr-login        login to ECR
  get              Get token or information with act
  version          Print the version information

Usage:
  act [flags] [options]

Use "act <command> --help" for more information about a given command.
```

## Contribution Guide
- Check [CONTRIBUTING.md](CONTRIBUTING.md)
