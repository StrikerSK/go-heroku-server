package fileClient

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

var testClient = NewFileClient()

func Test_UploadingFile(t *testing.T) {
	// Open the file
	file, _ := os.Open("./Test.json")
	attachmentID, err := testClient.uploadAttachment(file)
	assert.Nil(t, err, "There should be no error during attachment upload")
	assert.NotEmpty(t, attachmentID, "attachment should have id assigned")
}

type DeletingSuite struct {
	suite.Suite
	FileClient   FileClient
	AttachmentID string
}

func (suite *DeletingSuite) SetupTest() {
	suite.FileClient = NewFileClient()
	file, _ := os.Open("./Test.json")
	id, _ := suite.FileClient.uploadAttachment(file)
	suite.AttachmentID = id
}

func (suite *DeletingSuite) TestAttachmentDeletion() {
	err := suite.FileClient.deleteAttachment(suite.AttachmentID)
	suite.Nil(err, "There should be no error during attachment delete")
}

type ReadingSuite struct {
	suite.Suite
	FileClient   FileClient
	AttachmentID string
}

func (suite *ReadingSuite) SetupTest() {
	suite.FileClient = NewFileClient()
	file, _ := os.Open("./Test.json")
	id, _ := suite.FileClient.uploadAttachment(file)
	suite.AttachmentID = id
}

func (suite *ReadingSuite) TestAttachmentReading() {
	attachment, err := suite.FileClient.readAttachment(suite.AttachmentID)
	suite.Nil(err, "There should be no error during attachment read")
	suite.NotEmpty(attachment, "attachment should not be empty")

	var testMap map[string]string
	_ = json.Unmarshal(attachment, &testMap)
	suite.Equal("tester", testMap["firstName"])
}

func TestFileSuite(t *testing.T) {
	suite.Run(t, new(DeletingSuite))
	suite.Run(t, new(ReadingSuite))
}
