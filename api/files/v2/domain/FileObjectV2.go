package fileDomains

type FileObjectV2 struct {
	FileEntityV2
	FileMetadataV2
}

func (receiver *FileObjectV2) SetID(uuid string) {
	receiver.FileEntityV2.Id = uuid
	receiver.FileMetadataV2.Id = uuid
}
