package main
import (
	"net/http"
	"github.com/julienschmidt/httprouter"
//	"github.com/microcosm-cc/bluemonday"
	"log"
	"fmt"
	"io/ioutil"
	"strconv"
	"net/url"
	"path"
	"encoding/json"
)

// Good reading: http://www.alexedwards.net/blog/golang-response-snippets

func webserver(listen_port string) {
	router := httprouter.New()
	doRoutes(router)
	pf("Server listening on %s... Ctrl-C to quit", listen_port)
	log.Fatal(http.ListenAndServe(":" + listen_port, router))
}

// Handlers for httprouter
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/q/all/l/100", http.StatusFound)
//	fmt.Fprint(w, "Welcome!\n")
}

func Query(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	opts_str["q"] = p.ByName("query")  // Overwrite the query param
	limit, err := strconv.Atoi(p.ByName("limit"))
	if err == nil {
		opts_intf["l"] = limit
	}
	notes := queryNotes()
	RenderQuery(w, notes) //call Ego generated method
}
func QueryIdAsJson(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
	resetOptions()
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil { id = 0 }
	opts_intf["qi"] = id  // qi is the highest priority
	j_notes, err := json.Marshal(queryNotes())
	w.Header().Set("Content-Type", "application/json")
	w.Write(j_notes)
}

func WebDeleteNote(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, err := strconv.ParseInt(p.ByName("id"), 10, 64)
	if err != nil {
		fpl("Error deleting note.")
	} else {
		doDelete(find_note_by_id(id))
	}
	http.Redirect(w, r, "/q/all/l/100", http.StatusFound)
}

func ServeJS(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	http.ServeFile(w, r, path.Join("js", p.ByName("file")))
}

func resetOptions() {
	opts_intf["qi"] = nil // turn off unused option
	opts_intf["ql"] = false // turn off unused option
	opts_intf["l"] = -1 // turn off unused option
	opts_str["qg"] = "" // turn off higher priority option
	opts_str["qt"] = "" // turn off unused option
	opts_str["qd"] = "" // turn off unused option
	opts_str["qb"] = "" // turn off unused option
	opts_str["q"] = "" // turn off higher priority option
}

func HandleRequestErr(err error, w http.ResponseWriter) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
}
