package storage

import (
	"github.com/castyapp/grpc.server/config"
	"github.com/minio/minio-go"
)

var Client *minio.Client

func Configure(c *config.Map) (err error) {
	Client, err = minio.NewV4(c.S3.Endpoint, c.S3.AccessKey, c.S3.SecretKey, false)
	if err != nil {
		return err
	}
	return nil
}
