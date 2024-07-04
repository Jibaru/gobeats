package services

import (
	"encoding/json"
	"errors"
	"github.com/jibaru/gobeats/internal/entities"
	"os"
)

type CachedDriveFilesService struct {
}

func NewCachedDriveFilesService() *CachedDriveFilesService {
	return &CachedDriveFilesService{}
}

func (serv *CachedDriveFilesService) Get() (entities.DriveFileList, error) {
	jsonData, err := os.ReadFile("./storage/drive_files.json")
	if err != nil {
		return nil, err
	}

	var driveFiles entities.DriveFileList

	err = json.Unmarshal(jsonData, &driveFiles)
	if err != nil {
		return nil, err
	}

	if len(driveFiles) == 0 {
		return nil, errors.New("songs not found")
	}

	return driveFiles, nil
}

func (serv *CachedDriveFilesService) Set(driveFiles entities.DriveFileList) error {
	jsonData, err := json.Marshal(driveFiles)
	if err != nil {
		return err
	}

	err = os.WriteFile("./storage/drive_files.json", jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
