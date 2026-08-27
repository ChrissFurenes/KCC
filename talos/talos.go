package talos

func ImportConfig(path string) error {
	return nil
}

func Validate(path string) error {
	return nil
}

func GetServerVersion() string {
	return ""
}

func GetClientVersion() string {
	return ""
}

func ApplyConfig(path string) error {
	return nil
}

func ChangeClientTalosVersion(version string) error { // go big or go home XD
	return nil
}

func DownloadClientTalosVersion(version string) error {
	return nil
}

func FindAllClientTalosVersion() string {
	return ""
}

func FindAllTalosVersionInAllCluster() string {
	return ""
}

type TalosConfigInformation struct {
	Context  string                        `yaml:"context"`
	Contexts map[string]TalosContextConfig `yaml:"contexts"`
}
type TalosContextConfig struct {
	Target    string   `yaml:"target,omitempty"`
	Endpoints []string `yaml:"endpoints"`
	Nodes     []string `yaml:"nodes,omitempty"`
}
