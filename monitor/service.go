package monitor

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"
)

type Monitor struct {
	httpServer *http.Server
	httpMux    *http.ServeMux
	httpPort   int
}

func (p *Monitor) OnClose() {}
func (p *Monitor) OnInit() {
	p.httpMux = http.NewServeMux()
	p.httpServer = &http.Server{
		Handler:      p.httpMux,
		Addr:         fmt.Sprintf(":%d", p.httpPort),
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
	}
	p.httpServer.SetKeepAlivesEnabled(false)
}

func (p *Monitor) Run(closeChan chan struct{}) {
	p.httpMux.HandleFunc("/monitor", monitor)
	p.httpMux.HandleFunc("/debug/pprof/", pprof.Index)

	for pattern, handler := range insMonitorHandlerMap {
		p.httpMux.HandleFunc(pattern, handler)
	}
	go func() {
		p.httpServer.ListenAndServe()
	}()

	<-closeChan
}

func monitor(w http.ResponseWriter, r *http.Request) {
	var memstat runtime.MemStats
	runtime.ReadMemStats(&memstat)
	s := fmt.Sprintf("time\t%d\n", time.Now().Unix())
	s += fmt.Sprintf("goroutines\t%d\n", runtime.NumGoroutine())
	s += fmt.Sprintf("memAlloc(MB)\t%.2f\n", float64(memstat.Alloc)/1048576)
	s += fmt.Sprintf("heapAlloc(MB)\t%.2f\n", float64(memstat.HeapAlloc)/1048576)
	s += fmt.Sprintf("heapSys(MB)\t%.2f\n", float64(memstat.HeapSys)/1048576)
	s += fmt.Sprintf("heapIdle(MB)\t%.2f\n", float64(memstat.HeapIdle)/1048576)
	fmt.Fprintf(w, s)
}

func (p *Monitor) Name() string { return "monitor" }

type MonitorHandlerMap map[string]func(http.ResponseWriter, *http.Request)

var insMonitorHandlerMap MonitorHandlerMap = make(MonitorHandlerMap)

func Add(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	insMonitorHandlerMap[pattern] = handler
}

func Init(port int) {
	Service.httpPort = port
}

var Service = new(Monitor)
