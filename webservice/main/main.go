package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"net"
	"net/http"
	"strconv"
	"webservice/config"
	"webservice/routers"
)

func CheckIpLimit(ip net.IP) bool {
	for _, v := range config.GlobalConfig.Iplimit {
		if ip.String() == v.Ip {
			return true
		}
	}
	return false
}

type GlogLogger struct {
}

func (log *GlogLogger) Write(p []byte) (n int, err error) {
	glog.Info(string(p))
	return len(p), nil
}

func createGinRouterWithGlog() *gin.Engine {
	// export GIN_MODE=release or call gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)

	//replace gin.Default()
	router := gin.New()

	var logcfg = gin.LoggerConfig{
		Output:    &GlogLogger{},
		Formatter: nil,
	}
	router.Use(gin.LoggerWithConfig(logcfg), gin.RecoveryWithWriter(logcfg.Output))
	return router
}

// webservice -log_dir=.\log\
func main() {
	// glog need flag.Parse
	flag.Parse()

	router := createGinRouterWithGlog()
	defer glog.Flush()

	// load webservice config
	if err := config.LoadConfig(); err != nil {
		glog.Info("LoadConfig error", err)
		return
	}

	// access ip limit
	router.Use(func(c *gin.Context) {
		addr, _ := net.ResolveTCPAddr("tcp", c.Request.RemoteAddr)
		if CheckIpLimit(addr.IP) {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusOK, -2)
			glog.Info(addr, c.Request.URL.String(), "(invalid)")
		}
	})

	router.GET("/users", routers.OnUsers)
	router.GET("/getZT2Point", routers.OnGetZT2Point)
	router.GET("/getRolePoints", routers.OnGetRolePoints)
	router.GET("/reduceRolePoints", routers.OnReduceRolePoints)
	router.GET("/getRolePointsLog", routers.OnGetRolePointsLog)
	router.GET("/getBattleCamp", routers.OnGetBattleCamp)
	router.GET("/getCampCountry", routers.OnGetCampCountry)
	router.GET("/getCampGod", routers.OnGetCampGod)
	router.GET("/getCountryDareResult", routers.OnGetCountryDareResult)
	router.GET("/getGoldCountryPower", routers.OnGetGoldCountryPower)
	router.GET("/getGoldWorldKillNum", routers.OnGetGoldWorldKillNum)

	fmt.Println("start run zt2-webservice")

	glog.Info("Gin server listening on port:", config.GlobalConfig.ListenPort.Port)
	if err := router.Run(":" + strconv.Itoa(config.GlobalConfig.ListenPort.Port)); err != nil {
		glog.Info("Listen error ", err)
	}
}
