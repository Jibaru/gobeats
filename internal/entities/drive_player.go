package entities

import (
	"bytes"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"io"
	"net/http"
	"os"
)

type DrivePlayer struct {
	Context       *oto.Context
	Player        *oto.Player
	IsDownloading bool
	CurrentVolume int
}

func NewDrivePlayer(context *oto.Context) *DrivePlayer {
	return &DrivePlayer{
		context,
		nil,
		false,
		50,
	}
}

func (p *DrivePlayer) DownloadAndPlay(driveFile DriveFile) error {
	var fileBytes []byte
	var err error

	if driveFile.IsCached() {
		fileBytes, err = os.ReadFile(driveFile.CachedFilename())
		if err != nil {
			return err
		}
	} else {
		p.IsDownloading = true
		if err = downloadMP3(driveFile.Url(), driveFile.CachedFilename()); err != nil {
			p.IsDownloading = false
			return err
		}
		p.IsDownloading = false

		fileBytes, err = os.ReadFile(driveFile.CachedFilename())
		if err != nil {
			return err
		}
	}
	fileBytesReader := bytes.NewReader(fileBytes)

	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		return err
	}

	err = p.Stop()
	if err != nil {
		return err
	}

	p.Player = p.Context.NewPlayer(decodedMp3)
	p.Player.SetVolume(p.VolumeAsOne())
	p.Player.Play()

	return nil
}

func (p *DrivePlayer) VolumeAsOne() float64 {
	return float64(p.CurrentVolume) / 100
}

func (p *DrivePlayer) Pause() {
	if p.Player != nil && p.Player.IsPlaying() {
		p.Player.Pause()
	}
}

func (p *DrivePlayer) Resume() {
	if p.Player != nil && !p.Player.IsPlaying() {
		p.Player.Play()
	}
}

func (p *DrivePlayer) Stop() error {
	if p.Player != nil {
		err := p.Player.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DrivePlayer) Close() error {
	return nil
	/*err := os.Remove("audio.mp3")
	if err != nil {
		return err
	}

	return nil*/
}

func (p *DrivePlayer) IncreaseVolume() {
	if p.CurrentVolume+1 <= 100 {
		p.CurrentVolume += 1
		if p.Player != nil {
			p.Player.SetVolume(p.VolumeAsOne())
		}
	}
}

func (p *DrivePlayer) DecreaseVolume() {
	if p.CurrentVolume-1 >= 0 {
		p.CurrentVolume -= 1
		if p.Player != nil {
			p.Player.SetVolume(p.VolumeAsOne())
		}
	}
}

func downloadMP3(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
