package fileClient

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_UploadingFile(t *testing.T) {
	client := NewFileClient()
	// Open the file
	file, _ := os.Open("./Test.json")
	id, err := client.uploadAttachment(file)
	assert.Nil(t, err, "There should be no error during attachment upload")
	assert.NotEmpty(t, id, "attachment should have id assigned")
}

func Test_DeletingFile(t *testing.T) {
	client := NewFileClient()
	err := client.deleteAttachment("409c4410-4db6-477a-9cec-06a25c466775")
	assert.Nil(t, err, "There should be no error during attachment upload")
}
