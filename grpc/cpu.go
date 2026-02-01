package grpc

import (
	"database/sql"

	"github.com/google/uuid"
)

type CpuInfo struct {
	CpuUuid       string
	CpuCoreNum    int32
	CpuUsePercent float64
}

func selectCpuTable(db *sql.DB, clientId string) ([]CpuInfo, error) {
	query, err := db.Query("SELECT cpu_uuid, cpu_core_num, cpu_use_percent FROM public.tb_cpu WHERE client_id = $1", clientId)
	if err != nil {
		return []CpuInfo{}, err
	}
	defer query.Close()

	var dbCpuInfo []CpuInfo
	for query.Next() {
		var uuid string
		var coreNum int32
		var usePercent float64
		err = query.Scan(&uuid, &coreNum, &usePercent)
		if err != nil {
			return []CpuInfo{}, err
		}

		dbCpuInfo = append(dbCpuInfo, CpuInfo{
			CpuUuid:       uuid,
			CpuCoreNum:    coreNum,
			CpuUsePercent: usePercent,
		})
	}

	return dbCpuInfo, nil
}

func insertCpu(db *sql.DB, clientId string, coreNum int32, usePercent float64) error {
	cpuUuid := uuid.New().String()
	_, err := db.Exec("INSERT INTO public.tb_cpu (cpu_uuid, client_id, cpu_core_num, cpu_use_percent) VALUES ($1, $2, $3, $4)",
		cpuUuid, clientId, coreNum, usePercent)
	return err
}

func updateCpu(db *sql.DB, cpuUuid string, usePercent float64) error {
	_, err := db.Exec("UPDATE public.tb_cpu SET cpu_use_percent = $1 WHERE cpu_uuid = $2",
		usePercent, cpuUuid)
	return err
}

func deleteCpu(db *sql.DB, cpuUuid string) error {
	_, err := db.Exec("DELETE FROM public.tb_cpu WHERE cpu_uuid = $1", cpuUuid)
	return err
}

func deleteCpuByClientId(db *sql.DB, clientId string) error {
	_, err := db.Exec("DELETE FROM public.tb_cpu WHERE client_id = $1", clientId)
	return err
}
