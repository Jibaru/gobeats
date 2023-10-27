package main

import (
	"fmt"
	"github.com/ebitengine/oto/v3"
	"github.com/jibaru/gobeats/m/internal/entities"
	"github.com/jibaru/gobeats/m/internal/services"
	"github.com/jibaru/gobeats/m/internal/ui"
	"github.com/spf13/viper"
	"strconv"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	rootFolderKey := viper.GetString("google_drive.root_folder_key")
	driveApiKey := viper.GetString("google_drive.api_key")

	cachedDriveFilesService := services.NewCachedDriveFilesService()
	getDriveFilesService := services.NewGetDriveFilesService(
		rootFolderKey,
		driveApiKey,
	)

	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-readyChan
	player := entities.NewDrivePlayer(otoCtx)

	cmd, err := ui.NewCmdUserInterface(
		// on close
		func(cmd *ui.CmdUserInterface) {
			err = player.Close()
			if err != nil {
				cmd.ChangeStatus(err.Error())
			}
		},
		// on enter pressed
		func(cmd *ui.CmdUserInterface) {
			err := player.Stop()
			if err != nil {
				cmd.ChangeStatus(err.Error())
			}

			selectedSong := cmd.GetSelectedSong()
			cmd.ClearPlayerStatus()
			cmd.ChangeStatus("Downloading " + selectedSong.Name)

			err = player.DownloadAndPlay(selectedSong)
			if err != nil {
				cmd.ChangeStatus(err.Error())
			} else {
				cmd.ChangePlayerStatus("Playing " + selectedSong.Name)
			}
		},
		// on pause pressed
		func(cmd *ui.CmdUserInterface) {
			player.Pause()
			selectedSong := cmd.GetSelectedSong()
			cmd.ChangePlayerStatus("Paused " + selectedSong.Name)
		},
		// on resume pressed
		func(cmd *ui.CmdUserInterface) {
			player.Resume()
			selectedSong := cmd.GetSelectedSong()
			cmd.ChangePlayerStatus("Playing " + selectedSong.Name)
		},
		// on increase volume pressed
		func(cmd *ui.CmdUserInterface) {
			player.IncreaseVolume()
			cmd.ChangeStatus("Volume: " + strconv.Itoa(player.CurrentVolume))
		},
		// on decrease volume pressed
		func(cmd *ui.CmdUserInterface) {
			player.DecreaseVolume()
			cmd.ChangeStatus("Volume: " + strconv.Itoa(player.CurrentVolume))
		},
	)
	if err != nil {
		return
	}
	defer cmd.Close()

	cmd.ChangeStatus("Reading cached songs...")

	driveFiles, err := cachedDriveFilesService.Get()
	if err != nil {
		cmd.ChangeStatus("Cached songs not found, fetching...")

		driveFiles, err = getDriveFilesService.Do()
		if err != nil {
			return
		}

		cmd.ChangeStatus("Caching songs fetched, caching...")

		err := cachedDriveFilesService.Set(driveFiles)
		if err != nil {
			return
		}
	}

	cmd.SetSongList(driveFiles)
	cmd.ChangeStatus("Ready to play!")

	cmd.Loop()
}
