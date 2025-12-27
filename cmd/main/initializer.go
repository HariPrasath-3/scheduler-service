package main

import (
	"context"
	"log"
	"sync"

	"github.com/HariPrasath-3/scheduler-service/pkg/config"
	"github.com/HariPrasath-3/scheduler-service/pkg/dynamo"
	"github.com/HariPrasath-3/scheduler-service/pkg/env"
	"github.com/HariPrasath-3/scheduler-service/pkg/kafka"
	"github.com/HariPrasath-3/scheduler-service/pkg/redis"
	"github.com/HariPrasath-3/scheduler-service/pkg/shutdown"
)

type application struct {
	appConfig *config.AppConfig
	env       *env.Env

	shutdownCallbacksMx sync.Mutex
	shutdownOnce        sync.Once
	closed              chan struct{}
	shutdownCallbacks   []func()
	confCtx             context.Context
}

func (a *application) Close() error {

	a.shutdownOnce.Do(func() {
		func() {
			for _, cb := range a.shutdownCallbacks {
				defer cb()
			}
		}()
		close(a.closed)
	})
	// to handle for idempotence.
	// multiple Close will be blocked until application is actually closed
	a.WaitForClose()
	return nil
}

func (a *application) WaitForClose() {
	<-a.closed
}

func (a *application) AddConfCtx(ctx context.Context) {
	a.confCtx = ctx
}

func (a *application) AddShutdownCallback(cb func()) {
	a.shutdownCallbacksMx.Lock()
	defer a.shutdownCallbacksMx.Unlock()
	a.shutdownCallbacks = append(a.shutdownCallbacks, cb)
}

func initialize(app *application) {
	// 1️⃣ Root context
	ctx, cancel := context.WithCancel(context.Background())
	app.AddShutdownCallback(func() { cancel() })
	app.AddConfCtx(ctx)

	// 2️⃣ Load config
	cfg, err := config.Load(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	app.appConfig = cfg
	log.Println("application config loaded")

	// 3️⃣ Redis init (cluster)
	redisClient := redis.NewClient(&cfg.Redis)
	app.AddShutdownCallback(func() { redisClient.Close() })
	log.Println("redis client initialized")

	// 4️⃣ DynamoDB init
	dynamoClient, err := dynamo.NewDynamoClient(ctx, &cfg.Dynamo)
	if err != nil {
		log.Fatalf("failed to initialize dynamo client: %v", err)
	}
	log.Println("dynamo client initialized")

	// 5️⃣ Kafka Producer init
	producer, err := kafka.NewProducer(&cfg.Kafka)
	if err != nil {
		log.Fatalf("failed to initialize kafka producer: %v", err)
	}
	app.AddShutdownCallback(func() { producer.Close() })
	log.Println("kafka producer initialized")

	env := env.NewEnv(
		env.WithRedisClient(redisClient),
		env.WithDynamoClient(dynamoClient),
		env.WithKafkaProducer(producer),
	)
	app.env = env

	shutdown.Listen(func() {
		_ = app.Close()
	})

	log.Println("application initialized")
}
