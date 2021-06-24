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
	Name string

	// Full name of the lvm logical volume.
	FullName string

	// UUID denotes a unique identity of a lvm logical volume.
	UUID string

	// Size specifies the total size of logical volume in bytes
	Size int64

	// Path specifies LVM logical volume path
	Path string

	// DMPath specifies device mapper path
	DMPath string

	// LVM logical volume device
	Device string

	// Name of the VG in which LVM logical volume is created
	VGName string

	LVLayout string
}
