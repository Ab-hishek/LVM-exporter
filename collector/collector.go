package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"LVM-exporter/lvm"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type LvmCollector struct {
	vgSizeMetric *prometheus.Desc
	vgFreeMetric *prometheus.Desc
}

//You must create a constructor for your collector that
//initializes every descriptor and returns a pointer to the collector
func NewLvmCollector() *LvmCollector {
	return &LvmCollector{
		vgSizeMetric: prometheus.NewDesc("vg_size_metric",
			"Shows the total size of the LVM volume group in gb",
			[]string{"vg_name"}, nil,
		),
		vgFreeMetric: prometheus.NewDesc("lv_size_metric",
			"Shows the free size of the LVM volume group in gb",
			[]string{"vg_name"}, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *LvmCollector) Describe(ch chan<- *prometheus.Desc) {
	//vg := lvm.VolumeGroup{}
	//Update this section with the each metric you create for a given collector
	ch <- collector.vgSizeMetric
	ch <- collector.vgFreeMetric
}

//Collect implements required collect function for all prometheus collectors
func (collector *LvmCollector) Collect(ch chan<- prometheus.Metric) {

	out, err := exec.Command("/sbin/vgs", "--units", "g", "--separator", ",", "-o", "vg_name,vg_free,vg_size", "--noheadings").Output()
	if err != nil {
		log.Print(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		values := strings.Split(line, ",")
		if len(values) == 3 {
			free_size, err := strconv.ParseFloat(strings.Trim(values[1], "g"), 64)
			if err != nil {
				log.Print(err)
			} else {
				total_size, err := strconv.ParseFloat(strings.Trim(values[2], "g"), 64)
				if err != nil {
					log.Print(err)
				} else {
					vg_name := strings.Trim(values[0], " ")
					ch <- prometheus.MustNewConstMetric(collector.vgFreeMetric, prometheus.GaugeValue, free_size, vg_name)
					ch <- prometheus.MustNewConstMetric(collector.vgSizeMetric, prometheus.GaugeValue, total_size, vg_name)
				}
			}
		}
	}
}
