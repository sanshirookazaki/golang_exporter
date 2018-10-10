package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fukata/golang-stats-api-handler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port        = flag.String("port", ":9696", "The address to listen on for HTTP requests.")
	apiendpoint = flag.String("apiendpoint", "http://localhost:9696/api/stats", "api endpoint")
)

const (
	namespace = "golang"
)

type Gometrics struct {
	Time             int64         `json:"time"`
	GoVersion        string        `json:"go_version"`
	GoOs             string        `json:"go_os"`
	GoArch           string        `json:"go_arch"`
	CPUNum           int           `json:"cpu_num"`
	GoroutineNum     int           `json:"goroutine_num"`
	Gomaxprocs       int           `json:"gomaxprocs"`
	CgoCallNum       int           `json:"cgo_call_num"`
	MemoryAlloc      int           `json:"memory_alloc"`
	MemoryTotalAlloc int           `json:"memory_total_alloc"`
	MemorySys        int           `json:"memory_sys"`
	MemoryLookups    int           `json:"memory_lookups"`
	MemoryMallocs    int           `json:"memory_mallocs"`
	MemoryFrees      int           `json:"memory_frees"`
	MemoryStack      int           `json:"memory_stack"`
	HeapAlloc        int           `json:"heap_alloc"`
	HeapSys          int           `json:"heap_sys"`
	HeapIdle         int           `json:"heap_idle"`
	HeapInuse        int           `json:"heap_inuse"`
	HeapReleased     int           `json:"heap_released"`
	HeapObjects      int           `json:"heap_objects"`
	GcNext           int           `json:"gc_next"`
	GcLast           int           `json:"gc_last"`
	GcNum            int           `json:"gc_num"`
	GcPerSecond      int           `json:"gc_per_second"`
	GcPausePerSecond int           `json:"gc_pause_per_second"`
	GcPause          []interface{} `json:"gc_pause"`
}

type goCollector struct {
	Time             prometheus.Gauge
	GoVersion        prometheus.Gauge
	GoOs             prometheus.Gauge
	GoArch           prometheus.Gauge
	CPUNum           prometheus.Gauge
	GoroutineNum     prometheus.Gauge
	Gomaxprocs       prometheus.Gauge
	CgoCallNum       prometheus.Gauge
	MemoryAlloc      prometheus.Gauge
	MemoryTotalAlloc prometheus.Gauge
	MemorySys        prometheus.Gauge
	MemoryLookups    prometheus.Gauge
	MemoryMallocs    prometheus.Gauge
	MemoryFrees      prometheus.Gauge
	MemoryStack      prometheus.Gauge
	HeapAlloc        prometheus.Gauge
	HeapSys          prometheus.Gauge
	HeapIdle         prometheus.Gauge
	HeapInuse        prometheus.Gauge
	HeapReleased     prometheus.Gauge
	HeapObjects      prometheus.Gauge
	GcNext           prometheus.Gauge
	GcLast           prometheus.Gauge
	GcNum            prometheus.Gauge
	GcPerSecond      prometheus.Gauge
	GcPausePerSecond prometheus.Gauge
}

func newgoCollector() *goCollector {
	return &goCollector{
		Time: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up_time",
			Help:      "golang up time",
		}),
		GoVersion: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "version",
			Help:      "golang version",
		}),
		GoOs: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "os",
			Help:      "golang os",
		}),
		GoArch: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "arch",
			Help:      "golang arch",
		}),
		CPUNum: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cpu_num",
			Help:      "cpu num",
		}),
		GoroutineNum: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "goroutine_num",
			Help:      "goroutine num",
		}),
		Gomaxprocs: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gomaxproc",
			Help:      "golang max proc",
		}),
		CgoCallNum: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cgo_call_num",
			Help:      "golang cgo call num",
		}),
		MemoryAlloc: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_alloc",
			Help:      "golang memory alloc",
		}),
		MemoryTotalAlloc: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_total_alloc",
			Help:      "golang memory total alloc",
		}),
		MemorySys: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_sys",
			Help:      "golang memory sys",
		}),
		MemoryLookups: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_lookups",
			Help:      "golang memory lookups",
		}),
		MemoryMallocs: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_mallocs",
			Help:      "golang memory mallocs",
		}),
		MemoryFrees: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_free",
			Help:      "golang memory free",
		}),
		MemoryStack: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_stack",
			Help:      "golang memory stack",
		}),
		HeapAlloc: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heap_alloc",
			Help:      "golang heap alloc",
		}),
		HeapSys: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heap_sys",
			Help:      "golang heap sys",
		}),
		HeapIdle: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heap_idle",
			Help:      "golang heap idle",
		}),
		HeapInuse: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heap_inuse",
			Help:      "golang heap inuse",
		}),
		HeapReleased: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heap_released",
			Help:      "golang heap released",
		}),
		HeapObjects: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "heap_objects",
			Help:      "golang heap objects",
		}),
		GcNext: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gc_next",
			Help:      "golang gc next",
		}),
		GcLast: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gc_last",
			Help:      "golang gc last",
		}),
		GcNum: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gc_num",
			Help:      "golang gc num",
		}),
		GcPerSecond: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gc_per_second",
			Help:      "golang gc per second",
		}),
		GcPausePerSecond: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gc_pause_per_second",
			Help:      "golang gc pause per second",
		}),
	}
}

func (c *goCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Time.Desc()
	ch <- c.GoVersion.Desc()
	ch <- c.GoOs.Desc()
	ch <- c.GoArch.Desc()
	ch <- c.CPUNum.Desc()
	ch <- c.GoroutineNum.Desc()
	ch <- c.Gomaxprocs.Desc()
	ch <- c.CgoCallNum.Desc()
	ch <- c.MemoryAlloc.Desc()
	ch <- c.MemoryTotalAlloc.Desc()
	ch <- c.MemorySys.Desc()
	ch <- c.MemoryLookups.Desc()
	ch <- c.MemoryMallocs.Desc()
	ch <- c.MemoryFrees.Desc()
	ch <- c.MemoryStack.Desc()
	ch <- c.HeapAlloc.Desc()
	ch <- c.HeapSys.Desc()
	ch <- c.HeapIdle.Desc()
	ch <- c.HeapInuse.Desc()
	ch <- c.HeapReleased.Desc()
	ch <- c.HeapObjects.Desc()
	ch <- c.GcNext.Desc()
	ch <- c.GcLast.Desc()
	ch <- c.GcNum.Desc()
	ch <- c.GcPerSecond.Desc()
	ch <- c.GcPausePerSecond.Desc()
}

func (c *goCollector) Collect(ch chan<- prometheus.Metric) {
	//dummyStaticNumber := float64(1234)
	res, err := http.Get(*apiendpoint)

	if err != nil {
		return //err
	} else if res.StatusCode != 200 {
		return //fmt.Errorf("Unable to get this url : http status %d", res.StatusCode)
	}
	defer res.Body.Close()
	// jsonを読み込む
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return //nil
	}
	jsonBytes := ([]byte)(body)
	data := new(Gometrics)
	if err := json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal ERROR!", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.Time.Desc(),
		prometheus.GaugeValue,
		float64(data.Time),
	)
	ch <- prometheus.MustNewConstMetric(
		c.CPUNum.Desc(),
		prometheus.GaugeValue,
		float64(data.CPUNum),
	)
	ch <- prometheus.MustNewConstMetric(
		c.GoroutineNum.Desc(),
		prometheus.GaugeValue,
		float64(data.GoroutineNum),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Gomaxprocs.Desc(),
		prometheus.GaugeValue,
		float64(data.Gomaxprocs),
	)
	ch <- prometheus.MustNewConstMetric(
		c.CgoCallNum.Desc(),
		prometheus.GaugeValue,
		float64(data.CgoCallNum),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MemoryAlloc.Desc(),
		prometheus.GaugeValue,
		float64(data.MemoryAlloc),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MemoryTotalAlloc.Desc(),
		prometheus.GaugeValue,
		float64(data.MemoryTotalAlloc),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MemorySys.Desc(),
		prometheus.GaugeValue,
		float64(data.MemorySys),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MemoryLookups.Desc(),
		prometheus.GaugeValue,
		float64(data.MemoryLookups),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MemoryMallocs.Desc(),
		prometheus.GaugeValue,
		float64(data.MemoryMallocs),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MemoryFrees.Desc(),
		prometheus.GaugeValue,
		float64(data.MemoryFrees),
	)
	ch <- prometheus.MustNewConstMetric(
		c.MemoryStack.Desc(),
		prometheus.GaugeValue,
		float64(data.MemoryStack),
	)
	ch <- prometheus.MustNewConstMetric(
		c.HeapAlloc.Desc(),
		prometheus.GaugeValue,
		float64(data.HeapAlloc),
	)
	ch <- prometheus.MustNewConstMetric(
		c.HeapSys.Desc(),
		prometheus.GaugeValue,
		float64(data.HeapSys),
	)
	ch <- prometheus.MustNewConstMetric(
		c.HeapIdle.Desc(),
		prometheus.GaugeValue,
		float64(data.HeapIdle),
	)
	ch <- prometheus.MustNewConstMetric(
		c.HeapInuse.Desc(),
		prometheus.GaugeValue,
		float64(data.HeapReleased),
	)
	ch <- prometheus.MustNewConstMetric(
		c.HeapObjects.Desc(),
		prometheus.GaugeValue,
		float64(data.HeapObjects),
	)
	ch <- prometheus.MustNewConstMetric(
		c.GcNext.Desc(),
		prometheus.GaugeValue,
		float64(data.GcNext),
	)
	ch <- prometheus.MustNewConstMetric(
		c.GcLast.Desc(),
		prometheus.GaugeValue,
		float64(data.GcLast),
	)
	ch <- prometheus.MustNewConstMetric(
		c.GcNum.Desc(),
		prometheus.GaugeValue,
		float64(data.GcNum),
	)
	ch <- prometheus.MustNewConstMetric(
		c.GcPerSecond.Desc(),
		prometheus.GaugeValue,
		float64(data.GcPerSecond),
	)
	ch <- prometheus.MustNewConstMetric(
		c.GcPausePerSecond.Desc(),
		prometheus.GaugeValue,
		float64(data.GcPausePerSecond),
	)
}

func main() {

	flag.Parse()

	c := newgoCollector()
	prometheus.MustRegister(c)

	fmt.Println("access to", *port+"/metrics")
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/api/stats", stats_api.Handler)
	log.Fatal(http.ListenAndServe(*port, nil))

}
