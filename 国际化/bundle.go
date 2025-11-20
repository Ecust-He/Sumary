package i18n

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"aio.zte.com.cn/aio-scheduler/pkg/logger"
)

type TranslateDict struct {
	Language string         `yaml:"language"`
	Contents []*DictContent `yaml:"contents"`
}

type DictContent struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type dict map[string]string
type Dicts map[string]dict

type Bundle struct {
	Dicts Dicts
}

func NewBundle() *Bundle {
	var bundle = &Bundle{
		Dicts: map[string]dict{},
	}

	return bundle
}

func (i *Bundle) AddDict(key, valueStr string) {
	logger.Logger().Infof("Start reading %s dict", key)
	if err := i.parseSingleDict(valueStr); err != nil {
		logger.Logger().Errorf("Parse single file failed! err : %v", err)
	}
	logger.Logger().Info("Loading finished...")
}

func (i *Bundle) parseSingleDict(valueStr string) error {
	translateDict, err := parseTranslateDict([]byte(valueStr))
	if err != nil {
		return err
	}
	if i.Dicts[translateDict.Language] == nil {
		i.Dicts[translateDict.Language] = dict{}
	}
	for _, content := range translateDict.Contents {
		i.Dicts[translateDict.Language][content.Key] = content.Value
	}

	return nil
}

func parseTranslateDict(buf []byte) (*TranslateDict, error) {
	var translateDict TranslateDict
	if err := yaml.Unmarshal(buf, &translateDict); err != nil {
		return nil, errors.Wrapf(err, "unmarshal values failed. buf str : %s", string(buf))
	}

	return &translateDict, nil
}
