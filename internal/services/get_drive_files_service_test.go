package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDriveFilesService_Do(t *testing.T) {
	rootFolderKey := "1cbPWVPp-xFhNV2s540QHn8GuZ_Bwy9zM"
	driveApiKey := "AIzaSyBynS-g91QFAauy_0r1FN5hdQplNnCGsdM"

	service := NewGetDriveFilesService(
		rootFolderKey,
		driveApiKey,
	)

	files, err := service.Do()

	assert.Nil(t, err)
	assert.NotNil(t, files)
	fmt.Printf("%+v\n", files)
}
