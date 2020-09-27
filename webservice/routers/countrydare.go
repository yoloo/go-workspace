package routers

import (
	"github.com/gin-gonic/gin"
)

// 获取黄金国战联赛战斗结果
func OnGetCountryDareResult(c *gin.Context) {
	onQueryMyCacheInfo(c, "goldcountrydare", "zoneid", "dareresults")
}

// 获取各个国战电竞类型比赛结束时的各国实力
func OnGetGoldCountryPower(c *gin.Context) {
	onQueryMyCacheInfo(c, "goldcountrypower", "zoneid", "countrypowers")
}

// 获取黄金国战联赛世界杀敌榜
func OnGetGoldWorldKillNum(c *gin.Context) {
	onQueryMyCacheInfo(c, "goldcountrydareworldkillinfo", "zoneid", "worldkillinfo")
}
