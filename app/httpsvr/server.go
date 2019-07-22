package httpsvr

import (
    "app/httpsvr/handler"
    "fmt"
    "net/http"
    "time"

    "github.com/cihub/seelog"
    "github.com/gorilla/mux"
)

/*
StartHTTP func()
*/
func StartHTTP(httpPort int, writeTimeout int, readTimeout int) error {
    r := mux.NewRouter().StrictSlash(true)
    initRoutes(r)
    seelog.Info("Initialize HTTP Routers Success !")

    svr := &http.Server{
        Handler:      r,
        Addr:         fmt.Sprintf(":%v", httpPort),
        WriteTimeout: time.Duration(writeTimeout) * time.Second,
        ReadTimeout:  time.Duration(readTimeout) * time.Second,
    }

    seelog.Info("Listen HTTP Port And Serve ...")
    // err := http.ListenAndServe(":8888", r)
    err := svr.ListenAndServe()

    return err
}

func initRoutes(r *mux.Router) {

    r.HandleFunc("/exec", handler.ExecuteHandler).Methods("POST")
    r.HandleFunc("/cron", handler.CronAddHandler).Methods("POST")
    r.HandleFunc("/cron", handler.CronDeleteHandler).Methods("DELETE")
    r.HandleFunc("/cron", handler.CronUpdateHandler).Methods("PUT")
    r.HandleFunc("/test", handler.TestHandler).Methods("GET")
    //r.HandleFunc("/test/ajax", handler.TestAjaxHandler).Methods("POST")

    // static resource router
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

    r.HandleFunc("/{key}", handler.NotFoundHandler)
    /*
    	备查
    	// normal router
    	r.HandleFunc("/index", handler.IndexHandler)

    	// html router
    	h := r.PathPrefix("/html").Subrouter()
    	h.HandleFunc("/", handler.NotFoundHandler)
    	h.HandleFunc("/{key}", handler.HTMLHandler)

    	// ajax router
    	a := r.PathPrefix("/ajax").Subrouter()
    	a.HandleFunc("/", handler.NotFoundHandler)
    	a.HandleFunc("/{key}", handler.AjaxHandler)

    	// test router
    	t := r.PathPrefix("/test").Subrouter()
    	t.HandleFunc("/", handler.NotFoundHandler)
    	t.HandleFunc("/{key}", handler.TestHandler)
    	t.HandleFunc("/ajax/{key}", handler.TestAjaxHandler)

    	// static resource router
    	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    	// other root subrouter
    	root := r.PathPrefix("/").Subrouter()
    	root.HandleFunc("/", handler.IndexHandler)
    	root.HandleFunc("/{key}", handler.NotFoundHandler)
    */
}
