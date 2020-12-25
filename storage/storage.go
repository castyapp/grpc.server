package storage

import (
	"github.com/CastyLab/grpc.server/config"
	"github.com/minio/minio-go"
)

var Client *minio.Client

func Configure() (err error) {
	Client, err = minio.NewV4(
		config.Map.Secrets.ObjectStorage.Endpoint,
		config.Map.Secrets.ObjectStorage.AccessKey,
		config.Map.Secrets.ObjectStorage.SecretKey,
		false,
	)
	if err != nil {
		return err
	}
	return nil
}
