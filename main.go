package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"net/url"
)

func proxy_iocc(port string,proxyurl string) {
	e := echo.New()
	url1, err := url.Parse(proxyurl)
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
	}
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	fmt.Println(port)
	e.Logger.Fatal(e.Start(":" + port))
}

//go:embed static templates
var local embed.FS
//go:embed conf/config.json
var s []byte
func main() {
	//type hostconfig struct {
	//	IP string `json:"ip"`
	//	//HOSTNAME string `json:"hostname"`
	//}
	//
	//type hostconfigs struct {
	//	IP1   map[string]hostconfig   `json:"ip1"`
	//	IP2   map[string]hostconfig   `json:"ip2"`
	//	IP3   map[string]hostconfig   `json:"ip3"`
	//}
	//var hconfig hostconfigs
	//portArr := []string{"8001", "8002", "8003", "8004", "8005","8006"}
	//plan, err := ioutil.ReadFile(s)
	//if err != nil {
	//	fmt.Println(err)
	//}
	buf := bytes.NewBuffer(s)
	res, err := simplejson.NewFromReader(buf)
	if err != nil || res == nil {
		log.Fatal("something wrong when call NewFromReader")
	}
	confArr, _ := res.Get("proxy_arr").Array()
	for _, row := range confArr {
		if each_map, ok := row.(map[string]interface{}); ok {
			proxy_IP := each_map["IP"].(string)
			proxy_port := each_map["proxy_port"].(string)
			go proxy_iocc(proxy_port,proxy_IP)
		}
	}
		e := echo.New()
	    e.GET("/static", echo.WrapHandler(http.FileServer(http.FS(local))))
		e.File("/", "templates/login.html")
		e.GET("/config", func(c echo.Context) error {
			u := res
			return c.JSON(200,u)
		})
		e.Logger.Fatal(e.Start(":2323"))
}