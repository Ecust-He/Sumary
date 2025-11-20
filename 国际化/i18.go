package handlers

import (
	"aio.zte.com.cn/aio-scheduler/config"
	"aio.zte.com.cn/aio-scheduler/pkg/i18n"
)

var I18nIns i18n.Ii18nOpr

func NewI18nOpr() i18n.Ii18nOpr {
	bundle := i18n.NewBundle()
	bundle.AddDict("zh-cn", config.GetCnDict())
	bundle.AddDict("en-us", config.GetUsDict())

	I18nIns = i18n.NewLocalizer(bundle)
	return I18nIns
}
