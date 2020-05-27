package utils

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

// SaveUploadedFile to save the file contents in temp location
func SaveUploadedFile(fileToSave *multipart.File, UploadFileHeaders *multipart.FileHeader) (string, error) {
	fileNameSplit := strings.Split(UploadFileHeaders.Filename, ".")
	tempFileName := fileNameSplit[0]
	tempFileExt := fileNameSplit[len(fileNameSplit)-1]

	tempFile, err := ioutil.TempFile("./temp", tempFileName+"-*."+tempFileExt)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(*fileToSave)
	if err != nil {
		return "", err
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	return tempFile.Name(), nil
}

// RemoveUploadedFile to remove the temporary file
func RemoveUploadedFile(fileName string) error {
	err := os.Remove(fileName)
	return err
}
