package ui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/jibaru/gobeats/m/internal/entities"
	"math/rand"
)

type CmdUserInterface struct {
	driveFiles              entities.DriveFileList
	keyPressedWidget        *widgets.Paragraph
	songListWidget          *widgets.List
	playerStatusWidget      *widgets.Paragraph
	statusWidget            *widgets.Paragraph
	helpWidget              *widgets.Paragraph
	onClose                 func(cmd *CmdUserInterface)
	onEnterPressed          func(cmd *CmdUserInterface)
	onPausePressed          func(cmd *CmdUserInterface)
	onResumePressed         func(cmd *CmdUserInterface)
	onIncreaseVolumePressed func(cmd *CmdUserInterface)
	onDecreaseVolumePressed func(cmd *CmdUserInterface)
	onIncreaseSongsIndexSet func(cmd *CmdUserInterface)
}

func NewCmdUserInterface(
	onClose func(cmd *CmdUserInterface),
	onEnterPressed func(cmd *CmdUserInterface),
	onPausePressed func(cmd *CmdUserInterface),
	onResumePressed func(cmd *CmdUserInterface),
	onIncreaseVolumePressed func(cmd *CmdUserInterface),
	onDecreaseVolumePressed func(cmd *CmdUserInterface),
	onIncreaseSongsIndexSet func(cmd *CmdUserInterface),
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
		onIncreaseSongsIndexSet: onIncreaseSongsIndexSet,
	}

	cmd.songListWidget = widgets.NewList()
	cmd.songListWidget.Title = "Songs"
	cmd.songListWidget.TextStyle = ui.NewStyle(ui.ColorYellow)
	cmd.songListWidget.WrapText = false
	cmd.songListWidget.SetRect(0, 0, 100, 15)

	cmd.keyPressedWidget = widgets.NewParagraph()
	cmd.keyPressedWidget.Title = "Key"
	cmd.keyPressedWidget.SetRect(100, 0, 115, 4)

	cmd.helpWidget = widgets.NewParagraph()
	cmd.helpWidget.Title = "Help"
	cmd.helpWidget.SetRect(100, 4, 115, 23)
	cmd.helpWidget.Text = "s: shuffle\n" +
		"p: pause\n" +
		"r: resume\n" +
		"↑: scroll up\n" +
		"↓: scroll down\n" +
		"enter: select song\n" +
		"+: inc. volume\n" +
		"-: dec. volume\n" +
		"q: exit"

	cmd.playerStatusWidget = widgets.NewParagraph()
	cmd.playerStatusWidget.Title = "Player Status"
	cmd.playerStatusWidget.SetRect(0, 15, 100, 15+4)

	cmd.statusWidget = widgets.NewParagraph()
	cmd.statusWidget.Title = "Status"
	cmd.statusWidget.SetRect(0, 15+4, 100, 15+4+4)

	ui.Render(
		cmd.songListWidget,
		cmd.playerStatusWidget,
		cmd.statusWidget,
		cmd.keyPressedWidget,
		cmd.helpWidget,
	)

	return cmd, nil
}

func (cmd *CmdUserInterface) SetSongList(songs entities.DriveFileList) {
	cmd.driveFiles = songs
	cmd.reloadSongsWidget()
}

func (cmd *CmdUserInterface) reloadSongsWidget() {
	formattedFiles := make([]string, len(cmd.driveFiles))
	for i, file := range cmd.driveFiles {
		formattedFiles[i] = file.Name
	}

	cmd.songListWidget.Rows = formattedFiles
	ui.Render(cmd.songListWidget)
}

func (cmd *CmdUserInterface) GetSelectedSong() entities.DriveFile {
	selectedIndex := cmd.songListWidget.SelectedRow
	return cmd.driveFiles[selectedIndex]
}

func (cmd *CmdUserInterface) IncreaseSelectedSongIndex() {
	if cmd.songListWidget.SelectedRow+1 < len(cmd.driveFiles) {
		cmd.songListWidget.SelectedRow = cmd.songListWidget.SelectedRow + 1
		cmd.reloadSongsWidget()
		cmd.onIncreaseSongsIndexSet(cmd)
	}
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
		case "s":
			rand.Shuffle(
				len(cmd.driveFiles),
				func(i, j int) {
					cmd.driveFiles[i], cmd.driveFiles[j] = cmd.driveFiles[j], cmd.driveFiles[i]
				},
			)
			cmd.reloadSongsWidget()
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
