package handler

import (
	"discount/internal/locale"
	"discount/internal/serr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func getPaginationParams(c *gin.Context) (page, pageSize int) {
	if pageStr := c.Query("page"); pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if page <= 0 {
		page = 1
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}
	if pageSize == 0 {
		if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
			pageSize, _ = strconv.Atoi(pageSizeStr)
		}
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return
}

func getTraceID(ctx *gin.Context) string {
	rID, exist := ctx.Get("trace_id")
	if exist {
		switch rID.(type) {
		case string:
			return rID.(string)
		case []byte:
			return string(rID.([]byte))
		}
	}
	return ""
}

type Error struct {
	Message string         `json:"message"`
	Code    serr.ErrorCode `json:"code"`
	TraceID string         `json:"trace_id"`
}

func handleError(ctx *gin.Context, err error) {
	tID := getTraceID(ctx)
	lang := getLanguage(ctx)
	switch err.(type) {
	case *serr.ServiceError:
		var e *serr.ServiceError
		errors.As(err, &e)
		l := log.Error().Str("method", e.Method).Str("code", string(e.ErrorCode)).Str("trace_id", tID)
		if e.Cause != nil {
			l.Err(e.Cause)
		}
		l.Msg(e.Message)
		ctx.AbortWithStatusJSON(
			e.Code,
			Error{
				Message: locale.Localize(e.Message, lang),
				Code:    e.ErrorCode,
				TraceID: tID,
			},
		)
		return
	case *pq.Error:
		var e *pq.Error
		errors.As(err, &e)
		log.Error().Err(err).Str("trace_id", tID).Msg("pq error")
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			Error{Message: e.Detail, Code: "DB_ERROR", TraceID: tID},
		)
	default:
		log.Error().Err(err).Str("trace_id", getTraceID(ctx)).Msg("unknown error")
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			Error{Code: serr.ErrInternal, Message: "internal server error", TraceID: tID},
		)
		return
	}
}
