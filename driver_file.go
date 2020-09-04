package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jealone/goconf"
	"github.com/jealone/sli4go"
)

func NewYamlFileDriver(conf FileConfigRegister) (*YamlFileDriver, error) {

	abs, err := filepath.Abs(conf.GetPath())

	if nil != err {
		return nil, err
	}
	driver := new(YamlFileDriver)
	driver.path = abs
	driver.Name = conf.GetName()

	return driver, nil
}

type YamlFileDriver struct {
	Name    string
	path    string
	Entries sync.Map
}

func (d *YamlFileDriver) Boot() {

	var configs []*goconf.YamlConfig

	err := filepath.Walk(d.path, func(path string, info os.FileInfo, err error) error {

		if nil != err {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if ext := filepath.Ext(path); ".yml" != ext && ".yaml" != ext {
			return nil
		}

		confs, err := goconf.NewConfigFile(path)

		if nil != err {
			return err
		}

		configs = append(configs, confs...)

		return nil
	})

	if nil != err {
		sli4go.Fatalf("registry boot error : %s", err)
	}

	for _, config := range configs {
		key := fmt.Sprintf("%s:%s", config.GetKind(), config.GetKey())
		if _, ok := d.Entries.Load(key); ok {
			sli4go.Warnf("skip conflict config key:%s", key)
			continue
		}
		d.Entries.Store(key, config.RawByte)

	}

}

func (d *YamlFileDriver) GetName() string {
	if "" == d.Name {
		return "file"
	}
	return d.Name
}

func (d *YamlFileDriver) GetEntry(key string) []byte {

	e, ok := d.Entries.Load(key)

	if !ok {
		return nil
	}

	if entry, ok := e.([]byte); ok {
		return entry
	} else {
		sli4go.Error("registry store wrong type except raw []byte")
		return nil
	}

}
