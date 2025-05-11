package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDriveFilesService_Do(t *testing.T) {
	rootFolderKey := "put your root folder key here"
	driveApiKey := "put your drive api key here"

	service := NewGetDriveFilesService(
		rootFolderKey,
		driveApiKey,
	)

	files, err := service.Do()

	assert.Nil(t, err)
	assert.NotNil(t, files)
	fmt.Printf("%+v\n", files)
}
