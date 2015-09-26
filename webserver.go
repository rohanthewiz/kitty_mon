package main
import (
	"net/http"
	"github.com/julienschmidt/httprouter"
//	"github.com/microcosm-cc/bluemonday"
	"fmt"
	"strconv"
	"path"
	"encoding/json"
)

// Good reading: http://www.alexedwards.net/blog/golang-response-snippets

func webserver(listen_port string) {
	router := httprouter.New()
	doRoutes(router)
	pf("Webserver listening on %s... Ctrl-C to quit\n", listen_port)
	lf(http.ListenAndServe(":" + listen_port, router))
}

// Handlers for httprouter
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/q/all/l/100", http.StatusFound)
//	fmt.Fprint(w, "Welcome!\n")
}

func Query(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	opts.Q = p.ByName("query")  // Overwrite the query param
	limit, err := strconv.Atoi(p.ByName("limit"))
	if err == nil {
		opts.L = limit
	}
	readings, err := getRecentReadings()
	RenderQuery(w, readings) //call Ego generated method
}

func QueryAll(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	opts.Q = p.ByName("query")  // Overwrite the query param
	limit, err := strconv.Atoi(p.ByName("limit"))
	if err == nil {
		opts.L = limit
	}
	readings, err := getAllReadings()
	RenderQuery(w, readings) //call Ego generated method
}

func QueryAsJson(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	readings, err := getRecentReadings()
	if err != nil {
		lpl(err)
		return
	}
	j_readings, _ := json.Marshal(readings)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j_readings)
}

//func QueryIdAsJson(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
//	resetOptions()
//	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
//	if err != nil { id = 0 }
//	opts_intf["qi"] = id  // qi is the highest priority
//	readings, err := getRecentReadings()
//	if err != nil {
//		lpl(err)
//		return
//	}
//	j_readings, _ := json.Marshal(readings)
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(j_readings)
//}

func WebDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		fpl("Error parsing reading id for delete.")
	} else {
		doDelete(find_reading_by_id(id))
	}
	http.Redirect(w, r, "/q/all/l/100", http.StatusFound)
}

func ServeJS(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.ServeFile(w, r, path.Join("js", p.ByName("file")))
}

func resetOptions() {
	opts.Ql = false // turn off unused option
	opts.L = -1 // turn off unused option
	opts.Q = "" // turn off higher priority option
}

func HandleRequestErr(err error, w http.ResponseWriter) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
}
