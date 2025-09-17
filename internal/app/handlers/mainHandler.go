package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexPageHandler(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "index.html", nil)
}
