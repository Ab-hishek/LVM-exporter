package lvm

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
	"os/exec"
	"strconv"
	"strings"
)

const (
	VGList = "vgs"
)

func ListLVMVolumeGroup() ([]VolumeGroup, error) {
	args := []string{
		"--options", "all",
		"--reportformat", "json",
		"--units", "b",
	}
	cmd := exec.Command(VGList, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("lvm: list volume group cmd %v: %v", args, err)
		return nil, err
	}
	return decodeVgsJSON(output)
}

func decodeVgsJSON(raw []byte) ([]VolumeGroup, error) {
	output := &struct {
		Report []struct {
			VolumeGroups []map[string]string `json:"vg"`
		} `json:"report"`
	}{}
	var err error
	if err = json.Unmarshal(raw, output); err != nil {
		return nil, err
	}

	if len(output.Report) != 1 {
		return nil, fmt.Errorf("expected exactly one lvm report")
	}

	items := output.Report[0].VolumeGroups
	vgs := make([]VolumeGroup, 0, len(items))
	for _, item := range items {
		var vg VolumeGroup
		if vg, err = parseVolumeGroup(item); err != nil {
			return vgs, err
		}
		vgs = append(vgs, vg)
	}
	return vgs, nil
}

func parseVolumeGroup(m map[string]string) (VolumeGroup, error) {
	var vg VolumeGroup
	var count int
	var sizeBytes int64
	var err error
	var vgSize int64
	var vgFree int64

	vg.Name = m["vg_name"]
	vg.UUID = m["vg_uuid"]

	int32Map := map[string]*int32{
		"pv_count":            &vg.PVCount,
		"lv_count":            &vg.LVCount,
		"max_lv":              &vg.MaxLV,
		"max_pv":              &vg.MaxPV,
		"snap_count":          &vg.SnapCount,
		"vg_missing_pv_count": &vg.MissingPVCount,
		"vg_mda_count":        &vg.MetadataCount,
		"vg_mda_used_count":   &vg.MetadataUsedCount,
	}
	for key, value := range int32Map {
		count, err = strconv.Atoi(m[key])
		if err != nil {
			err = fmt.Errorf("invalid format of %v=%v for vg %v: %v", key, m[key], vg.Name, err)
			return vg, err
		}
		*value = int32(count)
	}

	resQuantityMap := map[string]*resource.Quantity{
		"vg_size":     &vg.Size,
		"vg_free":     &vg.Free,
		"vg_mda_size": &vg.MetadataSize,
		"vg_mda_free": &vg.MetadataFree,
	}

	for key, value := range resQuantityMap {
		sizeBytes, err = strconv.ParseInt(
			strings.TrimSuffix(strings.ToLower(m[key]), "b"),
			10, 64)
		if err != nil {
			err = fmt.Errorf("invalid format of %v=%v for vg %v: %v", key, m[key], vg.Name, err)
			return vg, err
		}
		if key == "vg_size" {
			vgSize = sizeBytes
		} else if key == "vg_free" {
			vgFree = sizeBytes
		}
		quantity := resource.NewQuantity(sizeBytes, resource.BinarySI)
		*value = *quantity
	}
	vg.Used = *resource.NewQuantity(vgSize-vgFree, resource.BinarySI)

	return vg, err
}
