package lvm

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	VGList = "vgs"
	LVList = "lvs"
	PVList = "pvs"

	PVScan = "pvscan"

	LVThinPool = "thin-pool"
)

var (
	Enums = map[string][]string{
		"lv_permissions":       {"unknown", "writeable", "read-only", "read-only-override"},
		"lv_when_full":         {"error", "queue"},
		"raid_sync_action":     {"idle", "frozen", "resync", "recover", "check", "repair"},
		"lv_health_status":     {"", "partial", "refresh needed", "mismatches exist"},
		"vg_allocation_policy": {"normal", "contiguous", "cling", "anywhere", "inherited"},
		"vg_permissions":       {"writeable", "read-only"},
	}
)

// ListLVMVolumeGroup invokes `vgs` to list all the available volume
// groups in the node.
func ListLVMVolumeGroup() ([]VolumeGroup, error) {
	if err := ReloadLVMMetadataCache(); err != nil {
		return nil, err
	}

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
		quantity := resource.NewQuantity(sizeBytes, resource.BinarySI)
		*value = *quantity
	}

	vg.Permission = getIntFieldValue("vg_permissions", m["vg_permissions"])
	vg.AllocationPolicy = getIntFieldValue("vg_allocation_policy", m["vg_allocation_policy"])

	return vg, err
}

// This function returns the integer equivalent for different string values for the LVM component(vg,lv) field
// -1 represents undefined
func getIntFieldValue(fieldName, fieldValue string) int {
	mv := -1
	for i, v := range Enums[fieldName] {
		if v == fieldValue {
			mv = i
		}
	}
	return mv
}

func ListLVMLogicalVolume() ([]LogicalVolume, error) {
	if err := ReloadLVMMetadataCache(); err != nil {
		return nil, err
	}

	args := []string{
		"--options", "lv_all,vg_name,segtype",
		"--reportformat", "json",
		"--units", "b",
	}
	cmd := exec.Command(LVList, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("lvm: error while running command %s %v: %v", LVList, args, err)
		return nil, err
	}
	return decodeLvsJSON(output)
}

// ReloadLVMMetadataCache refreshes lvmetad daemon cache used for
// serving vgs or other lvm utility.
func ReloadLVMMetadataCache() error {
	args := []string{"--cache"}
	cmd := exec.Command(PVScan, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("lvm: reload lvm metadata cache: %v - %v", string(output), err)
		return err
	}
	return nil
}

func parseLogicalVolume(m map[string]string) (LogicalVolume, error) {
	var lv LogicalVolume
	var err error
	var sizeBytes int64
	var count float64

	lv.Name = m["lv_name"]
	lv.FullName = m["lv_full_name"]
	lv.UUID = m["lv_uuid"]
	lv.Path = m["lv_path"]
	lv.DMPath = m["lv_dm_path"]
	lv.VGName = m["vg_name"]
	lv.ActiveStatus = m["lv_active"]

	resQuantityMap := map[string]*resource.Quantity{
		"lv_size":          &lv.Size,
		"lv_metadata_size": &lv.MetadataSize,
	}

	for key, value := range resQuantityMap {
		if m["segtype"] != LVThinPool && key == "lv_metadata_size" {
			sizeBytes = 0
		} else {
			sizeBytes, err = strconv.ParseInt(
				strings.TrimSuffix(strings.ToLower(m[key]), "b"),
				10, 64)
			if err != nil {
				err = fmt.Errorf("invalid format of %v=%v for lv %v: %v", key, m[key], lv.Name, err)
				return lv, err
			}
		}
		quantity := resource.NewQuantity(sizeBytes, resource.BinarySI)
		*value = *quantity
	}

	lv.SegType = m["segtype"]
	lv.Permission = getIntFieldValue("lv_permissions", m["lv_permissions"])
	lv.BehaviourWhenFull = getIntFieldValue("lv_when_full", m["lv_when_full"])
	lv.HealthStatus = getIntFieldValue("lv_health_status", m["lv_health_status"])
	lv.RaidSyncAction = getIntFieldValue("raid_sync_action", m["raid_sync_action"])
	lv.Host = m["lv_host"]
	lv.PoolName = m["pool_lv"]

	float64Map := map[string]*float64{
		"data_percent":     &lv.UsedSizePercent,
		"metadata_percent": &lv.MetadataUsedPercent,
		"snap_percent":     &lv.SnapshotUsedPercent,
	}
	for key, value := range float64Map {
		if m[key] == "" {
			count = 0
		} else {
			count, err = strconv.ParseFloat(m[key], 64)
			if err != nil {
				err = fmt.Errorf("invalid format of %v=%v for lv %v: %v", key, m[key], lv.Name, err)
				return lv, err
			}
		}
		*value = count
	}

	return lv, err
}

func decodeLvsJSON(raw []byte) ([]LogicalVolume, error) {
	output := &struct {
		Report []struct {
			LogicalVolumes []map[string]string `json:"lv"`
		} `json:"report"`
	}{}
	var err error
	if err = json.Unmarshal(raw, output); err != nil {
		return nil, err
	}

	if len(output.Report) != 1 {
		return nil, fmt.Errorf("expected exactly one lvm report")
	}

	items := output.Report[0].LogicalVolumes
	lvs := make([]LogicalVolume, 0, len(items))
	for _, item := range items {
		var lv LogicalVolume
		if lv, err = parseLogicalVolume(item); err != nil {
			return lvs, err
		}
		deviceName, err := getLvDeviceName(lv.Path)
		if err != nil {
			klog.Error(err)
			return nil, err
		}
		lv.Device = deviceName
		lvs = append(lvs, lv)
	}
	return lvs, nil
}

// Function to get LVM Logical volume device
// It returns LVM logical volume device(dm-*).
// This is used as a label in metrics which helps us to map lv_name to device.
//
// Example: my_lv(lv_name) -> dm-0(device)
func getLvDeviceName(path string) (string, error) {
	dmPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		klog.Errorf("failed to resolve device mapper from lv path %v: %v", path, err)
		return "", err
	}
	deviceName := strings.Split(dmPath, "/")
	return deviceName[len(deviceName)-1], nil
}

/*
ListLVMPhysicalVolume invokes `pvs` to list all the available LVM physical volumes in the node.
*/
func ListLVMPhysicalVolume() ([]PhysicalVolume, error) {
	if err := ReloadLVMMetadataCache(); err != nil {
		return nil, err
	}

	args := []string{
		"--options", "pv_all,vg_name",
		"--reportformat", "json",
		"--units", "b",
	}
	cmd := exec.Command(PVList, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		klog.Errorf("lvm: error while running command %s %v: %v", PVList, args, err)
		return nil, err
	}
	return decodePvsJSON(output)
}

/*
To parse the output of pvs command and store it in PhysicalVolume
*/
func parsePhysicalVolume(m map[string]string) (PhysicalVolume, error) {
	var pv PhysicalVolume
	var err error
	var sizeBytes int64

	pv.Name = m["pv_name"]
	pv.UUID = m["pv_uuid"]
	pv.InUse = m["pv_in_use"]
	pv.Allocatable = m["pv_allocatable"]
	pv.Missing = m["pv_missing"]
	pv.VGName = m["vg_name"]

	resQuantityMap := map[string]*resource.Quantity{
		"pv_size":     &pv.Size,
		"pv_free":     &pv.Free,
		"pv_used":     &pv.Used,
		"pv_mda_size": &pv.MetadataSize,
		"pv_mda_free": &pv.MetadataFree,
		"dev_size":    &pv.DeviceSize,
	}

	for key, value := range resQuantityMap {
		sizeBytes, err = strconv.ParseInt(
			strings.TrimSuffix(strings.ToLower(m[key]), "b"),
			10, 64)
		if err != nil {
			err = fmt.Errorf("invalid format of %v=%v for pv %v: %v", key, m[key], pv.Name, err)
			return pv, err
		}
		quantity := resource.NewQuantity(sizeBytes, resource.BinarySI)
		*value = *quantity
	}

	return pv, err
}

/*
Decode json format and store physical volumes in map[string]string
*/
func decodePvsJSON(raw []byte) ([]PhysicalVolume, error) {
	output := &struct {
		Report []struct {
			PhysicalVolume []map[string]string `json:"pv"`
		} `json:"report"`
	}{}
	var err error
	if err = json.Unmarshal(raw, output); err != nil {
		return nil, err
	}

	if len(output.Report) != 1 {
		return nil, fmt.Errorf("expected exactly one lvm report")
	}

	items := output.Report[0].PhysicalVolume
	pvs := make([]PhysicalVolume, 0, len(items))
	for _, item := range items {
		var pv PhysicalVolume
		if pv, err = parsePhysicalVolume(item); err != nil {
			return pvs, err
		}
		pvs = append(pvs, pv)
	}
	return pvs, nil
}
