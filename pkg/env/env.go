package env

import (
	"github.com/HariPrasath-3/scheduler-service/pkg/kafka"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/redis/go-redis/v9"
)

const ctxKeyEnv = "env"

type Env struct {
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

// func UnaryServerInterceptor(env *Env) grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		nctx := env.WithContext(ctx)
// 		return handler(nctx, req)
// 	}
// }

// func KafkaHandlerWithEnv(e *Env, next kafka.MessageHandler) kafka.MessageHandler {
// 	return kafka.MessageHandlerFunc(func(ctx context.Context, msg *sarama.ConsumerMessage) error {
// 		nctx := e.WithContext(ctx)
// 		return next.HandleMessage(nctx, msg)
// 	})
// }

// func (env *Env) WithContext(ctx context.Context) context.Context {
// 	nctx := context.WithValue(ctx, ctxKeyEnv, env)
// 	return nctx
// }

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

func (e *Env) Redis() *redis.ClusterClient {
	return e.redis
}

func (e *Env) Dynamo() *dynamodb.Client {
	return e.dynamo
}

func (e *Env) Producer() *kafka.Producer {
	return e.producer
}
