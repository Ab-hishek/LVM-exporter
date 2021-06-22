package lvm

import (
	"k8s.io/apimachinery/pkg/api/resource"
)
// VolumeGroup specifies attributes of a given vg exists on node.
type VolumeGroup struct {
	// Name of the lvm volume group.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// UUID denotes a unique identity of a lvm volume group.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	UUID string `json:"uuid"`

	// Size specifies the total size of volume group.
	// +kubebuilder:validation:Required
	Size resource.Quantity `json:"size"`
	// Free specifies the available capacity of volume group.
	// +kubebuilder:validation:Required
	Free resource.Quantity `json:"free"`

	// LVCount denotes total number of logical volumes in
	// volume group.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	LVCount int32 `json:"lvCount"`
	// PVCount denotes total number of physical volumes
	// constituting the volume group.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	PVCount int32 `json:"pvCount"`
}
