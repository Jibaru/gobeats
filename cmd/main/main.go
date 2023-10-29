package main

import (
	"fmt"
	"github.com/jibaru/gobeats/m/cmd/config"
	"github.com/jibaru/gobeats/m/internal/entities"
	"github.com/jibaru/gobeats/m/internal/services"
	"github.com/jibaru/gobeats/m/internal/ui"
	"github.com/jibaru/gobeats/m/pkg/time"
	"strconv"
)

func main() {
	cfg, err := config.NewAppConfig()
	if err != nil {
		fmt.Printf("Error reading cfg file: %s\n", err)
		return
	}

	cachedDriveFilesService := services.NewCachedDriveFilesService()
	getDriveFilesService := services.NewGetDriveFilesService(
		cfg.GoogleDriveRootFolderKey,
		cfg.GoogleDriveDriveApiKey,
	)

	otoCtx, readyChan, err := config.NewOtoContext()
	if err != nil {
		fmt.Println(err)
		return
	}
	<-readyChan

	player, err := entities.NewDrivePlayer(otoCtx, cfg.InitialVolume)
	if err != nil {
		fmt.Println(err)
		return
	}

	handlePlay := func(cmd *ui.CmdUserInterface) {
		cmd.ChangePlayerStatus("Stopping song...")
		err = player.Stop()
		if err != nil {
			cmd.ChangeStatus(err.Error())
		}

		cmd.ChangePlayerStatus("Loading song...")

		err = player.DownloadAndPlay(
			cmd.GetSelectedSong(),
			func(duration int64, song entities.DriveFile) {
				formattedDuration := time.FormatDurationInSeconds(duration)
				cmd.ChangePlayerStatus("[" + formattedDuration + "] Playing " + song.Name)
			},
			func() {
				cmd.ClearPlayerStatus()
				if cfg.UseAutoPlay {
					cmd.IncreaseSelectedSongIndex()
				}
			},
		)
		if err != nil {
			cmd.ChangeStatus(err.Error())
		}
	}

	cmd, err := ui.NewCmdUserInterface(
		// on close
		func(cmd *ui.CmdUserInterface) {
			cmd.ChangeStatus("Closing...")
		},
		// on enter pressed
		handlePlay,
		// on pause pressed
		func(cmd *ui.CmdUserInterface) {
			if player.CurrentSong != nil {
				player.Pause()
				cmd.ChangePlayerStatus("Paused " + player.CurrentSong.Name)
			}
		},
		// on resume pressed
		func(cmd *ui.CmdUserInterface) {
			if player.CurrentSong != nil {
				player.Resume()
				cmd.ChangePlayerStatus("Playing " + player.CurrentSong.Name)
			}
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
		// on increase songs index set
		handlePlay,
	)
	if err != nil {
		return
	}
	defer cmd.Close()

	cmd.ChangeStatus("Reading cached songs...")

	songs, err := cachedDriveFilesService.Get()
	if err != nil {
		cmd.ChangeStatus("Cached songs not found, fetching...")

		songs, err = getDriveFilesService.Do()
		if err != nil {
			return
		}

		cmd.ChangeStatus("Caching songs fetched, caching...")

		err = cachedDriveFilesService.Set(songs)
		if err != nil {
			return
		}
	}

	cmd.SetSongList(songs)
	cmd.ChangeStatus("Ready to play!")

	cmd.Loop()

	err = player.Stop()
	if err != nil {
		fmt.Println(err)
	}
}
