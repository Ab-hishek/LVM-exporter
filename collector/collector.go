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

	lvSizeMetric                *prometheus.Desc
	lvUsedSizePercentMetric     *prometheus.Desc
	lvPermissionMetric          *prometheus.Desc
	lvBehaviourWhenFullMetric   *prometheus.Desc
	lvHealthStatusMetric        *prometheus.Desc
	lvRaidSyncActionMetric      *prometheus.Desc
	lvMetadataSizeMetric        *prometheus.Desc
	lvMetadataUsedPercentMetric *prometheus.Desc
	lvSnapshotUsedPercentMetric *prometheus.Desc

	pvSizeMetric         *prometheus.Desc
	pvFreeMetric         *prometheus.Desc
	pvUsedMetric         *prometheus.Desc
	pvDeviceSizeMetric   *prometheus.Desc
	pvMetadataSizeMetric *prometheus.Desc
	pvMetadataFreeMetric *prometheus.Desc
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
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvUsedSizePercentMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "used_size_percent"),
			"LVM LV used size in percentage",
			[]string{"name", "path", "dm_path", "vg", "device", "host", "segtype", "pool", "active_status"}, nil,
		),
		lvPermissionMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "permission"),
			"VG permissions: [-1: undefined], [0: writeable], [1: read-only], [2: read-only-override]",
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
		lvMetadataSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "mda_size_bytes"),
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

		pvSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "total_size_bytes"),
			"LVM PV total size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in-use"}, nil,
		),
		pvFreeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "free_bytes"),
			"LVM PV free size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in-use"}, nil,
		),
		pvUsedMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "used_bytes"),
			"LVM PV used size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in-use"}, nil,
		),
		pvDeviceSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "device_size_bytes"),
			"LVM PV underlying device size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in-use"}, nil,
		),
		pvMetadataSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "mda_size_bytes"),
			"LVM PV device smallest metadata area size in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in-use"}, nil,
		),
		pvMetadataFreeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "pv", "mda_free_bytes"),
			"LVM PV device free metadata area space in bytes",
			[]string{"name", "allocatable", "vg", "missing", "in-use"}, nil,
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
	ch <- collector.lvUsedSizePercentMetric
	ch <- collector.lvPermissionMetric
	ch <- collector.lvBehaviourWhenFullMetric
	ch <- collector.lvHealthStatusMetric
	ch <- collector.lvRaidSyncActionMetric
	ch <- collector.lvMetadataSizeMetric
	ch <- collector.lvMetadataUsedPercentMetric
	ch <- collector.lvSnapshotUsedPercentMetric

	ch <- collector.pvSizeMetric
	ch <- collector.pvFreeMetric
	ch <- collector.pvUsedMetric
	ch <- collector.pvDeviceSizeMetric
	ch <- collector.pvMetadataSizeMetric
	ch <- collector.pvMetadataFreeMetric

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

	lvList, err := lvm.ListLVMLogicalVolume()
	if err != nil {
		klog.Errorf("error in getting the list of lvm logical volumes: %v", err)
	} else {
		for _, lv := range lvList {
			ch <- prometheus.MustNewConstMetric(collector.lvSizeMetric, prometheus.GaugeValue, lv.Size.AsApproximateFloat64(), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvUsedSizePercentMetric, prometheus.GaugeValue, lv.UsedSizePercent, lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvPermissionMetric, prometheus.CounterValue, float64(lv.Permission), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvBehaviourWhenFullMetric, prometheus.CounterValue, float64(lv.BehaviourWhenFull), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvHealthStatusMetric, prometheus.CounterValue, float64(lv.HealthStatus), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvRaidSyncActionMetric, prometheus.CounterValue, float64(lv.RaidSyncAction), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvMetadataSizeMetric, prometheus.GaugeValue, lv.MetadataSize.AsApproximateFloat64(), lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvMetadataUsedPercentMetric, prometheus.GaugeValue, lv.MetadataUsedPercent, lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
			ch <- prometheus.MustNewConstMetric(collector.lvSnapshotUsedPercentMetric, prometheus.GaugeValue, lv.SnapshotUsedPercent, lv.Name, lv.Path, lv.DMPath, lv.VGName, lv.Device, lv.Host, lv.SegType, lv.ActiveStatus)
		}
	}

	pvList, err := lvm.ListLVMPhysicalVolume()
	if err != nil {
		klog.Errorf("error in getting the list of lvm logical volumes: %v", err)
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
