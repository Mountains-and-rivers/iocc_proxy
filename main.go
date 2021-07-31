package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/url"
)

func proxy_iocc()  {
	// Echo instance
	e := echo.New()
	// Routes
	//e.Static("/static", "static")
	//e.File("/", "templates/login.html")
	// Setup proxy
	url1, err := url.Parse("http://localhost:9001")
	if err != nil {
		e.Logger.Fatal(err)
	}
	//url2, err := url.Parse("http://localhost:8082")
	//if err != nil {
	//	e.Logger.Fatal(err)
	//}
	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
		//{
		//	URL: url2,
		//},
	}
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func proxy_iocc1()  {
	// Echo instance
	e := echo.New()
	// Routes
	e.Static("/static", "static")
	e.File("/1", "templates/login.html")
	// Setup proxy
	url1, err := url.Parse("http://localhost:9001")
	if err != nil {
		e.Logger.Fatal(err)
	}
	//url2, err := url.Parse("http://localhost:8082")
	//if err != nil {
	//	e.Logger.Fatal(err)
	//}
	targets := []*middleware.ProxyTarget{
		{
			URL: url1,
		},
		//{
		//	URL: url2,
		//},
	}
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	// Start server
	e.Logger.Fatal(e.Start(":8001"))
}

func main() {
	go proxy_iocc()
	go proxy_iocc1()
	e := echo.New()
	// Routes
	e.Static("/static", "static")
	e.File("/1", "templates/login.html")
	e.Logger.Fatal(e.Start(":2323"))
}