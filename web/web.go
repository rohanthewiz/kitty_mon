package web

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"kitty_mon/config"
	"kitty_mon/reading"
	"kitty_mon/util"
	"net/http"
	"strconv"
)

func DoRoutes(router *httprouter.Router) {
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

func Webserver(listen_port string) {
	router := httprouter.New()
	DoRoutes(router)
	util.Pf("Webserver listening on %s... Ctrl-C to quit\n", listen_port)
	util.Lf(http.ListenAndServe(":"+listen_port, router))
}

// Handlers for httprouter
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/q/all/l/100", http.StatusFound)
	//	fmt.Fprint(w, "Welcome!\n")
}

func Query(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	config.Opts.Q = p.ByName("query") // Overwrite the query param
	limit, err := strconv.Atoi(p.ByName("limit"))
	if err == nil {
		config.Opts.L = limit
	}
	readings, err := reading.GetRecentReadings()
	RenderQuery(w, readings) //call Ego generated method
}

func QueryAll(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	config.Opts.Q = p.ByName("query") // Overwrite the query param
	limit, err := strconv.Atoi(p.ByName("limit"))
	if err == nil {
		config.Opts.L = limit
	}
	readings, err := reading.GetAllReadings()
	RenderQuery(w, readings) //call Ego generated method
}

func QueryAsJson(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	readings, err := reading.GetRecentReadings()
	if err != nil {
		util.Lpl(err)
		return
	}
	j_readings, _ := json.Marshal(readings)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j_readings)
}

func APIQuery(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	limit, err := strconv.Atoi(p.ByName("limit"))
	if err != nil {
		limit = 0
	}
	readings, err := reading.QueryReadings(limit)
	if err != nil {
		util.Lpl(err)
		return
	}

	readingsWithNames := reading.ReadingsPopulateNodeName(readings)
	jReadings, _ := json.Marshal(readingsWithNames)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jReadings)
}

func WebDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		util.Fpl("Error parsing reading id for delete.")
	} else {
		reading.DoDelete(reading.Find_reading_by_id(id))
	}
	http.Redirect(w, r, "/q/all/l/100", http.StatusFound)
}

func WebDelete2Weeks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if p.ByName("secret") == "gonotes" {
		reading.DeleteGt2weeks()
		util.Pl("Delete should be completed at this point")
	}
	http.Redirect(w, r, "/q/all/l/100", http.StatusFound)
}

func resetOptions() {
	config.Opts.Ql = false // turn off unused option
	config.Opts.L = -1     // turn off unused option
	config.Opts.Q = ""     // turn off higher priority option
}
