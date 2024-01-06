package handler

import (
	"discount/server"
	"discount/service/gift"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GiftHandler struct {
	gift *gift.Service
}

func NewGiftHandler(gift *gift.Service) GiftHandler {
	return GiftHandler{
		gift: gift,
	}
}

func SetupGiftRoutes(s *server.Server, h GiftHandler) {
	g := s.Engine.Group("/gift")
	g.POST("", h.InitGift)
	g.GET("/:giftCode", h.GetGift)
	g.POST("/use/:giftCode", h.UseGift)
}

// InitGift godoc
// @Summary			Initialize gift
// @Description		Initialize a new gift.
// @Tags			GiftDTO
// @Accept			json
// @Produce      	json
// @Param        body			body		gift.CreateRequest		true	"Gift init request"
// @Success      200			{object}	gift.DTO
// @Failure      	400  			{object}	Error
// @Failure      	500  			{object}  	Error
// @Router       	/gift		[post]
func (h GiftHandler) InitGift(ctx *gin.Context) {
	var req gift.CreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		handleError(ctx, err)
		return
	}
	result, err := h.gift.Create(&req)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// GetGift godoc
// @Summary      Get gift
// @Description  Get a gift by code.
// @Tags         GiftDTO
// @Accept       json
// @Produce      json
// @Param        giftCode		path		string				true	"Gift code"
// @Success      200			{object}	gift.DTO
// @Failure      400  			{object}	Error
// @Router       	/gift/{giftCode}		[get]
func (h GiftHandler) GetGift(ctx *gin.Context) {
	giftCode := ctx.Param("giftCode")
	if giftCode == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	result, err := h.gift.GetByCode(giftCode)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// UseGift godoc
// @Summary      Use gift
// @Description  Use a gift code.
// @Tags         GiftDTO
// @Accept       json
// @Produce      json
// @Param        giftCode		path		string				true	"Gift code"
// @Success      200			{object}	gift.DTO
// @Failure      400  			{object}	Error
// @Router       	/gift/use/{giftCode}		[post]
func (h GiftHandler) UseGift(ctx *gin.Context) {
	giftCode := ctx.Param("giftCode")
	if giftCode == "" {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	result, err := h.gift.UseGift(giftCode)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
