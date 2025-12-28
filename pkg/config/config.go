package config

type AppConfig struct {
	Grpc                  GrpcConfig            `yaml:"grpc"`
	Redis                 RedisConfig           `yaml:"redis"`
	Kafka                 KafkaConfig           `yaml:"kafka"`
	Dynamo                DynamoConfig          `yaml:"dynamo"`
	SchedulerConfig       SchedulerConfig       `yaml:"scheduler"`
	SchedulerWorkerConfig SchedulerWorkerConfig `yaml:"scheduler_worker"`
	RetryWorkerConfig     RetryWorkerConfig     `yaml:"retry_worker"`
}

type GrpcConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type RedisMode string

const (
	RedisStandalone RedisMode = "standalone"
	RedisCluster    RedisMode = "cluster"
)

type RedisConfig struct {
	Mode                          RedisMode `yaml:"mode"`
	Hosts                         []string  `yaml:"hosts"`
	Password                      string    `yaml:"password"`
	Cluster                       bool      `yaml:"cluster"`
	ServeReadsFromSlaves          bool      `yaml:"serve_reads_from_slaves"`
	ServeReadsFromMasterAndSlaves bool      `yaml:"serve_reads_from_master_and_slaves"`
	ReadTimeout                   int       `yaml:"read_timeout"`
	WriteTimeout                  int       `yaml:"write_timeout"`
	IdleTimeout                   int       `yaml:"idle_timeout"`
	DialTimeout                   int       `yaml:"dial_timeout"`
	PoolSize                      int       `yaml:"pool_size"`
	MinIdleConns                  int       `yaml:"min_idle_conns"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
}

type DynamoConfig struct {
	Region    string `yaml:"region"`
	Endpoint  string `yaml:"endpoint"`
	TableName string `yaml:"table_name"`
}

type SchedulerConfig struct {
	TotalPartitions int `yaml:"total_partitions"`
	BucketSizeSec   int `yaml:"bucket_size_sec"`
}

type SchedulerWorkerConfig struct {
	PastBucketsCount int `yaml:"past_buckets_count"`
	SemaphoreLimit   int `yaml:"semaphore_limit"`
	BatchSize        int `yaml:"batch_size"`
	BackoffMs        int `yaml:"backoff_ms"`
	MaxBackoffMs     int `yaml:"max_backoff_ms"`
}

type RetryWorkerConfig struct {
	PastBucketsCount    int `yaml:"past_buckets_count"`
	RetryScanIntervalMs int `yaml:"retry_scan_interval_ms"`
	RetryBatchSize      int `yaml:"retry_batch_size"`
}
