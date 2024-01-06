package handler

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func getLanguage(ctx *gin.Context) language.Tag {
	acceptLanguage, _, err := language.ParseAcceptLanguage(ctx.GetHeader("Accept-Language"))
	if err != nil {
		return language.Persian
	}
	for _, al := range acceptLanguage {
		switch al {
		case language.English:
			return language.English
		}
	}
	return language.Persian
}
