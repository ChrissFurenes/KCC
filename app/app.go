package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrissfurenes/kcc/cmd"
	"github.com/chrissfurenes/kcc/kcc"
	"github.com/chrissfurenes/kcc/kube"
	"github.com/chrissfurenes/kcc/talos"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	version       string
	Items         []kcc.Item
	KubePath      string
	TalosPath     string
	CurrentFolder string
	UI            *tview.Application
	ConfigList    *tview.List
	InfoData      *tview.TextView
	CommandList   *tview.TextView
	Grid          *tview.Grid
}

func NewApp(version string) *App { // OK
	a := &App{
		KubePath: cmd.KubePath(),
		//TalosPath:     cmd.TalosPath(),
		CurrentFolder: "",
		UI:            tview.NewApplication(),
		ConfigList:    tview.NewList().ShowSecondaryText(false),
		InfoData:      tview.NewTextView(),
		CommandList:   tview.NewTextView().SetText("[F5] Refresh").SetTextAlign(tview.AlignCenter),
		version:       version,
	}
	a.Grid = tview.NewGrid().
		SetRows(-1, 25).
		SetColumns(-1, -1).
		SetBorders(false).
		AddItem(a.ConfigList, 0, 0, 5, 1, 0, 0, true).
		AddItem(a.InfoData, 0, 1, 5, 1, 0, 0, false).
		AddItem(a.CommandList, 5, 0, 1, 2, 1, 0, false)

	return a
}

func (a *App) ConfigDir() string { // needs change
	return filepath.Join(a.KubePath, "configs", a.CurrentFolder)
}

func (a *App) GoBack() {
	if a.CurrentFolder == "" {
		return
	}

	parent := filepath.Dir(a.CurrentFolder)
	if parent == "." {
		parent = ""
	}

	previousFolder := a.CurrentFolder
	a.CurrentFolder = parent
	err := a.LoadEntities()

	if err != nil {
		a.CurrentFolder = previousFolder
		a.InfoData.SetText("[red]" + err.Error() + "[::-]")
		return
	}

	a.RefreshConfigList()
	if a.ConfigList.GetItemCount() > 0 {
		a.ConfigList.SetCurrentItem(0)
	}
	go a.RefreshClusterInfo()
}

func (a *App) LoadEntities() error { // need some changes or remove
	a.Items = nil
	entries, err := os.ReadDir(a.ConfigDir())
	if err != nil {
		return err
	}
	if a.CurrentFolder != "" {
		parent := filepath.Dir(a.CurrentFolder)
		name := "configs"
		if parent != "." && parent != "" {
			name = filepath.Base(parent)
		}
		a.Items = append(a.Items, kcc.Item{Name: name, IsBack: true})

	}
	for _, entry := range entries {
		var currentEntity kcc.Item

		if entry.IsDir() {
			currentEntity.Name = entry.Name()
			currentEntity.IsDir = true

		} else {
			k := kube.NewKube(filepath.Join(a.ConfigDir(), entry.Name()))
			currentEntity.Name = k.GetName()
			currentEntity.Path = filepath.Join(a.ConfigDir(), entry.Name())
			currentEntity.FileName = entry.Name()
			currentEntity.IsDir = false
			currentEntity.IsConfig = true
			currentEntity.ClusterData.Address = k.Address
			currentEntity.ClusterData.Port = k.Port
			currentEntity.ClusterData.Reachable = k.Reachable
			currentEntity.ClusterData.Status = k.Status
			currentEntity.ClusterData.User = k.User
		}
		a.Items = append(a.Items, currentEntity)
	}

	return nil
}

func (a *App) OpenItem(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	item := &a.Items[index]
	if item.IsBack {
		a.GoBack()
		return
	}

	if item.IsDir {
		previousFolder := a.CurrentFolder
		a.CurrentFolder = filepath.Join(a.CurrentFolder, item.Name)
		err := a.LoadEntities()
		if err != nil {
			a.CurrentFolder = previousFolder
			a.InfoData.SetText(kcc.RedText(err.Error()))
			return
		}
		a.RefreshConfigList()
		if a.ConfigList.GetItemCount() > 1 {
			a.ConfigList.SetCurrentItem(1)
		}
		go a.RefreshClusterInfo()
		return
	}
	if item.IsConfig {
		//err := item.Apply(a.KubePath, a.TalosPath)

		//if err != nil {
		//	item.ClusterData.Status = "[red]Apply failed: " + err.Error() + "[::-]"
		a.InfoData.SetText(item.GetInfoText())
		return
		//}
		a.UI.Stop()
	}
}
func (a *App) RefreshConfigList() { // beholde
	oldPosition := a.ConfigList.GetCurrentItem()
	a.ConfigList.Clear()
	for index := range a.Items {
		itemIndex := index
		a.ConfigList.AddItem(a.Items[index].GetDisplayName(), "", 0, func() { a.OpenItem(itemIndex) })
	}
	count := a.ConfigList.GetItemCount()
	if count == 0 {
		return
	}

	if oldPosition >= count {
		oldPosition = count - 1
	}

	if oldPosition < 0 {
		oldPosition = 0
	}
	a.ConfigList.SetCurrentItem(oldPosition)
}

func (a *App) RefreshClusterInfo() {
	folder := a.CurrentFolder
	items := make([]kcc.Item, len(a.Items))
	copy(items, a.Items)
	for index := range items {

		if !items[index].IsConfig {
			continue
		}
		//_ = items[index].RefreshClusterInfo(a.TalosPath)
	}

	a.UI.QueueUpdateDraw(func() {
		if folder != a.CurrentFolder {
			return
		}

		current := a.ConfigList.GetCurrentItem()
		a.Items = items
		a.RefreshConfigList()
		if current >= 0 && current < len(a.Items) {
			a.ConfigList.SetCurrentItem(current)
			a.InfoData.SetText(a.Items[current].GetInfoText())
		}
	})
}

func (a *App) HandleArgs(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	switch args[0] {

	case "version":
		fmt.Println("Version: " + a.version)
		return true, nil

	case "help", "h":
		cmd.Help()
		return true, nil

	case "import", "i":
		if len(args) != 3 {
			return true, fmt.Errorf("usage: kcc import [kube/talos] <from> <to>")
		}
		switch args[1] {
		case "kubeconfig":
			err := kube.ImportConfig(args[2], args[3])
			return true, err
		case "talos":
			err := talos.ImportConfig(args[2], args[3])
			return true, err
		}
		fmt.Println("Config imported")
		return true, nil

	default:
		return true, fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a *App) Run() error {
	//if err := a.EnsureConfigPath(); err != nil {
	//	return err
	//}
	if strings.Contains(a.version, "beta") {
		//kcc.Backup(a.version) // will be active in beta prod
	}
	if err := a.LoadEntities(); err != nil {
		return err
	}

	a.ConfigList.SetBorder(true).SetTitle("Configuration")
	a.InfoData.SetDynamicColors(true).SetBorder(true).SetTitle("Info").SetTitleAlign(tview.AlignCenter)
	a.ConfigList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index < 0 || index >= len(a.Items) {
			return
		}
		a.InfoData.SetText(a.Items[index].GetInfoText())
	})

	a.RefreshConfigList()
	go a.RefreshClusterInfo()

	a.UI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF5:
			for index := range a.Items { // need to change
				if a.Items[index].IsConfig {
					a.Items[index].ClusterData.Status = kcc.YellowText("Getting info from cluster....")
				}
			}
			current := a.ConfigList.GetCurrentItem()
			if current >= 0 && current < len(a.Items) {
				a.InfoData.SetText(a.Items[current].GetInfoText())
			}
			go a.RefreshClusterInfo()
		case tcell.KeyBackspace, tcell.KeyEsc:
			a.GoBack()
		}
		return event
	})
	return a.UI.SetRoot(a.Grid, true).Run()
}
