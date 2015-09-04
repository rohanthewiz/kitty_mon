// Generated by ego on Fri Sep  4 00:03:57 2015.
// DO NOT EDIT

package main
import (
"fmt"
"html"
"io"
"strconv"
)
//line templates/query.ego:1
 func RenderQuery(w io.Writer, readings []Reading) error  {
//line templates/query.ego:2
_, _ = fmt.Fprint(w, "\n")
//line templates/query.ego:3
_, _ = fmt.Fprint(w, "\n\n<html>\n<head>\n  <style>\n    body { background-color: tan }\n    ul { list-style-type:none; margin: 0; padding: 0; }\n    ul.topmost > li:first-child { border-top: 1px solid #531C1C}\n    ul.topmost > li { border-top:none; border-bottom: 1px solid #8A2E2E; padding: 0.3em 0.3em}\n    li { border-top: 1px solid #B89c72; line-height:1.2em; padding: 1.2em, 4em }\n    .h1 { font-size: 1.2em; margin-bottom: 0.1em; padding: 0.1em }\n    .title { font-weight: bold; color:darkgreen; padding-top: 0.4em }\n    .count { font-size: 0.8em; color:#401020; padding-left: 0.5em; padding-right: 0.5em }\n    .tool { font-size: 0.7em; color:#401020; padding-left: 0.5em }\n    code { -webkit-border-radius: 0.3em;\n          -moz-border-radius: 0.3em;\n          border-radius: 0.3em; }\n  </style>\n  <script type=\"text/javascript\" src=\"https://code.jquery.com/jquery-2.1.3.min.js\"></script>\n</head>\n<body>\n")
//line templates/query.ego:23
 readings_count := len(readings) 
//line templates/query.ego:24
_, _ = fmt.Fprint(w, "\n<p>\n  <span class=\"h1\">KittyMon</span> <span class=\"count\">")
//line templates/query.ego:25
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  readings_count )))
//line templates/query.ego:25
_, _ = fmt.Fprint(w, " found</span>  [<a class=\"tool\" href=\"http://127.0.0.1:")
//line templates/query.ego:25
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  opts_str["port"] )))
//line templates/query.ego:25
_, _ = fmt.Fprint(w, "/new\">New</a> |\n  <a class=\"tool\" href=\"http://127.0.0.1:")
//line templates/query.ego:26
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  opts_str["port"] )))
//line templates/query.ego:26
_, _ = fmt.Fprint(w, "/q/all\">All</a>]\n</p>\n\n<ul class=\"topmost\">\n  ")
//line templates/query.ego:30
 for _, reading := range readings { 
//line templates/query.ego:31
_, _ = fmt.Fprint(w, "\n      ")
//line templates/query.ego:31
 id_str := strconv.FormatInt(reading.Id, 10) 
//line templates/query.ego:32
_, _ = fmt.Fprint(w, "\n      <li><a class=\"title\" href=\"/show/")
//line templates/query.ego:32
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  id_str )))
//line templates/query.ego:32
_, _ = fmt.Fprint(w, "\">")
//line templates/query.ego:32
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  reading.CreatedAt )))
//line templates/query.ego:32
_, _ = fmt.Fprint(w, " - ")
//line templates/query.ego:32
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  short_sha(reading.SourceGuid) )))
//line templates/query.ego:32
_, _ = fmt.Fprint(w, "</a>\n      ")
//line templates/query.ego:33
 if readings_count == 1 { 
//line templates/query.ego:34
_, _ = fmt.Fprint(w, "\n      | <a class=\"tool\" href=\"/del/")
//line templates/query.ego:34
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  id_str )))
//line templates/query.ego:34
_, _ = fmt.Fprint(w, "\"\n            onclick=\"return confirm('Are you sure you want to delete this reading?')\">\n        delete</a>\n      ")
//line templates/query.ego:37
 } 
//line templates/query.ego:38
_, _ = fmt.Fprint(w, "\n      ")
//line templates/query.ego:38
 if reading.Temp != -255 { 
//line templates/query.ego:38
_, _ = fmt.Fprint(w, " - ")
//line templates/query.ego:38
_, _ = fmt.Fprint(w, html.EscapeString(fmt.Sprintf("%v",  reading.Temp )))
//line templates/query.ego:39
_, _ = fmt.Fprint(w, "\n      ")
//line templates/query.ego:39
 } 
//line templates/query.ego:40
_, _ = fmt.Fprint(w, "\n      </li>\n  ")
//line templates/query.ego:41
 } 
//line templates/query.ego:42
_, _ = fmt.Fprint(w, "\n</ul>\n\n<script type=\"text/javascript\">\n  $( function() {\n    //el = $('.reading-body');\n    //el.find(\"pre code\").each( function(i, block) {\n    //  hljs.highlightBlock( block );\n    //})\n  });\n</script>\n\n</body>\n</html>\n")
return nil
}
