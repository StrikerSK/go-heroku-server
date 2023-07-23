package fileDomains

type FileEntityV2 struct {
	Id       string `json:"-"`
	FileData []byte `json:"-"`
}
