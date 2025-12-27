package dynamo

import (
	"context"
	"fmt"

	appconfig "github.com/HariPrasath-3/scheduler-service/pkg/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoClient(
	ctx context.Context,
	cfg *appconfig.DynamoConfig,
) (*dynamodb.Client, error) {

	if cfg.Region == "" {
		return nil, fmt.Errorf("dynamo region must be set")
	}

	// Base AWS config
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, err
	}

	// DynamoDB Local support
	if cfg.Endpoint != "" {
		awsCfg.EndpointResolverWithOptions =
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					if service == dynamodb.ServiceID {
						return aws.Endpoint{
							URL:           cfg.Endpoint,
							SigningRegion: cfg.Region,
						}, nil
					}
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				},
			)

		// Dummy credentials required for DynamoDB Local
		awsCfg.Credentials = credentials.NewStaticCredentialsProvider(
			"dummy",
			"dummy",
			"",
		)
	}

	return dynamodb.NewFromConfig(awsCfg), nil
}
