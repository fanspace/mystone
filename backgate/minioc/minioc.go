package minioc

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "enlist_appapi/logger"
	"enlist_appapi/settings"
)

func GetMinioConn()(*minio.Client, error) {
	minioClient, err := minio.New(settings.Cfg.MinioSettings.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.Cfg.MinioSettings.KeyID, settings.Cfg.MinioSettings.AccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return minioClient, err
}
