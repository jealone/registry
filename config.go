package registry

type Config struct {
	Driver YamlNode `yaml:"driver"`
	Type   string   `yaml:"type"`
}

func (c *Config) GetDriver() *YamlNode {
	return &c.Driver
}

func (c *Config) GetType() string {
	return c.Type
}

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
