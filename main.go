package main

import (
	"github.com/gin-gonic/gin"
)

var (
	conf *Conf
)

func getServiceOrFail(c *gin.Context, sn string) *Service {
	service, ok := conf.Services[sn]

	if !ok {
		c.JSON(200, gin.H{
			"message": "unknown service",
			"service": sn,
			"error":   true,
		})
		return nil
	}
	return service
}

func getStatus(sn string, s *Service) (status_code int, json gin.H) {
	status, err := s.Status()
	if err != nil {
		return 500, gin.H{
			"message": "error retrieving process status (wrong service configuration?)",
			"service": sn,
			"error":   true,
		}
	}

	return 200, gin.H{
		"running": status,
		"service": sn,
	}
}

func status(c *gin.Context) {
	sn := c.Param("name")
	service := getServiceOrFail(c, sn)

	if service == nil {
		return
	}
	sc, result := getStatus(sn, service)
	c.JSON(sc, result)
}

func list(c *gin.Context) {
	var out []gin.H
	for sn, service := range conf.Services {
		_, result := getStatus(sn, service)
		out = append(out, result)
	}
	c.JSON(200, out)
}

func change(c *gin.Context) {
	sn := c.Param("name")
	service := getServiceOrFail(c, sn)
	service.Start()

	c.JSON(200, gin.H{
		"message": "started",
	})
}

func main() {
	conf = parseConf("webctrl.yml")

	r := gin.Default()
	r.GET("/service/", list)
	r.GET("/service/:name", status)
	r.PUT("/service/:name", change)
	r.Run() // listen and serve on 0.0.0.0:8080
}
