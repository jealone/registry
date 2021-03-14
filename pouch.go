package registry

import (
	"errors"

	"github.com/jealone/goconf"
	"gopkg.in/yaml.v3"
)

type (
	YamlNode   = yaml.Node
	YamlConfig = goconf.YamlConfig
)

var (
	errorEmptyRawBytes = errors.New("buf bytes is nil")
)

func ParseYamlConfig(buf []byte) (*YamlConfig, error) {

	if 0 == len(buf) {
		return nil, errorEmptyRawBytes
	}

	conf := &YamlConfig{}
	err := goconf.Unmarshal(append(buf, '\n'), conf)
	if nil != err {
		return nil, err
	}
	return conf, nil
}
