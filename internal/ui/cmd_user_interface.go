package ui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jibaru/gobeats/m/internal/entities"
)

type CmdUserInterface struct {
	driveFiles              []entities.DriveFile
	keyPressedWidget        *widgets.Paragraph
	songListWidget          *widgets.List
	playerStatusWidget      *widgets.Paragraph
	statusWidget            *widgets.Paragraph
	onClose                 func(cmd *CmdUserInterface)
	onEnterPressed          func(cmd *CmdUserInterface)
	onPausePressed          func(cmd *CmdUserInterface)
	onResumePressed         func(cmd *CmdUserInterface)
	onIncreaseVolumePressed func(cmd *CmdUserInterface)
	onDecreaseVolumePressed func(cmd *CmdUserInterface)
}

func NewCmdUserInterface(
	onClose func(cmd *CmdUserInterface),
	onEnterPressed func(cmd *CmdUserInterface),
	onPausePressed func(cmd *CmdUserInterface),
	onResumePressed func(cmd *CmdUserInterface),
	onIncreaseVolumePressed func(cmd *CmdUserInterface),
	onDecreaseVolumePressed func(cmd *CmdUserInterface),
) (*CmdUserInterface, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}

	cmd := &CmdUserInterface{
		onClose:                 onClose,
		onEnterPressed:          onEnterPressed,
		onPausePressed:          onPausePressed,
		onResumePressed:         onResumePressed,
		onIncreaseVolumePressed: onIncreaseVolumePressed,
		onDecreaseVolumePressed: onDecreaseVolumePressed,
	}

	cmd.songListWidget = widgets.NewList()
	cmd.songListWidget.Title = "Songs"
	cmd.songListWidget.TextStyle = ui.NewStyle(ui.ColorYellow)
	cmd.songListWidget.WrapText = false
	cmd.songListWidget.SetRect(0, 0, 100, 15)

	cmd.keyPressedWidget = widgets.NewParagraph()
	cmd.keyPressedWidget.Title = "Key"
	cmd.keyPressedWidget.SetRect(100, 0, 110, 4)

	cmd.playerStatusWidget = widgets.NewParagraph()
	cmd.playerStatusWidget.Title = "Current Song"
	cmd.playerStatusWidget.SetRect(0, 15, 100, 15+4)

	cmd.statusWidget = widgets.NewParagraph()
	cmd.statusWidget.Title = "Status"
	cmd.statusWidget.SetRect(0, 15+4, 100, 15+4+4)

	ui.Render(
		cmd.songListWidget,
		cmd.playerStatusWidget,
		cmd.statusWidget,
		cmd.keyPressedWidget,
	)

	return cmd, nil
}

func (cmd *CmdUserInterface) SetSongList(driveFiles []entities.DriveFile) {
	formattedFiles := make([]string, len(driveFiles))
	for i, file := range driveFiles {
		formattedFiles[i] = file.Name
	}

	cmd.songListWidget.Rows = formattedFiles
	cmd.driveFiles = driveFiles
	ui.Render(cmd.songListWidget)
}

func (cmd *CmdUserInterface) GetSelectedSong() entities.DriveFile {
	selectedIndex := cmd.songListWidget.SelectedRow
	return cmd.driveFiles[selectedIndex]
}

func (cmd *CmdUserInterface) ChangeStatus(status string) {
	cmd.statusWidget.Text = status
	ui.Render(cmd.statusWidget)
}

func (cmd *CmdUserInterface) ChangePlayerStatus(status string) {
	cmd.playerStatusWidget.Text = status
	ui.Render(cmd.playerStatusWidget)
}

func (cmd *CmdUserInterface) ClearPlayerStatus() {
	cmd.playerStatusWidget.Text = ""
	ui.Render(cmd.playerStatusWidget)
}

func (cmd *CmdUserInterface) Loop() {
	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		cmd.keyPressedWidget.Text = e.ID
		switch e.ID {
		case "q", "<C-c>":
			cmd.onClose(cmd)
			return
		case "j", "<Down>":
			cmd.songListWidget.ScrollDown()
		case "k", "<Up>":
			cmd.songListWidget.ScrollUp()
		case "<C-d>":
			cmd.songListWidget.ScrollHalfPageDown()
		case "<C-u>":
			cmd.songListWidget.ScrollHalfPageUp()
		case "<C-f>":
			cmd.songListWidget.ScrollPageDown()
		case "<C-b>":
			cmd.songListWidget.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				cmd.songListWidget.ScrollTop()
			}
		case "<Home>":
			cmd.songListWidget.ScrollTop()
		case "G", "<End>":
			cmd.songListWidget.ScrollBottom()
		case "<Enter>":
			cmd.onEnterPressed(cmd)
		case "p":
			cmd.onPausePressed(cmd)
		case "r":
			cmd.onResumePressed(cmd)
		case "+":
			cmd.onIncreaseVolumePressed(cmd)
		case "-":
			cmd.onDecreaseVolumePressed(cmd)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		cmd.RenderAll()
	}
}

func (cmd *CmdUserInterface) RenderAll() {
	ui.Render(
		cmd.songListWidget,
		cmd.playerStatusWidget,
		cmd.statusWidget,
		cmd.keyPressedWidget,
	)
}

func (cmd *CmdUserInterface) Close() {
	ui.Close()
}
