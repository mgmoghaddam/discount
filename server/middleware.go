package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func WithTraceID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, exist := ctx.Get("trace_id")
		if !exist {
			id, err := uuid.NewRandom()
			if err != nil {
				log.Error().Str("method", "server.WithTraceID").Err(err).Msg("failed to create uuid")
			} else {
				ctx.Set("trace_id", id.String())
			}
		}
		ctx.Next()
	}
}
