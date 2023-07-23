package filePorts

import (
	"go-heroku-server/api/files/v2/domain"
)

type IFileRepositoryV2 interface {
	CreateFile(fileDomains.FileEntityV2) error
	DeleteFile(uint) error
}
