package grpc

import (
	"database/sql"

	"github.com/google/uuid"
)

type IpDbInfo struct {
	IpUuid string
	IpV4   string
}

func selectInternetProtocolTable(db *sql.DB, clientId string) ([]IpDbInfo, error) {
	query, err := db.Query("SELECT internet_protocol_uuid, v4 FROM public.tb_internet_protocol WHERE client_id = $1", clientId)
	if err != nil {
		return []IpDbInfo{}, err
	}
	defer query.Close()

	var dbInternetProtocolInfo []IpDbInfo
	for query.Next() {
		var ipUuid string
		var internetProtocol string
		err = query.Scan(&ipUuid, &internetProtocol)
		if err != nil {
			return []IpDbInfo{}, err
		}

		dbInternetProtocolInfo = append(dbInternetProtocolInfo, IpDbInfo{IpUuid: ipUuid, IpV4: internetProtocol})
	}

	return dbInternetProtocolInfo, nil
}

func insertIp(db *sql.DB, clientId, ipV4 string) error {
	ipUuid := uuid.New().String()
	_, err := db.Exec("INSERT INTO public.tb_internet_protocol (internet_protocol_uuid, client_id, v4) VALUES ($1, $2, $3)",
		ipUuid, clientId, ipV4)
	return err
}

func updateIp(db *sql.DB, ipUuid, ipV4 string) error {
	_, err := db.Exec("UPDATE public.tb_internet_protocol SET v4 = $1 WHERE internet_protocol_uuid = $2",
		ipV4, ipUuid)
	return err
}

func deleteIp(db *sql.DB, ipUuid string) error {
	_, err := db.Exec("DELETE FROM public.tb_internet_protocol WHERE internet_protocol_uuid = $1", ipUuid)
	return err
}

func deleteIpByClientId(db *sql.DB, clientId string) error {
	_, err := db.Exec("DELETE FROM public.tb_internet_protocol WHERE client_id = $1", clientId)
	return err
}
