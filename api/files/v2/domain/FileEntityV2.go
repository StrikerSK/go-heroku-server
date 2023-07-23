package fileDomains

type FileEntityV2 struct {
	Id       uint   `json:"-"`
	FileData []byte `json:"-"`
}
