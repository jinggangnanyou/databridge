package kafka

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"databridge/common"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
)

const (
	defaultMessageMaxBytes = 20971520
)

type KafkaConfig struct {
	ClientID          string `yaml:"client_id"`
	MaxMessageBytes   int    `yaml:"max_message_bytes"`
	Host              string `yaml:"host"`
	SaslEnabled       bool   `yaml:"sasl_enabled"`
	SaslUsername      string `yaml:"sasl_username"`
	SaslPassword      string `yaml:"sasl_password"`
	NumPartitions     int32  `yaml:"num_partitions"`
	ReplicationFactor int16  `yaml:"replication_factor"`
	DialTimeout       int64  `yaml:"dial_timeout"` // 方便本地开发，连不上时快速启动
}

type AdminClient struct {
	Client            sarama.ClusterAdmin
	NumPartitions     int32
	ReplicationFactor int16
}

func (cf *KafkaConfig) newConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.ClientID = cf.ClientID
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.MaxMessageBytes = cf.MaxMessageBytes
	if cf.SaslEnabled {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cf.SaslUsername
		config.Net.SASL.Password = cf.SaslPassword
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}
	if cf.DialTimeout != 0 {
		config.Net.DialTimeout = time.Duration(cf.DialTimeout) * time.Second
	}
	return config
}

func (cf *KafkaConfig) newClusterAdmin() (AdminClient, error) {
	config := cf.newConfig()
	adminClient, err := sarama.NewClusterAdmin(strings.Split(cf.Host, ","), config)
	if err != nil {
		return AdminClient{}, err
	}
	return AdminClient{
		Client:            adminClient,
		NumPartitions:     cf.NumPartitions,
		ReplicationFactor: cf.ReplicationFactor,
	}, nil
}

func InitKafka(kafka *KafkaConfig, topicNames []string) error {
	tracer := otel.Tracer(common.ModuleName)
	_, span := tracer.Start(context.Background(), "init kafka")
	fmt.Printf("trace_id:%s,span_id:%s\n",
		span.SpanContext().TraceID(), span.SpanContext().SpanID())
	defer span.End()
	messageMaxBytes := defaultMessageMaxBytes // 20MB
	if kafka.MaxMessageBytes <= 0 {
		kafka.MaxMessageBytes = messageMaxBytes
	}
	messageMaxBytesStr := strconv.Itoa(kafka.MaxMessageBytes)
	clusterAdmin, err := kafka.newClusterAdmin()
	if err != nil {
		return err
	}
	defer clusterAdmin.Client.Close()

	clusterTopics, err := clusterAdmin.Client.ListTopics()
	if err != nil {
		return err
	}
	for _, topicName := range topicNames {
		if _, exists := clusterTopics[topicName]; exists {
			// alert config
			err = clusterAdmin.alertMaxMessageBytes(topicName, messageMaxBytesStr)
			if err != nil {
				fmt.Printf("AlertKafkaMaxMessageBytes err. topic: %s, error: %v", topicName, err)
			}
		} else {
			// create topic
			err = clusterAdmin.createTopic(topicName, messageMaxBytesStr)
			if err != nil {
				fmt.Printf("CreateTopic err. topic: %s, error: %v", topicName, err)
			}
		}
	}
	return nil
}

func (ac *AdminClient) alertMaxMessageBytes(topicName, messageMaxBytesStr string) error {
	configEntries := make(map[string]*string)
	configEntries["max.message.bytes"] = &messageMaxBytesStr
	return ac.Client.AlterConfig(sarama.TopicResource, topicName, configEntries, false)
}

func (ac *AdminClient) createTopic(topicName, messageMaxBytesStr string) error {
	configEntries := make(map[string]*string)
	configEntries["max.message.bytes"] = &messageMaxBytesStr
	detail := &sarama.TopicDetail{
		NumPartitions:     ac.NumPartitions,
		ReplicationFactor: ac.ReplicationFactor,
		ConfigEntries:     configEntries,
	}
	return ac.Client.CreateTopic(topicName, detail, false)
}
