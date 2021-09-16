package controller

import (
	"github.com/Luna-CY/v2ray-helper/common/deploy"
	"github.com/Luna-CY/v2ray-helper/common/http/code"
	"github.com/Luna-CY/v2ray-helper/common/http/response"
	"github.com/gin-gonic/gin"
)

func V2rayServerDeployStage(c *gin.Context) {
	stages := deploy.GetManager().GetStage()
	if nil == stages {
		response.Response(c, code.BadRequest, "没有部署中的任务", nil)
	}

	stageList := response.NewEmptyJsonList()
	for _, stage := range stages {
		stageList = append(stageList, gin.H{"stage": stage.Stage, "state": stage.State})
	}

	last := stages[len(stages)-1]
	if deploy.StageConfigGenerate == last.Stage && deploy.StateRunning != last.State {
		deploy.GetManager().Clean()
	}

	response.Success(c, code.OK, &gin.H{"stages": stageList, "cloudreve": gin.H{"email": "admin@cloudreve.org", "password": deploy.GetManager().GetCloudreveAdminPassword()}})
}
