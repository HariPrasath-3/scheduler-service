package env

import (
	"github.com/HariPrasath-3/scheduler-service/pkg/config"
	"github.com/HariPrasath-3/scheduler-service/pkg/kafka"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/redis/go-redis/v9"
)

type Env struct {
	config   *config.AppConfig
	redis    *redis.ClusterClient
	dynamo   *dynamodb.Client
	producer *kafka.Producer
}

func NewEnv(options ...func(env *Env)) *Env {
	env := &Env{}
	for _, option := range options {
		option(env)
	}
	return env
}

func WithAppConfig(cfg *config.AppConfig) func(*Env) {
	return func(env *Env) {
		env.config = cfg
	}
}

func WithRedisClient(redis *redis.ClusterClient) func(*Env) {
	return func(env *Env) {
		env.redis = redis
	}
}

func WithDynamoClient(dynamo *dynamodb.Client) func(*Env) {
	return func(env *Env) {
		env.dynamo = dynamo
	}
}

func WithKafkaProducer(producer *kafka.Producer) func(*Env) {
	return func(env *Env) {
		env.producer = producer
	}
}

func (e *Env) Config() *config.AppConfig {
	return e.config
}

func (e *Env) Redis() *redis.ClusterClient {
	return e.redis
}

func (e *Env) Dynamo() *dynamodb.Client {
	return e.dynamo
}

func (e *Env) Producer() *kafka.Producer {
	return e.producer
}
