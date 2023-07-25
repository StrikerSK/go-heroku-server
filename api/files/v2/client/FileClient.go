package client

import fileDomains "go-heroku-server/api/files/v2/domain"

type FileClient struct {
	baseURL string
}

func NewFileClient() FileClient {
	return FileClient{
		baseURL: "localhost:8080",
	}
}

func (r FileClient) uploadAttachment(attachment fileDomains.FileObjectV2) error {
	return nil
}
