package routers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func OnUsers(c *gin.Context){
	c.String(http.StatusOK, "respond with a resource")
}

// 获取征途2点数
// 参数：zoneid=(\d+)&charid=(\d+)
// 返回：
// 		正常：{"ret":0, "point":(\d+)}
//		无效角色：{"ret":-10}
func OnGetZT2Point(c *gin.Context){
	var charid = c.Query("charid")
	if 0 == len(charid) {
		c.JSON(http.StatusOK, gin.H{"ret":-10})
		return
	}
	var key = "user:" + charid
	if exist, err := getMyCache().Exists(key).Result(); err != nil || exist == 0 {
		c.JSON(http.StatusOK, gin.H{"ret":-10})
		return
	}

	if point, err := getMyCache().HGet(key, "zt2point").Result(); err != nil {
		c.JSON(http.StatusOK, gin.H{"ret":-10})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"ret" : 0,
			"point" : point,
		})
		return
	}
}
