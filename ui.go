package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"github.com/0x17de/uoproxy/packets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	MainWindow          fyne.Window
	MegaClilocTable     *widget.Table
	LastMegaCliloc      *packets.MegaClilocPacket
	LastMegaClilocMutex sync.Mutex
)

func setLastMegaCliloc(packet *packets.MegaClilocPacket) {
	LastMegaClilocMutex.Lock()
	LastMegaCliloc = packet
	LastMegaClilocMutex.Unlock()
	MegaClilocTable.Refresh()
}

func initMegaClilocUi() fyne.CanvasObject {
	table := widget.NewTable(
		func() (int, int) {
			LastMegaClilocMutex.Lock()
			defer LastMegaClilocMutex.Unlock()
			if LastMegaCliloc == nil || LastMegaCliloc.Server == nil {
				return 0, 3
			}
			return len(LastMegaCliloc.Server.Clilocs), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("placeholder")
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			LastMegaClilocMutex.Lock()
			defer LastMegaClilocMutex.Unlock()
			if tci.Row >= len(LastMegaCliloc.Server.Clilocs) {
				return
			}
			data := LastMegaCliloc.Server.Clilocs[tci.Row]
			if tci.Col == 0 {
				co.(*widget.Label).SetText("0x" + strconv.FormatUint(uint64(data.Id), 16))
			} else if tci.Col == 1 {
				co.(*widget.Label).SetText(packets.ClilocId(data.Id).String())
			} else {
				co.(*widget.Label).SetText(string(data.Text))
			}
		},
	)
	table.SetColumnWidth(1, 200)
	MegaClilocTable = table
	return table
}

func initSettingsUi() fyne.CanvasObject {
	lobbySelect := widget.NewSelect(
		[]string{
			"shard.uoex.net:60",
			"localhost:2593",
		},
		func(s string) {
			hostPort := strings.SplitN(s, ":", 2)
			LobbyHost = hostPort[0]
			LobbyPort, _ = strconv.Atoi(hostPort[1])
			log.Printf("Selected host %s:%d", LobbyHost, LobbyPort)
		},
	)
	lobbySelect.SetSelectedIndex(0)
	form := widget.NewForm(
		widget.NewFormItem("Host/Port", lobbySelect),
	)
	return form
}

func initUi() {
	a := app.New()
	MainWindow = a.NewWindow("UOProxy2")
	MainWindow.SetMaster()

	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", initSettingsUi()),
		container.NewTabItem("MegaCliloc", initMegaClilocUi()),
	)
	MainWindow.SetContent(tabs)
	MainWindow.Resize(fyne.NewSize(800, 600))
}

func runUi() {
	MainWindow.ShowAndRun()
}
