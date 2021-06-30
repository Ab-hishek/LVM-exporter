package collector

import (
	"github.com/Ab-hishek/LVM-exporter/lvm"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog"
)

// Define a struct for you collector that contains pointers
// to prometheus descriptors for each metric you wish to expose.
// Note you can also include fields of other types if they provide utility
// but we just won't be exposing them as metrics.
type pvCollector struct {
	pvSizeMetric         *prometheus.Desc
	pvFreeMetric         *prometheus.Desc
	pvUsedMetric         *prometheus.Desc
	pvDeviceSizeMetric   *prometheus.Desc
	pvMetadataSizeMetric *prometheus.Desc
	pvMetadataFreeMetric *prometheus.Desc
}

// You must create a constructor for your collector that
// initializes every descriptor and returns a pointer to the collector
func NewPvCollector() *pvCollector {
	return &pvCollector{
		pvSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "total_size_bytes"),
			"LVM PV total size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in_use"}, nil,
		),
		pvFreeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "free_size_bytes"),
			"LVM PV free size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in_use"}, nil,
		),
		pvUsedMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "used_size_bytes"),
			"LVM PV used size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in_use"}, nil,
		),
		pvDeviceSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "device_size_bytes"),
			"LVM PV underlying device size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in_use"}, nil,
		),
		pvMetadataSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "mda_total_size_bytes"),
			"LVM PV device smallest metadata area size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in_use"}, nil,
		),
		pvMetadataFreeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "mda_free_size_bytes"),
			"LVM PV device free metadata area space in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in_use"}, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *pvCollector) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.pvSizeMetric
	ch <- collector.pvFreeMetric
	ch <- collector.pvUsedMetric
	ch <- collector.pvDeviceSizeMetric
	ch <- collector.pvMetadataSizeMetric
	ch <- collector.pvMetadataFreeMetric
}

// Collect implements required collect function for all prometheus collectors
func (collector *pvCollector) Collect(ch chan<- prometheus.Metric) {
	pvList, err := lvm.ListLVMPhysicalVolume()
	if err != nil {
		klog.Errorf("error in getting the list of lvm physical volumes: %v", err)
	} else {
		for _, pv := range pvList {
			ch <- prometheus.MustNewConstMetric(collector.pvSizeMetric, prometheus.GaugeValue, pv.Size.AsApproximateFloat64(), pv.Name, pv.Allocatable, pv.VGName, pv.Missing, pv.InUse)
			ch <- prometheus.MustNewConstMetric(collector.pvFreeMetric, prometheus.GaugeValue, pv.Free.AsApproximateFloat64(), pv.Name, pv.Allocatable, pv.VGName, pv.Missing, pv.InUse)
			ch <- prometheus.MustNewConstMetric(collector.pvUsedMetric, prometheus.GaugeValue, pv.Used.AsApproximateFloat64(), pv.Name, pv.Allocatable, pv.VGName, pv.Missing, pv.InUse)
			ch <- prometheus.MustNewConstMetric(collector.pvDeviceSizeMetric, prometheus.GaugeValue, pv.DeviceSize.AsApproximateFloat64(), pv.Name, pv.Allocatable, pv.VGName, pv.Missing, pv.InUse)
			ch <- prometheus.MustNewConstMetric(collector.pvMetadataSizeMetric, prometheus.GaugeValue, pv.MetadataSize.AsApproximateFloat64(), pv.Name, pv.Allocatable, pv.VGName, pv.Missing, pv.InUse)
			ch <- prometheus.MustNewConstMetric(collector.pvMetadataFreeMetric, prometheus.GaugeValue, pv.MetadataFree.AsApproximateFloat64(), pv.Name, pv.Allocatable, pv.VGName, pv.Missing, pv.InUse)
		}
	}
}
