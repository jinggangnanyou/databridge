package minio

import (
	"context"
	"fmt"

	"databridge/common"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel"
)

type MinioConfig struct {
	BucketName      string `yaml:"bucket_name"`
	Region          string `yaml:"region"`
	Enabled         bool   `yaml:"enabled"`
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	UseSSL          bool   `yaml:"use_ssl"`
}

func InitMinioClient(minioConfig *MinioConfig) (*minio.Client, error) {
	tracer := otel.Tracer(common.ModuleName)
	_, span := tracer.Start(context.Background(), "init minio")
	fmt.Printf("trace_id:%s,span_id:%s\n",
		span.SpanContext().TraceID(), span.SpanContext().SpanID())
	defer span.End()
	minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
		Secure: minioConfig.UseSSL,
		Region: minioConfig.Region,
	})
	if err != nil {
		return nil, err
	}
	exists, errBucketExists := minioClient.BucketExists(context.Background(), minioConfig.BucketName)
	if errBucketExists != nil {
		return nil, errBucketExists
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), minioConfig.BucketName, minio.MakeBucketOptions{Region: minioConfig.Region})
		if err != nil {
			return nil, err
		}
	}
	return minioClient, nil
}
