package jgin

import (
	"github.com/gin-gonic/gin"
	"github.com/tommy351/gin-sessions"
)

func SetSession(ctx *gin.Context, k string, o interface{}) {
	session := sessions.Get(ctx)
	session.Set(k, o)
	session.Save()
}

func GetSession(ctx *gin.Context, k string) interface{} {
	session := sessions.Get(ctx)
	return session.Get(k)
}

func ClearAllSession(ctx *gin.Context) {
	session := sessions.Get(ctx)
	session.Clear()
	session.Save()
	return
}
