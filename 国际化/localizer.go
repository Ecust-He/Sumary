package i18n

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	"aio.zte.com.cn/aio-scheduler/pkg/logger"
)

const VarTemplate = "${var}"

// Localizer provides Translate methods that return localized messages.
type Localizer struct {
	Dicts Dicts
	lock  sync.Mutex
}

func NewLocalizer(bundle *Bundle) Ii18nOpr {
	return &Localizer{
		Dicts: bundle.Dicts,
	}
}

func (l *Localizer) Translate(language, key string, vars ...string) string {
	dict, ok := l.Dicts[language]
	if !ok {
		logger.Logger().Warnf("Text yaml in %s format is not loaded.", language)
		return ""
	}
	value, ok := dict[key]
	if !ok {
		logger.Logger().Warnf("The key is not found in the %s Dictionary, key is %s", language, key)
		return ""
	}
	for _, v := range vars {
		// 每次只替换1个，在循环的加持下，可以实现从左到右依次替换的目的
		value = strings.Replace(value, VarTemplate, v, 1)
	}

	return value
}

func (l *Localizer) PutValue(language, key, value string) error {
	dict, ok := l.Dicts[language]
	if !ok {
		return errors.Errorf("text yaml in %s format is not loaded.", language)
	}
	_, ok = dict[key]
	if ok {
		return errors.New("The corresponding key already exists!")
	}
	// add
	l.lock.Lock()
	defer l.lock.Unlock()
	dict[key] = value

	return nil
}

func (l *Localizer) RemoveValue(language, key string) error {
	dict, ok := l.Dicts[language]
	if !ok {
		return errors.Errorf("text yaml in %s format is not loaded.", language)
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(dict, key)

	return nil
}
