/*
Author: zuoguocai@126.com
*/

package main

import (
    "io"
    "os"
    log "github.com/sirupsen/logrus"
    "net/http"
    "net/http/httputil"
    "strings"
    "go.elastic.co/apm/module/apmhttp"
)




func main() {
    fs := http.FileServer(http.Dir("./nginx/live2d"))
    mux := http.NewServeMux()
    mux.Handle("/live2d/", http.StripPrefix("/live2d/", fs))
    mux.HandleFunc("/", GetRealIP)
    log.Info("Server starting ...")
    
    
    var logFile string = "/var/log/ipcat/ipcat.log"
    var file, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Errorf("Could Not Open Log File : " + err.Error())
    }

    log.SetOutput(file)
    //log.SetFormatter(&log.TextFormatter{})
    log.SetFormatter(&log.JSONFormatter{})

    
    
    if err :=  http.ListenAndServe(":5000", apmhttp.Wrap(mux));err != nil {
          log.Errorf("Httpserver: ListenAndServe() error: %s", err)
    }
}

func GetRealIP(w http.ResponseWriter, r *http.Request) {
    dump, _ := httputil.DumpRequest(r, false)
    //log.Printf("%q\n", dump)
    log.WithFields(log.Fields{
	"ip": r.Header.Get("X-Real-Ip"),
        "url": r.URL.String()   ,
	"referer": r.Header.Get("Referer"),
        "method": r.Method,
	"userAgent": r.Header.Get("User-Agent"),
    }).Info()
    
    
    
    
    head := `<!doctype html><html lang="zh"><head><meta charset="UTF-8"><title>DevOps Pipeline Demo</title></head><body><h1 align="center" style="color:red;">ipcat v4.0</h1>`
    r1 := strings.Join([]string{"<h3 style='background-color:powderblue;'>","RemoteAddr:  ",r.RemoteAddr,"</h3>"},"")
    r2 := strings.Join([]string{"<h3 style='background-color:#DDA0DD;'>","X-Original-Forwarded-For:  ",r.Header.Get("X-Original-Forwarded-For"),"</h3>"},"")
    r3 := strings.Join([]string{"<h3 style='background-color:powderblue;'>","X-Forwarded-For:  ",r.Header.Get("X-Forwarded-For"),"</h3>"},"")
    r4 := strings.Join([]string{"<h3 style='background-color:#DDA0DD;'>","X-Real-Ip:  ",r.Header.Get("X-Real-Ip"),"</h3>"},"")

footer := `<script src="./live2d/L2Dwidget.min.js"></script><script type="text/javascript">L2Dwidget.init({
    model: {
      jsonPath: './live2d/tororo/assets/tororo.model.json',
    },
    display: {
      superSample: 2,
      width: 150,
      height: 150,
      position: 'right',
      hOffset: 0,
      vOffset: 0,
    },
    mobile: {
      show: true,
      scale: 1,
      motion: true,
    },
    react: {
      opacityDefault: 0.8,
      opacityOnHover: 0.2,
    }
  })
</script></html>`

    html := strings.Join([]string{head,r1,r2,r3,r4,footer},"")
    if _,err := io.WriteString(w, html); err != nil {
         log.Errorf("Httpserver: io WriteString error: %s", err)
    }
}
