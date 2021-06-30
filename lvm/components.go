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

	// MaxLV denotes maximum number of Logical volumes allowed
	// in volume group or 0 if unlimited.
	MaxLV int32 `json:"max_lv"`

	// MaxPV denotes Maximum number of physical volumes allowed
	// in volume group or 0 if unlimited.
	MaxPV int32 `json:"max_pv"`

	// SnapCount denotes number of snapshots in volume group.
	SnapCount int32 `json:"snap_count"`

	// MissingPVCount denotes number of physical volumes in
	// volume group which are missing.
	MissingPVCount int32 `json:"vg_missing_pv_count"`

	// MetadataCount denotes number of metadata areas on the
	// volume group.
	MetadataCount int32 `json:"vg_mda_count"`

	// MetadataUsedCount denotes number of metadata areas in
	// volume group
	MetadataUsedCount int32 `json:"vg_mda_used_count"`

	// MetadataFree specifies the available metadata area space
	// for the volume group
	MetadataFree resource.Quantity `json:"vg_mda_free"`

	// MetadataSize specifies size of smallest metadata area
	// for the volume group
	MetadataSize resource.Quantity `json:"vg_mda_size"`

	// Permission indicates the volume group permission
	// which can be writable or read-only
	Permission int `json:"vg_permissions"`

	// AllocationPolicy indicates the volume group allocation
	// policy(normal/contiguous/cling/anywhere/inherited)
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

	// SegType specifies the type of Logical volume segment
	SegType string `json:"segtype"`

	// Permission indicates the logical volume permission
	// which can be either one of these- (unknown/writeable/read-only/read-only-override)
	Permission int `json:"lv_permissions"`

	// BehaviourWhenFull indicates the behaviour of thin pools when it is full
	BehaviourWhenFull int `json:"lv_when_full"`

	// HealthStatus indicates the health status of logical volumes.
	// This can be any one among these - (""/partial/refresh needed/mismatches exist)
	HealthStatus int `json:"lv_health_status"`

	// RaidSyncAction indicates the current synchronization action being performed for RAID
	// action can be any one of these - (idle/frozen/resync/recover/check/repair)
	RaidSyncAction int `json:"raid_sync_action"`

	// ActiveStatus indicates the active state of logical volume
	ActiveStatus string `json:"lv_active"`

	// Host specifies the creation host of the logical volume, if known
	Host string `json:"lv_host"`

	// For thin volumes, the thin pool Logical volume for that volume
	PoolName string `json:"pool_lv"`

	// UsedSizePercent specifies the percentage full for snapshot, cache
	// and thin pools and volumes if logical volume is active.
	UsedSizePercent float64 `json:"data_percent"`

	// MetadataSize specifies the size of the logical volume that holds
	// the metadata for thin and cache pools.
	MetadataSize resource.Quantity `json:"lv_metadata_size"`

	// MetadataUsedPercent specifies the percentage of metadata full if logical volume
	// is active for cache and thin pools.
	MetadataUsedPercent float64 `json:"metadata_percent"`

	// SnapshotUsedPercent specifies the percentage full for snapshots  if
	// logical volume is active.
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

	// DeviceSize specifies the size of underlying device in bytes
	DeviceSize resource.Quantity `json:"dev_size"`

	// MetadataSize specifies the size of smallest metadata area on this device in bytes
	MetadataSize resource.Quantity `json:"pv_mda_size"`

	// MetadataFree specifies the free metadata area space on the device in bytes
	MetadataFree resource.Quantity `json:"pv_mda_free"`

	// Free specifies the physical volume unallocated space in bytes
	Free resource.Quantity `json:"pv_free"`

	// Used specifies the physical volume allocated space in bytes
	Used resource.Quantity `json:"pv_used"`

	// Allocatable indicates whether the device can be used for allocation
	Allocatable string `json:"pv_allocatable"`

	// Missing indicates whether the device is missing in the system
	Missing string `json:"pv_missing"`

	// InUse indicates whether or not the physical volume is in use
	InUse string `json:"pv_in_use"`

	// Name of the volume group which uses this physical volume
	VGName string `json:"vg_name"`
}
