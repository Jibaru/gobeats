package services

import (
	"encoding/json"
	"github.com/jibaru/gobeats/m/internal/entities"
	"os"
)

type CachedDriveFilesService struct {
}

func NewCachedDriveFilesService() *CachedDriveFilesService {
	return &CachedDriveFilesService{}
}

func (serv *CachedDriveFilesService) Get() ([]entities.DriveFile, error) {
	jsonData, err := os.ReadFile("./storage/drive_files.json")
	if err != nil {
		return nil, err
	}

	var driveFiles []entities.DriveFile

	err = json.Unmarshal(jsonData, &driveFiles)
	if err != nil {
		return nil, err
	}

	return driveFiles, nil
}

func (serv *CachedDriveFilesService) Set(driveFiles []entities.DriveFile) error {
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
