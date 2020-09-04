package registry

type FileDriverConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func (c *FileDriverConfig) GetName() string {
	if "" == c.Name {
		return "file"
	}
	return c.Name
}

func (c *FileDriverConfig) GetPath() string {
	if "" == c.Path {
		return "config/registry"
	}
	return c.Path
}

type FileConfigRegister interface {
	GetName() string
	GetPath() string
}
