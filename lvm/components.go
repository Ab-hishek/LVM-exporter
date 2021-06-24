package lvm

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

// VolumeGroup specifies attributes of a given vg exists on node.
type VolumeGroup struct {
	// Name of the lvm volume group.
	Name string `json:"vg_name"`

	// UUID denotes a unique identity of a lvm volume group.
	UUID string `json:"vg_uuid"`

	// Size specifies the total size of volume group in bytes.
	Size resource.Quantity `json:"vg_size"`

	// Free specifies the available capacity of volume group in bytes.
	Free resource.Quantity `json:"vg_free"`

	// LVCount denotes total number of logical volumes in
	// volume group.
	LVCount int32 `json:"lv_count"`

	// PVCount denotes total number of physical volumes
	// constituting the volume group.
	PVCount int32 `json:"pv_count"`

	MaxLV int32 `json:"max_lv"`

	MaxPV int32 `json:"max_pv"`

	SnapCount int32 `json:"snap_count"`

	MissingPVCount int32 `json:"vg_missing_pv_count"`

	MetadataCount int32 `json:"vg_mda_count"`

	MetadataUsedCount int32 `json:"vg_mda_used_count"`

	MetadataFree resource.Quantity `json:"vg_mda_free"`

	MetadataSize resource.Quantity `json:"vg_mda_size"`

	Permission int `json:"vg_permissions"`

	AllocationPolicy int `json:"vg_allocation_policy"`
}

// LogicalVolume specifies attributes of a given lv that exists on the node.
type LogicalVolume struct {
	// Name of the lvm logical volume.
	Name string `json:"lv_name"`

	// Full name of the lvm logical volume.
	FullName string `json:"lv_full_name"`

	// UUID denotes a unique identity of a lvm logical volume.
	UUID string `json:"lv_uuid"`

	// Size specifies the total size of logical volume in bytes
	Size resource.Quantity `json:"lv_size"`

	// Path specifies LVM logical volume path
	Path string `json:"lv_path"`

	// DMPath specifies device mapper path
	DMPath string `json:"lv_dm_path"`

	// LVM logical volume device
	Device string

	// Name of the VG in which LVM logical volume is created
	VGName string `json:"vg_name"`

	SegType string `json:"segtype"`

	Permission int `json:"lv_permissions"`

	BehaviourWhenFull int `json:"lv_when_full"`

	HealthStatus int `json:"lv_health_status"`

	RaidSyncAction int `json:"raid_sync_action"`

	ActiveStatus string `json:"lv_active"`

	Host string `json:"lv_host"`

	PoolName string `json:"pool_lv"`

	UsedSizePercent float64 `json:"data_percent"`

	MetadataSize resource.Quantity `json:"lv_metadata_size"`

	MetadataUsedPercent float64 `json:"metadata_percent"`

	SnapshotUsedPercent float64 `json:"snap_percent"`
}

// PhysicalVolume specifies attributes of a given pv that exists on the node.
type PhysicalVolume struct {
	// Name of the lvm physical volume.
	Name string `json:"pv_name"`

	// UUID denotes a unique identity of a lvm physical volume.
	UUID string `json:"pv_uuid"`

	// Size specifies the total size of physical volume in bytes
	Size resource.Quantity `json:"pv_size"`

	DeviceSize resource.Quantity `json:"dev_size"`

	MetadataSize resource.Quantity `json:"pv_mda_size"`

	MetadataFree resource.Quantity `json:"pv_mda_free"`

	Free resource.Quantity `json:"pv_free"`

	Used resource.Quantity `json:"pv_used"`

	Allocatable string `json:"pv_allocatable"`

	Missing string `json:"pv_missing"`

	InUse string `json:"pv_in_use"`

	VGName string `json:"vg_name"`
}
