package entities

import "os"

type DriveFile struct {
	Name      string `json:"name"`
	MimeType  string `json:"mimeType"`
	ID        string `json:"id"`
	Extension string `json:"extension"`
}

type DriveFileList = []DriveFile

func (d *DriveFile) Url() string {
	return "https://docs.google.com/uc" + "?export=download&id=" + d.ID
}

func (d *DriveFile) CachedFilename() string {
	return "./storage/" + d.ID + "." + d.Extension
}

func (d *DriveFile) IsCached() bool {
	if _, err := os.Stat(d.CachedFilename()); err == nil {
		return true
	}

	return false
}
