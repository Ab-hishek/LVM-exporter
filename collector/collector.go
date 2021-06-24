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
type LvmCollector struct {
	vgSizeMetric              *prometheus.Desc
	vgFreeMetric              *prometheus.Desc
	vgLvCountMetric           *prometheus.Desc
	vgPvCountMetric           *prometheus.Desc
	vgMaxLvMetric             *prometheus.Desc
	vgMaxPvMetric             *prometheus.Desc
	vgSnapCountMetric         *prometheus.Desc
	vgMissingPvCountMetric    *prometheus.Desc
	vgMetadataCountMetric     *prometheus.Desc
	vgMetadataUsedCountMetric *prometheus.Desc
	vgMetadataFreeMetric      *prometheus.Desc
	vgMetadataSizeMetric      *prometheus.Desc
	vgMetadataCopiesMetric    *prometheus.Desc
	vgPermissionsMetric       *prometheus.Desc
	vgAllocationPolicyMetric  *prometheus.Desc

	lvSizeMetric *prometheus.Desc
}

// You must create a constructor for your collector that
// initializes every descriptor and returns a pointer to the collector
func NewLvmCollector() *LvmCollector {
	return &LvmCollector{
		vgFreeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "free_bytes"),
			"LVM VG free size in bytes",
			[]string{"name"}, nil,
		),
		vgSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "total_bytes"),
			"LVM VG total size in bytes",
			[]string{"name"}, nil,
		),
		vgLvCountMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "lv_count"),
			"Number of LVs in VG",
			[]string{"name"}, nil,
		),
		vgPvCountMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "pv_count"),
			"Number of PVs in VG",
			[]string{"name"}, nil,
		),
		vgMaxLvMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "max_lv_count"),
			"LMaximum number of LVs allowed in VG or 0 if unlimited",
			[]string{"name"}, nil,
		),
		vgMaxPvMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "max_pv_count"),
			"Maximum number of PVs allowed in VG or 0 if unlimited",
			[]string{"name"}, nil,
		),
		vgSnapCountMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "snap_count"),
			"Number of snapshots in VG",
			[]string{"name"}, nil,
		),
		vgMissingPvCountMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "missing_pv_count"),
			"Number of PVs in VG which are missing",
			[]string{"name"}, nil,
		),
		vgMetadataCountMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "mda_count"),
			"Number of metadata areas on this VG",
			[]string{"name"}, nil,
		),
		vgMetadataUsedCountMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "mda_used_count"),
			"Number of metadata areas in use on this VG",
			[]string{"name"}, nil,
		),
		vgMetadataFreeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "mda_free_bytes"),
			"Free metadata area space for this VG in bytes",
			[]string{"name"}, nil,
		),
		vgMetadataSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "mda_total_size_bytes"),
			"Size of smallest metadata area for this VG in bytes",
			[]string{"name"}, nil,
		),
		vgPermissionsMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "permission"),
			"VG permissions: [-1: undefined], [0: writeable], [1: read-only]",
			[]string{"name"}, nil,
		),
		vgAllocationPolicyMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "allocation_policy"),
			"VG allocation policy: [-1: undefined], [0: normal], [1: contiguous], [2: cling], [3: anywhere], [4: inherited]",
			[]string{"name"}, nil,
		),

		lvSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "total_size_bytes"),
			"LVM LV total size in bytes",
			[]string{"name", "path", "dm_path", "vg", "device"}, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *LvmCollector) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.vgSizeMetric
	ch <- collector.vgFreeMetric
	ch <- collector.vgLvCountMetric
	ch <- collector.vgPvCountMetric
	ch <- collector.vgMaxLvMetric
	ch <- collector.vgMaxPvMetric
	ch <- collector.vgSnapCountMetric
	ch <- collector.vgMissingPvCountMetric
	ch <- collector.vgMetadataCountMetric
	ch <- collector.vgMetadataUsedCountMetric
	ch <- collector.vgMetadataFreeMetric
	ch <- collector.vgMetadataSizeMetric
	ch <- collector.vgPermissionsMetric
	ch <- collector.vgAllocationPolicyMetric

	ch <- collector.lvSizeMetric
}

// Collect implements required collect function for all prometheus collectors
func (collector *LvmCollector) Collect(ch chan<- prometheus.Metric) {
	vgList, err := lvm.ListLVMVolumeGroup()
	if err != nil {
		klog.Errorf("error in getting the list of lvm volume groups: %v", err)
	} else {
		for _, vg := range vgList {
			ch <- prometheus.MustNewConstMetric(collector.vgFreeMetric, prometheus.GaugeValue, vg.Free.AsApproximateFloat64(), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgSizeMetric, prometheus.GaugeValue, vg.Size.AsApproximateFloat64(), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgLvCountMetric, prometheus.CounterValue, float64(vg.LVCount), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgPvCountMetric, prometheus.CounterValue, float64(vg.PVCount), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgMaxLvMetric, prometheus.CounterValue, float64(vg.MaxLV), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgMaxPvMetric, prometheus.CounterValue, float64(vg.MaxPV), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgSnapCountMetric, prometheus.CounterValue, float64(vg.SnapCount), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgMissingPvCountMetric, prometheus.CounterValue, float64(vg.MissingPVCount), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgMetadataCountMetric, prometheus.CounterValue, float64(vg.MetadataCount), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgMetadataUsedCountMetric, prometheus.CounterValue, float64(vg.MetadataUsedCount), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgMetadataFreeMetric, prometheus.GaugeValue, vg.MetadataFree.AsApproximateFloat64(), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgMetadataSizeMetric, prometheus.GaugeValue, vg.MetadataSize.AsApproximateFloat64(), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgPermissionsMetric, prometheus.CounterValue, float64(vg.Permission), vg.Name)
			ch <- prometheus.MustNewConstMetric(collector.vgAllocationPolicyMetric, prometheus.CounterValue, float64(vg.AllocationPolicy), vg.Name)
		}
	}
}
