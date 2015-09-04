package main

import (
	"github.com/julienschmidt/httprouter"
)

func doRoutes(router * httprouter.Router) {
	router.GET("/", Index)
//	router.GET("/ql", QueryLast)
//	router.GET("/qi/:id", QueryId)
	router.GET("/json/qi/:id", QueryIdAsJson)
	router.GET("/q/:query/l/:limit", Query)
	router.GET("/del/:id", WebDelete)
	router.GET("/js/:file", ServeJS)
}
