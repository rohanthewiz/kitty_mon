package main
import (
"fmt"
"io"
"strconv"
)
//line query.ego:1
 func RenderQuery(w io.Writer, readings []Reading) error  {
//line query.ego:2
_, _ = fmt.Fprintf(w, "\n")
//line query.ego:3
_, _ = fmt.Fprintf(w, "\n\n<html>\n<head>\n  <style>\n    body { background-color: tan }\n    ul { list-style-type:none; margin: 0; padding: 0; }\n    ul.topmost > li:first-child { border-top: 1px solid #531C1C}\n    ul.topmost > li { border-top:none; border-bottom: 1px solid #8A2E2E; padding: 0.3em 0.3em}\n    li { border-top: 1px solid #B89c72; line-height:1.2em; padding: 1.2em, 4em }\n    .h1 { font-size: 1.2em; margin-bottom: 0.1em; padding: 0.1em }\n    .title { font-weight: bold; color:darkgreen; padding-top: 0.4em }\n    .count { font-size: 0.8em; color:#401020; padding-left: 0.5em; padding-right: 0.5em }\n    .tool { font-size: 0.7em; color:#401020; padding-left: 0.5em }\n    code { -webkit-border-radius: 0.3em;\n          -moz-border-radius: 0.3em;\n          border-radius: 0.3em; }\n  </style>\n  <script type=\"text/javascript\" src=\"https://code.jquery.com/jquery-2.1.3.min.js\"></script>\n</head>\n<body>\n")
//line query.ego:23
 readings_count := len(readings) 
//line query.ego:24
_, _ = fmt.Fprintf(w, "\n<p>\n  <span class=\"h1\">KittyMon</span> <span class=\"count\">")
//line query.ego:25
_, _ = fmt.Fprintf(w, "%v",  readings_count )
//line query.ego:25
_, _ = fmt.Fprintf(w, " found</span>  [<a class=\"tool\" href=\"http://127.0.0.1:")
//line query.ego:25
_, _ = fmt.Fprintf(w, "%v",  opts_str["port"] )
//line query.ego:25
_, _ = fmt.Fprintf(w, "/new\">New</a> |\n  <a class=\"tool\" href=\"http://127.0.0.1:")
//line query.ego:26
_, _ = fmt.Fprintf(w, "%v",  opts_str["port"] )
//line query.ego:26
_, _ = fmt.Fprintf(w, "/q/all\">All</a>]\n</p>\n\n<ul class=\"topmost\">\n  ")
//line query.ego:30
 for _, reading := range readings { 
//line query.ego:31
_, _ = fmt.Fprintf(w, "\n      ")
//line query.ego:31
 id_str := strconv.FormatInt(reading.Id, 10) 
//line query.ego:32
_, _ = fmt.Fprintf(w, "\n      <li><a class=\"title\" href=\"/show/")
//line query.ego:32
_, _ = fmt.Fprintf(w, "%v",  id_str )
//line query.ego:32
_, _ = fmt.Fprintf(w, "\">")
//line query.ego:32
_, _ = fmt.Fprintf(w, "%v",  reading.CreatedAt )
//line query.ego:32
_, _ = fmt.Fprintf(w, " - ")
//line query.ego:32
_, _ = fmt.Fprintf(w, "%v",  short_sha(reading.SourceGuid) )
//line query.ego:32
_, _ = fmt.Fprintf(w, "</a>\n      ")
//line query.ego:33
 if readings_count == 1 { 
//line query.ego:34
_, _ = fmt.Fprintf(w, "\n      | <a class=\"tool\" href=\"/del/")
//line query.ego:34
_, _ = fmt.Fprintf(w, "%v",  id_str )
//line query.ego:34
_, _ = fmt.Fprintf(w, "\"\n            onclick=\"return confirm('Are you sure you want to delete this reading?')\">\n        delete</a>\n      ")
//line query.ego:37
 } 
//line query.ego:38
_, _ = fmt.Fprintf(w, "\n      ")
//line query.ego:38
 if reading.Temp != -255 { 
//line query.ego:38
_, _ = fmt.Fprintf(w, " - ")
//line query.ego:38
_, _ = fmt.Fprintf(w, "%v",  reading.Temp )
//line query.ego:39
_, _ = fmt.Fprintf(w, "\n      ")
//line query.ego:39
 } 
//line query.ego:40
_, _ = fmt.Fprintf(w, "\n      ")
//line query.ego:40
 if reading.IPs != "" { 
//line query.ego:40
_, _ = fmt.Fprintf(w, " - ")
//line query.ego:40
_, _ = fmt.Fprintf(w, "%v",  reading.IPs )
//line query.ego:41
_, _ = fmt.Fprintf(w, "\n      ")
//line query.ego:41
 } 
//line query.ego:42
_, _ = fmt.Fprintf(w, "\n      </li>\n  ")
//line query.ego:43
 } 
//line query.ego:44
_, _ = fmt.Fprintf(w, "\n</ul>\n\n<script type=\"text/javascript\">\n  $( function() {\n    //el = $('.reading-body');\n    //el.find(\"pre code\").each( function(i, block) {\n    //  hljs.highlightBlock( block );\n    //})\n  });\n</script>\n\n</body>\n</html>\n")
return nil
}
