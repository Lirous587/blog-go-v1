package controller

import (
	"blog/logic"
	"blog/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func ResponseDataAboutIndexHandler(c *gin.Context) {
	//得到各大分类种类以及相应的名称
	var DataAboutIndex = models.DataAboutIndex{}
	if err := logic.ResponseDataAboutIndex(&DataAboutIndex); err != nil {
		zap.L().Error("logic.ResponseDataAboutIndex(&DataAboutIndex) failed,err:%v", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, DataAboutIndex)
}

func ResponseDataAboutEssayHandler(c *gin.Context) {
	//1.参数处理
	query := c.Query("id")
	id, err := strconv.Atoi(query)
	if err != nil {
		zap.L().Error("strconv.Atoi(query) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	//2.业务处理
	var essay = new(models.EssayData)
	if err = logic.GetEssayData(essay, id); err != nil {
		zap.L().Error("logic.GetEssayData(essay, id) failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, essay)
}
