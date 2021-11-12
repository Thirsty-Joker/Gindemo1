package controller

import (
	"demo1/src/model"
	"demo1/src/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Config(c *gin.Context) {
	action := c.Param("action")
	if action == "" {
		c.String(http.StatusBadRequest, "params failed")
		return
	} else if action == "add" {
		AddConfig(c)
	} else if action == "remove" {
		RmConfig(c)
	} else if action == "edit" {
		EditConfig(c)
	}
}

func AddConfig(c *gin.Context) {
	var config model.Rule
	if err := c.ShouldBind(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if service.CheckConfig(config) {
		modelID := model.AddRule(config)
		fmt.Printf("%d", modelID)
		if modelID != 0 {
			c.JSON(http.StatusOK, gin.H{"message": "please remember this id", "rule_id": modelID})
			return
		}
	}
	c.String(http.StatusBadRequest, "add config fail")
}

func RmConfig(c *gin.Context) {
	//todo
	c.JSON(http.StatusOK, gin.H{"message": "hello Rm"})
}
func EditConfig(c *gin.Context) {
	//todo
	c.JSON(http.StatusOK, gin.H{"message": "hello Edit"})
}

func Check(c *gin.Context) {
	var form model.Client
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	downloadUrl, updateVersionCode, md5, title, updateTips := service.Hit(&form)
	c.JSON(http.StatusOK, gin.H{
		"download_url:":       downloadUrl,
		"update_version_code": updateVersionCode,
		"md5":                 md5,
		"title":               title,
		"update_tips":         updateTips,
	})
}
