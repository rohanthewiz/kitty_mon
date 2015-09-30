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
_, _ = fmt.Fprintf(w, "\n<!DOCTYPE html>\n<html>\n<head>\n  <style>\n    @import url(https://fonts.googleapis.com/css?family=Arimo:400,700);\n    * {position:relative;margin:0;padding:0; border: 0 none;}\n    *, *:before, *:after {box-sizing: inherit;}\n    html{\n        box-sizing: border-box;\n        font-family: Arimo,sans-serif;\n        font-size: 0.9rem;\n        color: #e6eff0;\n        background: #32382b;\n    }\n    table * {height: auto; min-height: none} /* fixed ie9 & <*/\n    table {\n        background: #bec28b;\n        table-layout: fixed;\n        margin: 2px auto;\n        width: 98%%;\n        box-shadow: 0 0 3px 1px rgba(0,0,0,.4);\n        border-collapse: collapse;\n        border: 1px solid rgba(0,0,0,.5);\n        border-top: 0 none;\n    }\n    thead {\n        background-color: #788056;\n        text-align: center;\n        z-index: 2;\n    }\n    thead tr {\n        padding-right: 16px;\n        box-shadow: 0 0px 1px rgba(0,0,0,0.1);\n        z-index: 2;\n    }\n    th {\n        border-right: 2px solid rgba(0,0,0,.2);\n        padding: .7rem 0;\n        font-size: 1.1rem;\n        font-weight: normal;\n        font-variant: small-caps;\n    }\n    tbody {\n        display: block;\n        height: calc(100vh - 5rem);\n        min-height: calc(200px + 1px); /* ie9 and < fix */\n        overflow-Y: scroll;\n        color: #000;\n    }\n    tr {\n        display: block;\n        overflow: hidden;\n    }\n    tbody tr:nth-child(odd) {\n        background: rgba(0,0,0,.1);\n    }\n    th, td {\n        width: 22%%;\n        float:left;\n    }\n    td {\n        padding: .5rem 0 .5rem 1rem;\n        border-right: 2px solid rgba(0,0,0,.2);\n    }\n    td:nth-child(2n) {color: #fff;}\n    th:last-child, td:last-child {\n        text-align: center;\n        border-right: 0 none;\n        padding-left: 0;\n    }\n    .h1 { font-size: 1.2em; margin-bottom: 0.1em; padding: 0.1em }\n    .temp {width: 9%%; font-weight: bold; }\n    .temp_fail { color:red; }\n    .temp_hot { color:#a07030; }\n    .temp_warm { color:yellow; }\n    .temp_good { color:darkgreen; }\n    .created_at { width: 25%% }\n    .title { font-weight: bold; color:darkgreen; padding-top: 0.4em }\n    .count { font-size: 0.8em; color:#f0e9e6; padding-left: 0.5em; padding-right: 0.5em }\n    .tool { font-size: 0.7em; color:#401020; padding-left: 0.5em }\n    .heading { margin: 1px auto; padding-left:2px; width: 98%% }\n  </style>\n  <!-- <script type=\"text/javascript\" src=\"https://code.jquery.com/jquery-2.1.3.min.js\"></script> -->\n</head>\n<body>\n")
//line query.ego:88
 readings_count := len(readings) 
//line query.ego:89
_, _ = fmt.Fprintf(w, "\n<div class=\"heading\">\n  <span class=\"h1\">KittyMon</span> <span class=\"count\">")
//line query.ego:90
_, _ = fmt.Fprintf(w, "%v",  readings_count )
//line query.ego:90
_, _ = fmt.Fprintf(w, " found</span>  <a class=\"tool\" href=\"/q/all\">All</a>\n</div>\n\n<table>\n  <thead>\n  <tr>\n  <th class=\"temp\">Temp</th><th class=\"created_at\">Created At</th><th>Reading Guid</th><th>Source Guid</th><th>Source IP</th>\n      ")
//line query.ego:97
 if readings_count == 1 { 
//line query.ego:98
_, _ = fmt.Fprintf(w, "\n        <th>Delete</th>\n      ")
//line query.ego:99
 } 
//line query.ego:100
_, _ = fmt.Fprintf(w, "\n  </tr>\n  </thead>\n  <tbody>\n  ")
//line query.ego:103
 var temp_class string
  var temp int 
//line query.ego:105
_, _ = fmt.Fprintf(w, "\n  ")
//line query.ego:105
 for _, reading := range readings { 
//line query.ego:106
_, _ = fmt.Fprintf(w, "\n      ")
//line query.ego:107

        temp = reading.Temp /1000
        switch {
        case temp > 64:
          temp_class = "temp_fail"
        case temp > 62:
          temp_class = "temp_hot"
        case temp > 59:
          temp_class = "temp_warm"
        default:
          temp_class = "temp_good"
        }
      
//line query.ego:120
_, _ = fmt.Fprintf(w, "\n      <tr>\n\n      ")
//line query.ego:122
 if reading.Temp == -255 { 
//line query.ego:123
_, _ = fmt.Fprintf(w, "\n       <td>N/A</td>\n      ")
//line query.ego:124
 } else { 
//line query.ego:125
_, _ = fmt.Fprintf(w, "\n       <td class=\"temp ")
//line query.ego:125
_, _ = fmt.Fprintf(w, "%v", temp_class )
//line query.ego:125
_, _ = fmt.Fprintf(w, "\">")
//line query.ego:125
_, _ = fmt.Fprintf(w, "%v",  temp )
//line query.ego:125
_, _ = fmt.Fprintf(w, "</td>\n      ")
//line query.ego:126
 } 
//line query.ego:127
_, _ = fmt.Fprintf(w, "\n\n      <td class=\"created_at\">")
//line query.ego:128
_, _ = fmt.Fprintf(w, "%v",  reading.CreatedAt.Format("15:04:05 Mon Jan _2") )
//line query.ego:128
_, _ = fmt.Fprintf(w, "</td>\n      <td>")
//line query.ego:129
_, _ = fmt.Fprintf(w, "%v",  short_sha(reading.Guid) )
//line query.ego:129
_, _ = fmt.Fprintf(w, "</td>\n      <td>")
//line query.ego:130
_, _ = fmt.Fprintf(w, "%v",  short_sha(reading.SourceGuid) )
//line query.ego:130
_, _ = fmt.Fprintf(w, "</td>\n      <td>")
//line query.ego:131
_, _ = fmt.Fprintf(w, "%v",  reading.IPs )
//line query.ego:131
_, _ = fmt.Fprintf(w, "</td>\n      ")
//line query.ego:132
 if readings_count == 1 { 
//line query.ego:133
_, _ = fmt.Fprintf(w, "\n      ")
//line query.ego:133
 id_str := strconv.FormatInt(reading.Id, 10) 
//line query.ego:134
_, _ = fmt.Fprintf(w, "\n      <td><a class=\"tool\" href=\"/del/")
//line query.ego:134
_, _ = fmt.Fprintf(w, "%v",  id_str )
//line query.ego:134
_, _ = fmt.Fprintf(w, "\"\n            onclick=\"return confirm('Are you sure you want to delete this reading?')\">\n        delete</a>\n      </td>\n      ")
//line query.ego:138
 } 
//line query.ego:139
_, _ = fmt.Fprintf(w, "\n\n      </tr>\n  ")
//line query.ego:141
 } 
//line query.ego:142
_, _ = fmt.Fprintf(w, "\n  </tbody>\n</table>\n\n<script type=\"text/javascript\">\n  //$( function() {\n    //el = $('.reading-body');\n    //el.find(\"pre code\").each( function(i, block) {\n    //  hljs.highlightBlock( block );\n    //})\n  //});\n</script>\n\n</body>\n</html>\n")
return nil
}
