package grpc

import (
	"database/sql"

	"github.com/google/uuid"
)

type RamInfo struct {
	RamUuid       string
	RamUsePercent float64
	RamUseGb      float64
	RamTotalGb    float64
}

func selectRamTable(db *sql.DB, clientId string) (*RamInfo, error) {
	query, err := db.Query("SELECT ram_uuid, ram_use_percent, ram_use_gb, ram_total_gb FROM public.tb_ram WHERE client_id = $1", clientId)

	if err != nil {
		return nil, err
	}
	defer query.Close()

	var checkRamInfo *RamInfo
	for query.Next() {
		var ramUuid string
		var usePercent, useGb, totalGb float64

		err := query.Scan(&ramUuid, &usePercent, &useGb, &totalGb)
		if err != nil {
			return nil, err
		}

		checkRamInfo = &RamInfo{
			RamUuid:       ramUuid,
			RamUsePercent: usePercent,
			RamUseGb:      useGb,
			RamTotalGb:    totalGb,
		}
	}

	return checkRamInfo, nil
}

func insertRam(db *sql.DB, clientId string, usePercent, useGb, totalGb float64) error {
	ramUuid := uuid.New().String()
	_, err := db.Exec("INSERT INTO public.tb_ram (ram_uuid, client_id, ram_use_percent, ram_use_gb, ram_total_gb) VALUES ($1, $2, $3, $4, $5)",
		ramUuid, clientId, usePercent, useGb, totalGb)
	return err
}

func updateRam(db *sql.DB, ramUuid string, usePercent, useGb, totalGb float64) error {
	_, err := db.Exec("UPDATE public.tb_ram SET ram_use_percent = $1, ram_use_gb = $2, ram_total_gb = $3 WHERE ram_uuid = $4",
		usePercent, useGb, totalGb, ramUuid)
	return err
}

func deleteRam(db *sql.DB, ramUuid string) error {
	_, err := db.Exec("DELETE FROM public.tb_ram WHERE ram_uuid = $1", ramUuid)
	return err
}

func deleteRamByClientId(db *sql.DB, clientId string) error {
	_, err := db.Exec("DELETE FROM public.tb_ram WHERE client_id = $1", clientId)
	return err
}
