package routers

import (
	"net/http"
	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
)

func onQueryMyCacheInfo(c *gin.Context, reqfield,reqkey,retfiled string){
	var key = reqfield + ":" + c.Query(reqkey)
	if exist, err := getMyCache().Exists(key).Result(); err != nil || exist != 1 {
		c.JSON(http.StatusOK, gin.H{"ret":-1, "msg":"invalid access"})
	}else{ // exist == 1
		if reply, err := getMyCache().HGet(key, retfiled).Result(); err != nil{
			c.JSON(http.StatusOK, gin.H{"ret":0, retfiled:0})
		}else if gbk := mahonia.NewDecoder("gbk"); gbk != nil {
			c.JSON(http.StatusOK, gin.H{"ret":1, retfiled : gbk.ConvertString(reply)})
		}
	}
}

// 获取战区阵营
func OnGetBattleCamp(c *gin.Context) {
	onQueryMyCacheInfo(c, "zonebattle", "battleid", "members")
}

// 获取阵营国家
func OnGetCampCountry(c *gin.Context){
	onQueryMyCacheInfo(c, "campcountry", "campid", "members")
}

// 获取阵营大神
func OnGetCampGod(c *gin.Context) {
	onQueryMyCacheInfo(c, "campgod", "campid", "members")
}