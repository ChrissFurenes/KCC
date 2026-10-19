package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrissfurenes/kcc/cli"
	"github.com/chrissfurenes/kcc/cmd"
	"github.com/chrissfurenes/kcc/kcc"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
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
		KubePath:      cmd.KubePath(),
		TalosPath:     cmd.TalosPath(),
		CurrentFolder: "",
		UI:            tview.NewApplication(),
		ConfigList:    tview.NewList().ShowSecondaryText(false),
		InfoData:      tview.NewTextView(),
		CommandList:   tview.NewTextView().SetText("[ENTER] Apply KCC file [F5] Refresh  [^K] Apply kubeconfig  [^T] Apply talosconfig").SetTextAlign(tview.AlignCenter),
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

func (a *App) LoadEntities() error { // TODO: need some changes or remove
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

		} else { // TODO Need to validate
			kc := kcc.NewKubeConfig(filepath.Join(a.ConfigDir(), entry.Name()))

			//k := kube.NewKube(filepath.Join(a.ConfigDir(), entry.Name()))
			currentEntity.Kube = *kc
			currentEntity.Name = kc.GetName()
			currentEntity.Path = filepath.Join(a.ConfigDir(), entry.Name())
			currentEntity.FileName = entry.Name()
			currentEntity.IsDir = false
			currentEntity.IsConfig = true
			currentEntity.Kube.Address = kc.Address
			currentEntity.Kube.Port = kc.Port
			currentEntity.Kube.Reachable = kc.Reachable
			currentEntity.Kube.Status = kc.Status
			currentEntity.Kube.User = kc.User
		}
		a.Items = append(a.Items, currentEntity)
	}

	return nil
}

func (a *App) OpenItem(index int) {
	if index < 0 || index >= len(a.Items) {
		return
	}
	item := a.Items[index]
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
			a.InfoData.SetText(cmd.RedText(err.Error()))
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
		//return
		//}
		a.UI.Stop()
	}
}
func addLine(view *tview.TextView, left, right string) string { // TODO: NEED fixing
	_, _, width, _ := view.GetInnerRect()

	leftWidth := runewidth.StringWidth(left)
	rightWidth := runewidth.StringWidth(right)

	spaces := width - leftWidth - rightWidth

	if spaces < 1 {
		spaces = 1
	}

	return left + strings.Repeat(" ", spaces) + right
}
func (a *App) RefreshConfigList() { // beholde
	oldPosition := a.ConfigList.GetCurrentItem()
	a.ConfigList.Clear()
	for index := range a.Items {
		itemIndex := index
		displayName := ""

		displayName = a.Items[index].GetDisplayName()

		a.ConfigList.AddItem(displayName, "", 0, func() { a.OpenItem(itemIndex) })
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
	for index := range a.Items {
		item := items[index]
		item.Kube.Ping()
	}
	a.UI.QueueUpdateDraw(func() {
		if folder != a.CurrentFolder {
			return
		}
	})

	//for index := range items { // TODO: NEED Fixing.
	//
	//	if !items[index].IsConfig {
	//		continue
	//	}
	//
	//	a.RefreshConfigList()
	//	if items[index].Kube.Ping() {
	//		//a.Items[index].StatusText = "🟢"
	//	} else {
	//		//a.Items[index].StatusText = "🔴"
	//	}
	//	a.RefreshConfigList()
	//	//_ = items[index].RefreshClusterInfo(a.TalosPath)
	//}

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
	case "backup", "-b", "--backup":
		err := cmd.CreateBackup()
		return true, err

	case "import", "i", "-i", "--import":

		if len(args) >= 5 {
			return true, fmt.Errorf("usage: kcc import <kubeconfig> <talosconfig> <to> (to is in ~/.kcc/configs/)")
		}

		if len(args) < 3 {
			return true, fmt.Errorf("usage: kcc import <kubeconfig> <talosconfig> <to> (to is in ~/.kcc/configs/)")
		}
		importexport := cli.Importexport{
			KubeConfigFilePath:  args[1],
			TalosConfigFilePath: args[2],
			OutputFilePath:      args[3],
		}
		importexport.Import()
		fmt.Println("Config imported")
		return true, nil

	case "export", "e":
		if len(args) >= 4 {
			return true, fmt.Errorf("usage: kcc export <kcc config> <dest folder>")
		}
		if len(args) < 2 {
			return true, fmt.Errorf("usage: kcc export <kcc config> <dest folder>")
		}
		exportconf := cli.Importexport{
			InputFilePath:  args[1],
			OutputFilePath: args[2],
		}
		exportconf.Export()

	}

	return true, nil

}

func (a *App) Run() error {
	//if err := a.EnsureConfigPath(); err != nil {
	//	return err
	//}
	if strings.Contains(a.version, "beta") {
		//kcc.Backup(a.version) // TODO: will be active in beta prod
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
			for index := range a.Items { // TODO: need to change
				if a.Items[index].IsConfig {
					a.Items[index].Kube.Status = cmd.YellowText("Getting info from cluster....")
					a.Items[index].StatusText = "🟡"
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
