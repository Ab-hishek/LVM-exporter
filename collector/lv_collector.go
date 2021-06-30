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
type lvCollector struct {
	lvSizeMetric                *prometheus.Desc
	lvUsedSizePercentMetric     *prometheus.Desc
	lvPermissionMetric          *prometheus.Desc
	lvBehaviourWhenFullMetric   *prometheus.Desc
	lvHealthStatusMetric        *prometheus.Desc
	lvRaidSyncActionMetric      *prometheus.Desc
	lvMetadataSizeMetric        *prometheus.Desc
	lvMetadataUsedPercentMetric *prometheus.Desc
	lvSnapshotUsedPercentMetric *prometheus.Desc
}

// You must create a constructor for your collector that
// initializes every descriptor and returns a pointer to the collector
func NewLvCollector() *lvCollector {
	return &lvCollector{
		lvSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "total_size_bytes"),
			"LVM LV total size in bytes",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvUsedSizePercentMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "used_percent"),
			"LVM LV used size in percentage",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvPermissionMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "permission"),
			"VG permissions: [-1: undefined], [0: unknown], [1: writeable], [2: read-only], [3: read-only-override]",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvBehaviourWhenFullMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "when_full"),
			"For thin pools, behavior when full: [-1: undefined], [0: error], [1: queue]",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvHealthStatusMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "health_status"),
			"LV health status: [-1: undefined], [0: \"\"], [1: partial], [2: refresh needed], [3: mismatches exist]",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvRaidSyncActionMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "raid_sync_action"),
			"For LV RAID, the current synchronization action being performed: [-1: undefined], [0: idle], [1: frozen], [2: resync], [3: recover], [4: check], [5: repair]",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvMetadataSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "mda_total_size_bytes"),
			"LVM LV metadata size in bytes",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvMetadataUsedPercentMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "mda_used_percent"),
			"LVM LV metadata used size in percentage",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvSnapshotUsedPercentMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "snap_percent"),
			"LVM LV snap used size in percentage",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *lvCollector) Describe(ch chan<- *prometheus.Desc) {
	//Update this section with the each metric you create for a given collector
	ch <- collector.lvSizeMetric
	ch <- collector.lvUsedSizePercentMetric
	ch <- collector.lvPermissionMetric
	ch <- collector.lvBehaviourWhenFullMetric
	ch <- collector.lvHealthStatusMetric
	ch <- collector.lvRaidSyncActionMetric
	ch <- collector.lvMetadataSizeMetric
	ch <- collector.lvMetadataUsedPercentMetric
	ch <- collector.lvSnapshotUsedPercentMetric
}

// Collect implements required collect function for all prometheus collectors
func (collector *lvCollector) Collect(ch chan<- prometheus.Metric) {
	lvList, err := lvm.ListLVMLogicalVolume()
	if err != nil {
		klog.Errorf("error in getting the list of lvm logical volumes: %v", err)
	} else {
		for _, lv := range lvList {
			ch <- prometheus.MustNewConstMetric(collector.lvSizeMetric, prometheus.GaugeValue, lv.Size.AsApproximateFloat64(), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvUsedSizePercentMetric, prometheus.GaugeValue, lv.UsedSizePercent, lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvPermissionMetric, prometheus.GaugeValue, float64(lv.Permission), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvBehaviourWhenFullMetric, prometheus.GaugeValue, float64(lv.BehaviourWhenFull), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvHealthStatusMetric, prometheus.GaugeValue, float64(lv.HealthStatus), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvRaidSyncActionMetric, prometheus.GaugeValue, float64(lv.RaidSyncAction), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvMetadataSizeMetric, prometheus.GaugeValue, lv.MetadataSize.AsApproximateFloat64(), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvMetadataUsedPercentMetric, prometheus.GaugeValue, lv.MetadataUsedPercent, lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvSnapshotUsedPercentMetric, prometheus.GaugeValue, lv.SnapshotUsedPercent, lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.PoolName, lv.ActiveStatus)
		}
	}
}
