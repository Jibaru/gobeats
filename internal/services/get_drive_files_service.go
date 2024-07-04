package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jibaru/gobeats/internal/entities"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	getFilesUrl           = "https://www.googleapis.com/drive/v3/files"
	googleDriveFolderType = "application/vnd.google-apps.folder"
	googleDriveAudioType  = "audio"
)

type GetDriveFilesService struct {
	rootDriveFolderKey string
	driveApiKey        string
}

func NewGetDriveFilesService(
	rootDriveFolderKey string,
	driveApiKey string,
) *GetDriveFilesService {
	return &GetDriveFilesService{
		rootDriveFolderKey,
		driveApiKey,
	}
}

func (serv *GetDriveFilesService) Do() (entities.DriveFileList, error) {
	return serv.fetchFilesInFolder(serv.rootDriveFolderKey)
}

func (serv *GetDriveFilesService) fetchFilesInFolder(folderID string) (entities.DriveFileList, error) {
	url := getFilesUrl + "?q=%22" + folderID + "%22%20in%20parents&key=" + serv.driveApiKey + "&pageSize=500"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files struct {
		Files entities.DriveFileList `json:"files"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, err
	}

	if files.Error.Message != "" {
		fmt.Println(files.Error.Message)
		return nil, errors.New("cannot fetch data")
	}

	var filesInFolder entities.DriveFileList

	for _, file := range files.Files {
		if strings.Contains(file.MimeType, googleDriveAudioType) && filepath.Ext(file.Name) == ".mp3" {
			filesInFolder = append(
				filesInFolder,
				entities.DriveFile{
					Name:      serv.removeExtension(file.Name),
					ID:        file.ID,
					MimeType:  file.MimeType,
					Extension: filepath.Ext(file.Name),
				},
			)
		} else if file.MimeType == googleDriveFolderType {
			filesInSubFolder, err := serv.fetchFilesInFolder(file.ID)
			if err != nil {
				return nil, err
			}

			filesInFolder = append(filesInFolder, filesInSubFolder...)
		}
	}

	return filesInFolder, nil
}

func (serv *GetDriveFilesService) removeExtension(str string) string {
	dotIndex := strings.LastIndex(str, ".")
	if dotIndex == -1 {
		return str
	}
	return str[:dotIndex]
}

func (serv *GetDriveFilesService) getExtension(str string) string {
	dotIndex := strings.LastIndex(str, ".")
	if dotIndex == -1 {
		return str
	}
	return str[:dotIndex]
}
