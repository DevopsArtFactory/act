package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds"
)

func GetRDSClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *rds.RDS {
	if creds == nil {
		return rds.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return rds.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

func (c Client) StartLoadtestRDS(dbClusterID string) (*string, error) {
	input := &rds.StartDBClusterInput{
		DBClusterIdentifier: aws.String(dbClusterID),
	}

	result, err := c.RDSClient.StartDBCluster(input)
	if err != nil {
		return nil, err
	}

	return result.DBCluster.Status, nil
}

func (c Client) StopLoadtestRDS(dbClusterID string) (*string, error) {
	input := &rds.StopDBClusterInput{
		DBClusterIdentifier: aws.String(dbClusterID),
	}

	result, err := c.RDSClient.StopDBCluster(input)
	if err != nil {
		return nil, err
	}

	return result.DBCluster.Status, nil
}

func (c Client) GetLoadtestRDSStatus(dbClusterID string) (*string, error) {
	input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(dbClusterID),
	}

	result, err := c.RDSClient.DescribeDBClusters(input)
	if err != nil {
		return nil, err
	}

	return result.DBClusters[0].Status, nil
}
