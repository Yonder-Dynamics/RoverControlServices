package main

import(
  "time"
  "fmt"
  "html"
  "log"
  "net/http"
  "strings"
)

type GroundControlServer struct{}

func (f *GroundControlServer) ServeHTTP(w http.ResponseWriter, r *http.Request){
  log.Printf("Got a request on: %q\n", html.EscapeString(r.URL.Path))
  log.Printf("Method: %v\n",r.Method)
  switch r.Method{
    case "GET":
			f.handleGet(w,r)
    default:
      w.WriteHeader(501)
  }
}

func (f *GroundControlServer) handleGet(w http.ResponseWriter, r *http.Request){
  switch path := r.URL.Path; path{
		case "/test":
			fmt.Fprintf(w,"hello react!\n")
		default:
      f.serveFiles(w,r,path)
	}
}

func (f *GroundControlServer) serveFiles(w http.ResponseWriter, r *http.Request, path string){
  if(path == "/"){
    path = "index.html"
  }
  http.ServeFile(w,r,"base/" + strings.TrimLeft(path,"./"))
}

func main(){
  port := ":8080"
  fmt.Printf("Starting go server on port %v\n",port);
  var handler GroundControlServer
  s := &http.Server{
    Addr:           port,
    Handler:        &handler,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }
  log.Fatal(s.ListenAndServe())
  fmt.Printf("Shutting down go server...\n");
}
