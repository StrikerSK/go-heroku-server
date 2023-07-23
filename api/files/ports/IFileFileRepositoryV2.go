package filePorts

import fileDomains "go-heroku-server/api/files/domain"

type IFileRepositoryV2 interface {
	CreateFile(fileDomains.FileEntityV2) error
	DeleteFile(uint) error
}
