package grpc

import (
	"context"
	"database/sql"
	"io"
	"log"
	"remoter-boy-server/common"
	remoter "remoter-boy-server/proto_go"
	"time"
)

func UpdateClientInfo(db *sql.DB, clientId string, msg *remoter.InfoRequestMsg) error {
	dbCpuList, err := selectCpuTable(db, clientId)
	if err != nil {
		return err
	}

	dbCpuMap := make(map[int32]CpuInfo)
	for _, cpu := range dbCpuList {
		dbCpuMap[cpu.CpuCoreNum] = cpu
	}

	msgCpuCores := make(map[int32]bool)
	for _, cpuStat := range msg.Cpu {
		msgCpuCores[cpuStat.Core] = true

		if dbCpu, exists := dbCpuMap[cpuStat.Core]; exists {
			if err := updateCpu(db, dbCpu.CpuUuid, float64(cpuStat.Load)); err != nil {
				return err
			}
		} else {
			if err := insertCpu(db, clientId, cpuStat.Core, float64(cpuStat.Load)); err != nil {
				return err
			}
		}
	}

	for coreNum, dbCpu := range dbCpuMap {
		if !msgCpuCores[coreNum] {
			if err := deleteCpu(db, dbCpu.CpuUuid); err != nil {
				return err
			}
		}
	}

	dbRam, err := selectRamTable(db, clientId)
	if err != nil {
		return err
	}

	if msg.Ram != nil {
		if dbRam != nil {
			if err := updateRam(db, dbRam.RamUuid, float64(msg.Ram.UsedPercent), float64(msg.Ram.UsedGb), float64(msg.Ram.TotalGb)); err != nil {
				return err
			}
		} else {
			if err := insertRam(db, clientId, float64(msg.Ram.UsedPercent), float64(msg.Ram.UsedGb), float64(msg.Ram.TotalGb)); err != nil {
				return err
			}
		}
	} else if dbRam != nil {
		if err := deleteRam(db, dbRam.RamUuid); err != nil {
			return err
		}
	}

	dbDiskList, err := selectDiskTable(db, clientId)
	if err != nil {
		return err
	}

	dbDiskMap := make(map[string]DiskInfo)
	for _, disk := range dbDiskList {
		key := disk.Device + "|" + disk.MountPoint
		dbDiskMap[key] = disk
	}

	msgDiskKeys := make(map[string]bool)
	for _, diskStat := range msg.Disk {
		key := diskStat.Device + "|" + diskStat.MountPoint
		msgDiskKeys[key] = true

		if dbDisk, exists := dbDiskMap[key]; exists {
			if err := updateDisk(db, dbDisk.DiskUuid, float64(diskStat.UsedPercent), float64(diskStat.UsedGb), float64(diskStat.TotalGb)); err != nil {
				return err
			}
		} else {
			if err := insertDisk(db, clientId, diskStat.Device, diskStat.MountPoint, float64(diskStat.UsedPercent), float64(diskStat.UsedGb), float64(diskStat.TotalGb)); err != nil {
				return err
			}
		}
	}

	for key, dbDisk := range dbDiskMap {
		if !msgDiskKeys[key] {
			if err := deleteDisk(db, dbDisk.DiskUuid); err != nil {
				return err
			}
		}
	}

	dbIpList, err := selectInternetProtocolTable(db, clientId)
	if err != nil {
		return err
	}

	dbIpMap := make(map[string]IpDbInfo)
	for _, ip := range dbIpList {
		dbIpMap[ip.IpV4] = ip
	}

	msgIps := make(map[string]bool)
	if msg.Ips != nil {
		for _, ipAddr := range msg.Ips.Ips {
			msgIps[ipAddr] = true

			if _, exists := dbIpMap[ipAddr]; !exists {
				if err := insertIp(db, clientId, ipAddr); err != nil {
					return err
				}
			}
		}
	}

	for ipAddr, dbIp := range dbIpMap {
		if !msgIps[ipAddr] {
			if err := deleteIp(db, dbIp.IpUuid); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) Connect(stream remoter.Remoter_ConnectServer) error {
	_, cancel := context.WithCancel(stream.Context())
	defer cancel()

	db := common.InitDatabase()

	if db == nil {
		panic("Database Connect Error")
	}

	defer db.Close()

	recv, err := stream.Recv()

	if err != nil {
		log.Println("stream.Recv() Error:" + err.Error())
		return err
	}

	_, err = CreateClient(db, recv.ClientId)

	if err != nil {
		log.Println("Client Insert Error:" + err.Error())
		return err
	}

	client := &ClientInfo{
		ID:       recv.ClientId,
		Stream:   stream,
		LastPing: time.Now(),
		Cancel:   cancel,
	}
	s.clients.Store(client.ID, client)
	log.Printf("[Connect] Client: %s", client.ID)

	// Clean up at the end of the connection
	defer func() {
		s.clients.Delete(client.ID)
		_, err = DeleteClient(db, client.ID)

		if err != nil {
			log.Println("Client Delete Error: " + err.Error())
		}
		log.Printf("[Disconnet] Client: %s", client.ID)
	}()

	for {
		recv, err := stream.Recv() // ← 루프 안에서 매번 새 메시지 수신

		if client.LastPing.Before(time.Now().Add(-20 * time.Second)) {
			_, err = DeleteClient(db, client.ID)
			if err != nil {
				log.Println("Client Delete Error: " + err.Error())
			}
			log.Printf("[Disconnet] Client: %s", client.ID)
			return nil
		}

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		client.LastPing = time.Now()

		if err := UpdateClientInfo(db, client.ID, recv); err != nil {
			log.Println("UpdateClientInfo Error: " + err.Error())
			s.clients.Delete(client.ID)
			client.Cancel()
			return err
		}

		// send Response
		if err := stream.Send(&remoter.NilResponseMsg{}); err != nil {
			return err
		}
	}
}
