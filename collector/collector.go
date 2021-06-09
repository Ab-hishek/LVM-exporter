package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type LvmCollector struct {
	vgSizeMetric *prometheus.Desc
	lvSizeMetric *prometheus.Desc
}

//You must create a constructor for your collector that
//initializes every descriptor and returns a pointer to the collector
func NewLvmCollector() *LvmCollector {
	return &LvmCollector{
		vgSizeMetric: prometheus.NewDesc("vg_size_metric",
			"Shows the size of the volume group present in the node",
			nil, nil,
		),
		lvSizeMetric: prometheus.NewDesc("lv_size_metric",
			"Shows size of the logical volume present in the node",
			nil, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *LvmCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.vgSizeMetric
	ch <- collector.lvSizeMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *LvmCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	if 1 == 1 {
		metricValue = 1
	}

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(collector.vgSizeMetric, prometheus.CounterValue, metricValue)
	ch <- prometheus.MustNewConstMetric(collector.lvSizeMetric, prometheus.CounterValue, metricValue)

}