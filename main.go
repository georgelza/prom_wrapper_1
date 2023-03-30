/*****************************************************************************
*
*	File			: main.go
*
* 	Created			: 27 March 2023
*
*	Description		: Quick Dirty wrapper for Prometheus and golang library to figure out how to back port it into fs_loader
*
*	Modified		: 27 March 2023	- Start
*
*	By			: George Leonard (georgelza@gmail.com)
*
*
*
*****************************************************************************/

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var sql_duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "fs_sql_duration_seconds",
	Help:    "Histogram for the duration in seconds.",
	Buckets: []float64{1, 2, 5, 6, 10},
},
	[]string{"endpoint"},
)

var api_duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "fs_api_duration_seconds",
	Help:    "Histogram for the duration in seconds.",
	Buckets: []float64{1, 2, 5, 6, 10},
},
	[]string{"endpoint"},
)

var requestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "fs_etl_operations_count",
	Help: "The total number of processed records",
})

var info = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "txn_count",
	Help: "Information about the batch size",
},
	[]string{"batch"},
)

func loadEFT() {

	var x int = 2
	var txn_count float64

	txn_count = 9752395 // this will be the recordcount of the records returned by the sql query
	info.With(prometheus.Labels{"batch": "eft"}).Set(txn_count)

	txn_count = 104565
	info.With(prometheus.Labels{"batch": "acc"}).Set(txn_count)

	////////////////////////////////
	// start a timer
	sTime := time.Now()

	// Execute a large sql #1 execute
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10000) // if vGeneral.sleep = 10000, 10 second
	fmt.Printf("EFT SQL Sleeping %d Millisecond...\n", n)
	time.Sleep(time.Duration(n) * time.Millisecond)

	// post to Prometheus
	sqlDuration := time.Since(sTime)
	sql_duration.WithLabelValues("eft_sql_seconds").Observe(sqlDuration.Seconds())

	////////////////////////////////
	// start timer
	sTime = time.Now()

	// Execute a large sql #2 execute
	rand.Seed(time.Now().UnixNano())
	n = rand.Intn(10000) // if vGeneral.sleep = 10000, 10 second
	fmt.Printf("ACC SQL Sleeping %d Millisecond...\n", n)
	time.Sleep(time.Duration(n) * time.Millisecond)

	// post to Prometheus
	sqlDuration = time.Since(sTime)
	sql_duration.WithLabelValues("acc_sql_seconds").Observe(sqlDuration.Seconds())

	for count := 0; count < x; count++ {

		// EFT
		sTime = time.Now()
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(5000) // if vGeneral.sleep = 1000, then n will be random value of 0 -> 1000  aka 0 and 1 second
		fmt.Printf("Sleeping %d Millisecond...\n", n)
		time.Sleep(time.Duration(n) * time.Millisecond)

		//determine the duration and log to prometheus
		fsDuration := time.Since(sTime)
		api_duration.WithLabelValues("payment_event_nrt_eft_seconds").Observe(fsDuration.Seconds())

		// ACC
		sTime = time.Now()
		rand.Seed(time.Now().UnixNano())
		n = rand.Intn(5000) // if vGeneral.sleep = 1000, then n will be random value of 0 -> 1000  aka 0 and 1 second
		fmt.Printf("Sleeping %d Millisecond...\n", n)
		time.Sleep(time.Duration(n) * time.Millisecond)

		//determine the duration and log to prometheus
		fsDuration = time.Since(sTime)
		api_duration.WithLabelValues("payment_event_nrt_acc_seconds").Observe(fsDuration.Seconds())

		//increment a counter for number of requests processed
		requestsProcessed.Inc()

		println(count)
	}
	os.Exit(0)
}

func main() {

	fmt.Println("starting...")

	prometheus.MustRegister(sql_duration, api_duration, info)

	go loadEFT()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9000", nil)

}
