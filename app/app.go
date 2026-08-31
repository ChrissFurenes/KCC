package app

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/chrissfurenes/kcc/cmd"
	"github.com/chrissfurenes/kcc/kube"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type App struct {
	version       string
	Items         []Item
	KubePath      string
	TalosPath     string
	CurrentFolder string
	UI            *tview.Application
	ConfigList    *tview.List
	InfoData      *tview.TextView
	CommandList   *tview.TextView
	Grid          *tview.Grid
	Kube          *kube.Kube
}

type Item struct {
	Name        string
	Path        string
	FileName    string
	File        []byte
	IsDir       bool
	IsConfig    bool
	IsActive    bool
	IsTalos     bool
	IsBack      bool
	Config      kube.KubeConfigInformation
	ClusterData ClusterData
}

type ClusterData struct {
	User         string
	Address      string
	Port         string
	Reachable    bool
	Nodes        int
	Pods         int
	TalosVersion string
	Status       string
	Test         string
}

func NewApp(version string) *App {
	a := &App{
		KubePath:      cmd.KubePath(),
		TalosPath:     cmd.TalosPath(),
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

func (a *App) ConfigDir() string {
	return filepath.Join(a.KubePath, "configs", a.CurrentFolder)
}

func talosConfigPath(talosPath, fileName string) string {
	return filepath.Join(talosPath, "configs", fileName)
}

func (i *Item) DisplayName() string {
	if i.IsBack {
		return " << Back to folder: " + i.Name
	}
	if i.IsDir {
		return "📁 " + i.Name
	}
	if !i.IsConfig {
		return "⚠ " + i.Name
	}
	name := "☸  " + i.Name
	if i.IsActive {
		name += " - [green]ACTIVE[::-]"
	}
	return name
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
	err := a.LoadConfigs()

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

func statusColorIcon(ok bool) (color, icon string) {
	if ok {
		return "[green]", " 🟢"
	}
	return "[red]", "🔴"
}

func (i *Item) InfoText() string {

	if i.IsBack {
		return ""
	}

	if i.IsDir {
		entries, err := os.ReadDir(i.Path)
		if err != nil {
			return "Folder path: " + i.Path + "\n\n[red]" + err.Error() + "[::-]"
		}
		return fmt.Sprintf("Folder path: %s\nItems: %d", i.Path, len(entries))
	}

	if !i.IsConfig {
		information := "Name:.. " + i.Name + "\nPath:.. " + filepath.Base(i.Path)
		if i.ClusterData.Status != "" {
			information += "\n\nStatus: " + i.ClusterData.Status
		}
		return information
	}

	data := i.ClusterData
	color, statusIcon := statusColorIcon(data.Reachable)
	talosColor, talosIcon := statusColorIcon(i.IsTalos)
	information := "Name:.. " + i.Name +
		"\n\nUser:.. " + data.User +
		"\nIP:.... " + data.Address +
		"\nPort:.. " + data.Port +
		"\nPing:.. " + color + strings.ToUpper(strconv.FormatBool(data.Reachable)) + "[::-] [white]" + statusIcon +
		"\nTalos:. " + talosColor + strings.ToUpper(strconv.FormatBool(i.IsTalos)) + "[::-] [white]" + talosIcon +
		"\nKube Path:.. " + filepath.Base(i.Path) + "\n"

	if data.Reachable {
		information += "\nKubernetes:\nNodes:. " + strconv.Itoa(data.Nodes) + "\nPods:.. " + strconv.Itoa(data.Pods)
		if data.TalosVersion != "" {
			information += "\nOS:.... " + data.TalosVersion
		}
	}
	if data.Status != "" {
		information += "\n\nStatus: " + data.Status
	}
	if data.Test != "" {
		information += "\n\n\nTests:. " + data.Test
	}
	return information
}

func (a *App) LoadConfigs() error {
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
		a.Items = append(a.Items, Item{Name: name, IsBack: true})
	}
	for _, entry := range entries {
		fullPath := filepath.Join(a.ConfigDir(), entry.Name())
		var item Item
		err := item.Load(fullPath)
		if err != nil {
			item.Name = entry.Name()
			item.Path = fullPath
			item.FileName = entry.Name()
			item.IsDir = entry.IsDir()
			item.ClusterData.Status = "[red]" + err.Error() + "[::-]"
			a.Items = append(a.Items, item)
			continue
		}
		if item.IsConfig {
			item.IsActive = item.IsCurrent(a.KubePath)
			if _, err := os.Stat(talosConfigPath(a.TalosPath, item.FileName)); err == nil {
				item.IsTalos = true
			}
		}
		a.Items = append(a.Items, item)
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
		err := a.LoadConfigs()
		if err != nil {
			a.CurrentFolder = previousFolder
			a.InfoData.SetText("[red]" + err.Error() + "[::-]")
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
		err := item.Apply(a.KubePath, a.TalosPath)

		if err != nil {
			item.ClusterData.Status = "[red]Apply failed: " + err.Error() + "[::-]"
			a.InfoData.SetText(item.InfoText())
			return
		}
		a.UI.Stop()
	}
}
func (a *App) RefreshConfigList() {
	oldPosition := a.ConfigList.GetCurrentItem()
	a.ConfigList.Clear()
	for index := range a.Items {
		itemIndex := index
		a.ConfigList.AddItem(a.Items[index].DisplayName(), "", 0, func() { a.OpenItem(itemIndex) })
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
	items := make([]Item, len(a.Items))
	copy(items, a.Items)
	for index := range items {

		if !items[index].IsConfig {
			continue
		}
		_ = items[index].RefreshClusterInfo(a.TalosPath)
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
			a.InfoData.SetText(a.Items[current].InfoText())
		}
	})
}

func (i *Item) Apply(kubePath, talosPath string) error {
	if !i.IsConfig {
		return fmt.Errorf("%s is not a config", i.Name)
	}

	dest := filepath.Join(kubePath, "config")
	backup := dest + ".kcc-backup"
	_ = os.Remove(backup)

	if _, err := os.Stat(dest); err == nil {
		if err := os.Rename(dest, backup); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.Link(i.Path, dest); err != nil {
		if _, backupErr := os.Stat(backup); backupErr == nil {
			_ = os.Rename(backup, dest)
		}
		return err
	}
	_ = os.Remove(backup)
	i.IsActive = true

	if i.IsTalos {
		talosSource := talosConfigPath(talosPath, i.FileName)
		talosDest := filepath.Join(talosPath, "config")
		_ = os.Remove(talosDest)
		_ = os.Link(talosSource, talosDest)
	}
	return nil
}

func (i *Item) RefreshClusterInfo(talosPath string) error {
	if !i.IsConfig || i.IsDir || i.IsBack {
		return nil
	}

	i.ClusterData.Status = "[yellow]Getting info from cluster....[::-]"

	if !i.ClusterData.Ping() {
		i.ClusterData.Status = "[red]Offline[::-]"
		return nil
	}
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", i.Path)
	if err != nil {
		i.ClusterData.Status = "[red]Config error: " + err.Error() + "[::-]"
		return err
	}
	clientSet, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		i.ClusterData.Status = "[red]Client error: " + err.Error() + "[::-]"
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	nodes, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		i.ClusterData.Status = "[red]Node API error: " + err.Error() + "[::-]"
		return err
	}
	i.ClusterData.Nodes = len(nodes.Items)
	pods, err := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		i.ClusterData.Status = "[red]Pod API error: " + err.Error() + "[::-]"
		return err
	}
	i.ClusterData.Pods = len(pods.Items)

	if i.IsTalos {
		i.ClusterData.TalosVersion = talosVersion(talosPath, i.FileName, i.ClusterData.Address)
	}

	i.ClusterData.Status = ""
	return nil
}

func talosVersion(talosPath, fileName, address string) string {
	talosConfig := talosConfigPath(talosPath, fileName)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "talosctl", "version",
		"--talosconfig", talosConfig, "--nodes", address, "--short").Output()

	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func ensureConfigsDir(base string, required bool) error {
	if _, err := os.Stat(base); err != nil {
		if required {
			return err
		}
		return nil
	}
	configsPath := filepath.Join(base, "configs")
	if _, err := os.Stat(configsPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(configsPath, 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join(base, "config"))
	if err != nil {
		if required {
			return err
		}
		return nil
	}
	return os.WriteFile(filepath.Join(configsPath, "config"), data, 0644)
}

func (a *App) EnsureConfigPath() error {
	if err := ensureConfigsDir(a.KubePath, true); err != nil {
		return err
	}
	return ensureConfigsDir(a.TalosPath, false)
}

func (a *App) Import(ftype string, from string, to string) error {
	source := from

	if !filepath.IsAbs(source) {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		source = filepath.Join(dir, source)
	}

	var testItem Item
	if slices.Contains([]string{"kube", "kubeconfig", "k8s"}, strings.ToLower(ftype)) {
		err := testItem.Load(source)
		if err != nil {
			return fmt.Errorf("invalid kubeconfig: %w", err)
		}

		if !testItem.IsConfig {
			return fmt.Errorf("%s is not a config file", source)
		}

		destination := filepath.Join(a.KubePath, "configs", to)
		err = os.MkdirAll(filepath.Dir(destination), 0755)
		if err != nil {
			return err
		}
		file, err := os.ReadFile(source)
		if err != nil {
			return err
		}
		return os.WriteFile(destination, file, 0644)
	}
	return nil
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
			return true, fmt.Errorf("usage: kcc import <kube/talos> <from> <to>")
		}
		err := a.Import(args[1], args[2], args[3])
		if err != nil {
			return true, err
		}
		fmt.Println("Config imported")
		return true, nil

	default:
		return true, fmt.Errorf("unknown command: %s", args[0])
	}
}

func (a *App) Run() error {
	if err := a.EnsureConfigPath(); err != nil {
		return err
	}

	if err := a.LoadConfigs(); err != nil {
		return err
	}

	a.ConfigList.SetBorder(true).SetTitle("Configuration")
	a.InfoData.SetDynamicColors(true).SetBorder(true).SetTitle("Info").SetTitleAlign(tview.AlignCenter)
	a.ConfigList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index < 0 || index >= len(a.Items) {
			return
		}
		a.InfoData.SetText(a.Items[index].InfoText())
	})

	a.RefreshConfigList()
	go a.RefreshClusterInfo()

	a.UI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF5:
			for index := range a.Items {
				if a.Items[index].IsConfig {
					a.Items[index].ClusterData.Status = "[yellow]Getting info from cluster....[::-]"
				}
			}
			current := a.ConfigList.GetCurrentItem()
			if current >= 0 && current < len(a.Items) {
				a.InfoData.SetText(a.Items[current].InfoText())
			}
			go a.RefreshClusterInfo()
		case tcell.KeyBackspace, tcell.KeyEsc:
			a.GoBack()
		}
		return event
	})
	return a.UI.SetRoot(a.Grid, true).Run()
}

func (c *ClusterData) SetServer(server string) error {
	u, err := url.Parse(server)
	if err != nil {
		return err
	}
	if u.Hostname() == "" {
		return fmt.Errorf("invalid server address: %s", server)
	}

	c.Address = u.Hostname()
	c.Port = u.Port()

	if c.Port == "" {
		switch u.Scheme {
		case "https":
			c.Port = "443"
		case "http":
			c.Port = "80"
		default:
			return fmt.Errorf(
				"server has no port: %s",
				server,
			)
		}
	}
	return nil
}

func (c *ClusterData) Ping() bool {
	address := net.JoinHostPort(c.Address, c.Port)

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)

	if err != nil {
		c.Reachable = false
		return false
	}
	defer conn.Close()
	c.Reachable = true
	return true
}

func (i *Item) Load(path string) error {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return err
	}

	i.Path = filepath.Clean(path)
	i.FileName = fileInfo.Name()
	i.Name = fileInfo.Name()

	i.IsDir = fileInfo.IsDir()
	i.IsConfig = false
	if i.IsDir {
		return nil
	}
	file, err := os.ReadFile(i.Path)
	if err != nil {
		return err
	}
	i.File = file
	var config kube.KubeConfigInformation
	if err := yaml.Unmarshal(file, &config); err != nil {
		return err
	}

	if len(config.Clusters) == 0 {
		return fmt.Errorf("config contains no clusters")
	}

	if len(config.Contexts) == 0 {
		return fmt.Errorf("config contains no contexts")
	}

	i.Config = config
	i.Name = config.Clusters[0].Name
	i.IsConfig = true

	i.ClusterData.User = config.Contexts[0].Context.User
	i.ClusterData.Status = "[yellow]Getting info from cluster....[::-]"
	err = i.ClusterData.SetServer(config.Clusters[0].Cluster.Server)
	if err != nil {
		return err
	}
	return nil
}

func (i *Item) GetFile() ([]byte, error) {
	if !i.IsConfig || i.IsDir {
		return nil, fmt.Errorf("%s is not a config file", i.Path)
	}
	file, err := os.ReadFile(i.Path)
	if err != nil {
		return nil, err
	}
	i.File = file
	return file, nil
}

func (i *Item) IsCurrent(kubePath string) bool {
	if !i.IsConfig {
		return false
	}
	currentPath := filepath.Join(kubePath, "config")

	sourceInfo, sourceErr := os.Stat(i.Path)
	currentInfo, currentErr := os.Stat(currentPath)
	if sourceErr == nil && currentErr == nil && os.SameFile(sourceInfo, currentInfo) {
		return true
	}

	current, err := os.ReadFile(currentPath)
	if err != nil {
		return false
	}
	file, err := os.ReadFile(i.Path)
	if err != nil {
		return false
	}
	return bytes.Equal(file, current)
}
