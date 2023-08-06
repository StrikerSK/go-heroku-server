package fileClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	userClient "go-heroku-server/api/user/client"
	"io"
	"net/http"
)

type FileClient struct {
	userClient.UserClient
	baseURL string
}

func NewFileClient() FileClient {
	baseURL := "http://localhost:8080"

	return FileClient{
		UserClient: userClient.NewUserClient(),
		baseURL:    baseURL,
	}
}

func (r FileClient) uploadAttachment(attachment io.Reader) (string, error) {
	fullUrl := r.baseURL + "/file/upload?name=Attachment"

	fileBytes, err := io.ReadAll(attachment)
	if err != nil {
		fmt.Println("error reading file:", err)
		return "", err
	}

	// Create the HTTP request with the file content in the body
	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewReader(fileBytes))
	if err != nil {
		fmt.Println("error creating request: ", err)
		return "", err
	}

	contentType := http.DetectContentType(fileBytes)

	// Set the appropriate content type header for the request
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", r.Token)

	// Make the request to upload the file
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error calling request: ", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body as a string
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response data:", err)
		return "", err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status code: %d", resp.StatusCode)
	}

	var mapStruct map[string]string
	err = json.Unmarshal(responseData, &mapStruct)
	if err != nil {
		fmt.Println("Error unmarshalling data:", err)
		return "", err
	}

	return mapStruct["id"], nil
}

func (r FileClient) deleteAttachment(attachmentID string) error {
	fullUrl := r.baseURL + "/file/" + attachmentID

	// Create the HTTP request with the file content in the body
	req, err := http.NewRequest(http.MethodDelete, fullUrl, nil)
	if err != nil {
		fmt.Println("error creating request: ", err)
		return err
	}
	req.Header.Set("Authorization", r.Token)

	// Make the request to upload the file
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error calling request: ", err)
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (r FileClient) readAttachment(attachmentID string) ([]byte, error) {
	fullUrl := r.baseURL + "/file/" + attachmentID

	// Create the HTTP request with the file content in the body
	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		fmt.Println("error creating request: ", err)
		return nil, err
	}
	req.Header.Set("Authorization", r.Token)

	// Make the request to upload the file
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error calling request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body as a string
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response data:", err)
		return nil, err
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("attachment read error: %s", responseData)
	}

	return responseData, nil
}
