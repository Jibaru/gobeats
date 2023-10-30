package entities

import (
	"bytes"
	"errors"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type DrivePlayer struct {
	context       *oto.Context
	player        *oto.Player
	mutex         sync.Mutex
	CurrentSong   *DriveFile
	CurrentVolume int
	isDownloading bool
	isCanceled    bool
	isPaused      bool
}

func NewDrivePlayer(
	context *oto.Context,
	initialVolume int,
) (*DrivePlayer, error) {
	if initialVolume > 100 || initialVolume < 0 {
		return nil, errors.New("volume must be betweeen 0 and 100")
	}

	return &DrivePlayer{
		context,
		nil,
		sync.Mutex{},
		nil,
		initialVolume,
		false,
		false,
		true,
	}, nil
}

func (p *DrivePlayer) DownloadAndPlay(
	driveFile DriveFile,
	onPlaying func(duration int64, song DriveFile),
	onFinish func(),
) error {
	p.Cancel()

	p.CurrentSong = &driveFile

	var fileBytes []byte
	var err error

	if driveFile.IsCached() {
		fileBytes, err = os.ReadFile(driveFile.CachedFilename())
		if err != nil {
			return err
		}
	} else {
		if err = p.downloadMP3(driveFile.Url(), driveFile.CachedFilename()); err != nil {
			return err
		}

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

	p.player = p.context.NewPlayer(decodedMp3)
	p.player.SetVolume(p.VolumeAsOne())
	p.player.Play()
	p.UnCancel()
	p.isPaused = false

	var currentTime int64 = 0
	go func() {
		for {
			p.mutex.Lock()
			if p.isCanceled {
				p.mutex.Unlock()
				return
			}
			p.mutex.Unlock()

			if !p.IsFinished() {
				if !p.isPaused {
					currentTime++
					onPlaying(currentTime, driveFile)
					time.Sleep(1 * time.Second)
				}
			} else {
				onFinish()
				return
			}
		}
	}()

	return nil
}

func (p *DrivePlayer) VolumeAsOne() float64 {
	return float64(p.CurrentVolume) / 100
}

func (p *DrivePlayer) Pause() {
	if p.player != nil && p.player.IsPlaying() {
		p.player.Pause()
		p.isPaused = true
	}
}

func (p *DrivePlayer) Resume() {
	if p.player != nil && !p.player.IsPlaying() {
		p.player.Play()
		p.isPaused = false
	}
}

func (p *DrivePlayer) Cancel() {
	p.mutex.Lock()
	p.isCanceled = true
	p.mutex.Unlock()
	time.Sleep(1 * time.Second)
}

func (p *DrivePlayer) UnCancel() {
	p.mutex.Lock()
	p.isCanceled = false
	p.mutex.Unlock()
}

func (p *DrivePlayer) Stop() error {
	if p.player != nil {
		err := p.player.Close()
		if err != nil {
			p.player = nil
			return err
		}
	}

	p.player = nil
	return nil
}

func (p *DrivePlayer) IsFinished() bool {
	return (p.player != nil && !p.player.IsPlaying() && !p.isPaused) || p.player == nil
}

func (p *DrivePlayer) IncreaseVolume() {
	if p.CurrentVolume+1 <= 100 {
		p.CurrentVolume += 1
		if p.player != nil {
			p.player.SetVolume(p.VolumeAsOne())
		}
	}
}

func (p *DrivePlayer) DecreaseVolume() {
	if p.CurrentVolume-1 >= 0 {
		p.CurrentVolume -= 1
		if p.player != nil {
			p.player.SetVolume(p.VolumeAsOne())
		}
	}
}

func (p *DrivePlayer) downloadMP3(url, filename string) error {
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
