package controller

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liuhongdi/digv10/global"
	"github.com/liuhongdi/digv10/pkg/result"
)

type SetController struct{}

func NewSetController() SetController {
	return SetController{}
}

//发布一条消息到redis
func (u *SetController) Pub(c *gin.Context) {
	resultRes := result.NewResult(c)
	data:=c.Query("id")

	err := global.RedisDb.Publish("articleMsg", data).Err()
	if err != nil {
		//return errors.New("发布失败")
		//result.Error(400,errs.Error())
		resultRes.Error(400,err.Error())
		return
	}
}

