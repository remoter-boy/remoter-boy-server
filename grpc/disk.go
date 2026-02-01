package grpc

import (
	"database/sql"

	"github.com/google/uuid"
)

type DiskInfo struct {
	DiskUuid   string
	Device     string
	MountPoint string
	UsePercent float64
	UseGb      float64
	TotalGb    float64
}

func selectDiskTable(db *sql.DB, clientId string) ([]DiskInfo, error) {
	query, err := db.Query("SELECT disk_uuid, disk_device, disk_mount_point, disk_use_percent, disk_use_gb, disk_total_gb FROM public.tb_disk WHERE client_id = $1", clientId)
	if err != nil {
		return []DiskInfo{}, err
	}
	defer query.Close()

	var dbDiskInfo []DiskInfo
	for query.Next() {
		var diskUuid string
		var device string
		var mountPoint string
		var usePercent float64
		var useGb float64
		var totalGb float64
		err = query.Scan(&diskUuid, &device, &mountPoint, &usePercent, &useGb, &totalGb)
		if err != nil {
			return []DiskInfo{}, err
		}

		dbDiskInfo = append(dbDiskInfo, DiskInfo{
			DiskUuid:   diskUuid,
			Device:     device,
			MountPoint: mountPoint,
			UsePercent: usePercent,
			UseGb:      useGb,
			TotalGb:    totalGb,
		})
	}

	return dbDiskInfo, nil
}

func insertDisk(db *sql.DB, clientId, device, mountPoint string, usePercent, useGb, totalGb float64) error {
	diskUuid := uuid.New().String()
	_, err := db.Exec("INSERT INTO public.tb_disk (disk_uuid, client_id, disk_device, disk_mount_point, disk_use_percent, disk_use_gb, disk_total_gb) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		diskUuid, clientId, device, mountPoint, usePercent, useGb, totalGb)
	return err
}

func updateDisk(db *sql.DB, diskUuid string, usePercent, useGb, totalGb float64) error {
	_, err := db.Exec("UPDATE public.tb_disk SET disk_use_percent = $1, disk_use_gb = $2, disk_total_gb = $3 WHERE disk_uuid = $4",
		usePercent, useGb, totalGb, diskUuid)
	return err
}

func deleteDisk(db *sql.DB, diskUuid string) error {
	_, err := db.Exec("DELETE FROM public.tb_disk WHERE disk_uuid = $1", diskUuid)
	return err
}

func deleteDiskByClientId(db *sql.DB, clientId string) error {
	_, err := db.Exec("DELETE FROM public.tb_disk WHERE client_id = $1", clientId)
	return err
}
