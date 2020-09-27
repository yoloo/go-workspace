package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"webservice/config"
	"webservice/util"
)

func getMyRedis() *redis.Client {
	return config.MyRedis
}

func getMyCache() *redis.Client {
	return config.MyCache
}

const (
	RolePointRetCode_NotExist  = "-20"
	RolePointRetCode_NotEnough = "-21"
)

// 获取用户积分
// 参数：charid=(\d+)
// 返回：retcode | {"points":xxxx}
func OnGetRolePoints(c *gin.Context) {
	if point, err := getMyRedis().Get("score:" + c.Query("charid")).Result(); err != nil {
		c.String(http.StatusOK, RolePointRetCode_NotExist+"|")
	} else {
		c.String(http.StatusOK, "0|"+util.JsonStringify(gin.H{"points": point}))
	}
}

// 扣除用户积分
// 参数：charid=(\d+)&points=(\d+)
// 返回：retcode | {"points":xxxx}
var reduceRolePointMtx sync.Mutex

func OnReduceRolePoints(c *gin.Context) {
	reduceRolePointMtx.Lock()
	defer reduceRolePointMtx.Unlock()

	if reply, err := getMyRedis().Get("score" + c.Query("charid")).Result(); err != nil {
		c.String(http.StatusOK, RolePointRetCode_NotExist+"|")
	} else {
		costpoint, _ := strconv.ParseInt(c.Query("points"), 10, 64)
		if points, _ := strconv.ParseInt(reply, 10, 64); points < costpoint {
			c.String(http.StatusOK, RolePointRetCode_NotEnough+"|")
			return
		}

		if leftpoints, err := getMyRedis().DecrBy("score:"+c.Query("charid"), int64(costpoint)).Result(); err != nil {
			c.String(http.StatusOK, RolePointRetCode_NotExist+"|")
		} else {
			c.String(http.StatusOK, "0|"+util.JsonStringify(gin.H{"points": leftpoints}))
			getMyRedis().Append("histroy:"+c.Query("charid"), fmt.Sprintf("[%d,%d,2],", time.Now().Unix(), leftpoints))
		}
	}
}

// 获取用户积分记录
// 参数：
// 		charid=(\d+)	角色ID
//		type=(\d+) 		0-所有日志 1-增加日志 2-消费日志
//		start=(\d+)		分页起始量
//		limit=(\d+)		每页显示的记录数
//
// 返回：retcode | [{"datetime":"YYYY-MM-DD hh:mm:ss","amount":xxx, "type":x}, ...]
// CURL：curl 'http://localhost:3000/getRolePointsLog?charid=1001&start=0&limit=2'
func OnGetRolePointsLog(c *gin.Context) {
	ntype, _ := strconv.ParseUint(c.Query("type"), 10, 32)
	nstart, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	nlimit, _ := strconv.ParseUint(c.Query("limit"), 10, 32)

	if reply, err := getMyRedis().Get("histroy:" + c.Query("charid")).Result(); err != nil {
		c.String(http.StatusOK, RolePointRetCode_NotExist+"|")
	} else {
		//reply [11,100,1],[12,100,2],
		reply = "[" + strings.TrimRight(reply, ",") + "]"

		var replylogs, vvlogs [][]uint64
		if e := json.Unmarshal([]byte(reply), &replylogs); e != nil || 0 == len(replylogs) {
			c.String(http.StatusOK, RolePointRetCode_NotExist+"|")
			return
		}

		for _, v := range replylogs {
			if len(v) >= 3 && (ntype == 0 || ntype == v[2]) {
				vvlogs = append(vvlogs, v)
			}
		}

		type RetObject struct {
			DateTime string `json:"datetime"`
			Amount   uint64 `json:"amount"`
			Type     int    `json:"type"`
		}

		var rets []RetObject
		if int(nstart) < len(vvlogs) {
			for i := int(nstart); i < len(vvlogs) && len(rets) < int(nlimit) && len(vvlogs[i]) >= 3; i++ {
				obj := &RetObject{
					DateTime: time.Unix(int64(vvlogs[i][0]), 0).Format("2006-01-02 15:04:05"),
					Amount:   vvlogs[i][1],
					Type:     int(vvlogs[i][2]),
				}
				rets = append(rets, *obj)
			}
		}

		c.String(http.StatusOK, "0|"+util.JsonStringify(rets))
	}
}
