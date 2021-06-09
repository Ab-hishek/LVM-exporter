package main

import (
	"LVM-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	//"github.com/joho/godotenv"
	//"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {

	//Create a new instance of the foocollector and
	//register it with the prometheus client.
	lvm := collector.NewLvmCollector()
	prometheus.MustRegister(lvm)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}