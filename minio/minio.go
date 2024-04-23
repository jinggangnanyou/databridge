package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	//BucketName BucketName
	BucketName = "code"
	//Region Region
	Region = "cn-north-1"
)

type MinioConfig struct {
	Enabled         bool   `yaml:"enabled"`
	EndPoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	UseSSL          bool   `yaml:"use_ssl"`
}

func InitMinioClient(minioConfig *MinioConfig) (*minio.Client, error) {
	minioClient, err := minio.New(minioConfig.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
		Secure: minioConfig.UseSSL,
		Region: Region,
	})
	if err != nil {
		return nil, err
	}
	exists, errBucketExists := minioClient.BucketExists(context.Background(), BucketName)
	if errBucketExists != nil {
		return nil, errBucketExists
	}
	if !exists {
		err = minioClient.MakeBucket(context.Background(), BucketName, minio.MakeBucketOptions{Region: Region})
		if err != nil {
			return nil, err
		}
	}
	return minioClient, nil
}
