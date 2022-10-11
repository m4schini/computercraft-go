package computer

type SettingsOption struct {
	description string
	defaultVal  string
	typeVal     string
}

type Settings interface {
	Define(name string, option ...SettingsOption) error
	Undefine(name string) error
	Set(name, value string) error
	Unset(name string) error
	Get(name string) (string, error)
	Clear() error
	Names() ([]string, error)
	Load(path string) (bool, error)
	Save(path string) (bool, error)
}
