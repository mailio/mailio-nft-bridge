package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mailio/mailio-nft-server/model"
)

type MailioUserStatsAPI struct {
}

func NewMailioUserStatsAPI() *MailioUserStatsAPI {
	return &MailioUserStatsAPI{}
}

// Get Mailio User Stats
// @Summary      Get Mailio User Stats
// @Description  Caliing Mailio server to retrieve user stats
// @Tags         Stats
// @Param        mailioaddress  path      string  true  "mailioaddress"
// @Success      200            {object}  model.MailioUserStats
// @Failure      404            {object}  api.JSONError  "user not found"
// @Failure      500            {object}  api.JSONError  "internal server error"
// @Accept       json
// @Produce      json
// @Router       /v1/user/{mailioaddress}/stats [get]
func (msa *MailioUserStatsAPI) GetMailioUserStats(c *gin.Context) {
	stats := model.MailioUserStats{
		Address: c.Param("mailioaddress"),
	}
	c.JSON(http.StatusOK, stats)
}
