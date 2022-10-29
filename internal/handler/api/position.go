package api

import (
	"github.com/convee/adcreative/internal/pkg/app"
	"github.com/convee/adcreative/internal/pkg/code"
	"github.com/convee/adcreative/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Position struct {
}
type BatchPosition struct {
	Position []string `json:"position"`
}

func (c *Position) List(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	publisher := ctx.Query("publisher")
	position := ctx.Query("position")
	positionService := &service.Position{
		Publisher: publisher,
		Position:  position,
	}
	data := positionService.GetApiList()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}
func (c *Position) Batch(ctx *gin.Context) {
	var (
		appG = app.Gin{C: ctx}
	)
	publisher := ctx.Query("publisher")
	var batchPosition BatchPosition
	errMsg := app.BindJson(ctx, &batchPosition)
	if len(errMsg) > 0 {
		appG.Response(http.StatusOK, code.INVALID_PARAMS, errMsg)
		return
	}
	data := make(map[string]interface{})
	var lists []interface{}
	for _, position := range batchPosition.Position {
		positionService := &service.Position{
			Publisher: publisher,
			Position:  position,
		}
		res := positionService.GetApiList()

		if list, ok := res["lists"].([]map[string]interface{}); ok {
			lists = append(lists, list[0])
		}
	}
	data["lists"] = lists
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}

func (c *Position) Material(ctx *gin.Context) {

	var (
		appG = app.Gin{C: ctx}
	)
	publisher := ctx.Query("publisher")
	position := ctx.Query("position")
	positionService := &service.Position{
		Publisher: publisher,
		Position:  position,
	}
	data := positionService.GetApiMaterial()
	appG.Response(http.StatusOK, code.SUCCESS, data)
	return
}
