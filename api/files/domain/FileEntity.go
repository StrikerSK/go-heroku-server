package fileDomains

type FileEntity struct {
	Id       uint   `json:"-"`
	FileData []byte `json:"-"`
}
