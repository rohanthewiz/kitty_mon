package main

import (
	"github.com/julienschmidt/httprouter"
)

func doRoutes(router *httprouter.Router) {
	router.GET("/", Index)
	router.GET("/q/:query", QueryAll)
	//	router.GET("/ql", QueryLast)
	//	router.GET("/qi/:id", QueryId)
	router.GET("/api/v1", APIQuery)
	router.GET("/api/v1/l/:limit", APIQuery)
	router.GET("/json/q", QueryAsJson)
	//	router.GET("/json/qi/:id", QueryIdAsJson)
	router.GET("/q/:query/l/:limit", Query)
	// router.GET("/del/:id", WebDelete)
	router.GET("/del2weeks/:secret", WebDelete2Weeks)
	//router.GET("/js/:file", ServeJS)
}
