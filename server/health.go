package server

import (
	"github.com/gin-gonic/gin"
)

// Health godoc
// @Summary Health check
// @Schemes
// @Description Health check
// @Tags Health
// @Accept json
// @Produce json
// @Success 200
// @Router /health [get]
func Health(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}
