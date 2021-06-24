package lvm

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

// VolumeGroup specifies attributes of a given vg exists on node.
type VolumeGroup struct {
	// Name of the lvm volume group.
	Name string `json:"name"`

	// UUID denotes a unique identity of a lvm volume group.
	UUID string `json:"uuid"`

	// Size specifies the total size of volume group.
	Size resource.Quantity `json:"size"`

	// Free specifies the available capacity of volume group.
	Free resource.Quantity `json:"free"`

	// Free specifies the used capacity of volume group.
	Used resource.Quantity

	// LVCount denotes total number of logical volumes in
	// volume group.
	LVCount int32 `json:"lvCount"`

	// PVCount denotes total number of physical volumes
	// constituting the volume group.
	PVCount int32 `json:"pvCount"`

	MaxLV int32 `json:"max_lv"`

	MaxPV int32 `json:"max_pv"`

	SnapCount int32 `json:"snap_count"`

	MissingPVCount int32 `json:"vg_missing_pv_count"`

	MetadataCount int32 `json:"vg_mda_count"`

	MetadataUsedCount int32 `json:"vg_mda_used_count"`

	MetadataFree resource.Quantity `json:"vg_mda_free"`

	MetadataSize resource.Quantity `json:"vg_mda_size"`
}

// LogicalVolume specifies attributes of a given lv that exists on the node.
type LogicalVolume struct {
	// Name of the lvm logical volume(name: pvc-213ca1e6-e271-4ec8-875c-c7def3a4908d)
	Name string

	// Full name of the lvm logical volume (fullName: linuxlvmvg/pvc-213ca1e6-e271-4ec8-875c-c7def3a4908d)
	FullName string

	// UUID denotes a unique identity of a lvm logical volume.
	UUID string

	// Size specifies the total size of logical volume in Bytes
	Size int64

	// Path specifies LVM logical volume path
	Path string

	// DMPath specifies device mapper path
	DMPath string

	// LVM logical volume device
	Device string

	// Name of the VG in which LVM logical volume is created
	VGName string
}
