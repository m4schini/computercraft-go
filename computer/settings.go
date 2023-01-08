package computer

import "context"

type SettingsOption struct {
	description string
	defaultVal  string
	typeVal     string
}

type Settings interface {
	Define(ctx context.Context, name string, option ...SettingsOption) error
	Undefine(ctx context.Context, name string) error
	Set(ctx context.Context, name, value string) error
	Unset(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (string, error)
	Clear(ctx context.Context) error
	Names(ctx context.Context) ([]string, error)
	Load(ctx context.Context, path string) (bool, error)
	Save(ctx context.Context, path string) (bool, error)
}
