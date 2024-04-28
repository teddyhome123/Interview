package main

import (
	"fmt"
)

/*
RAID0：將數據分散存儲在多個硬碟上，提高讀寫速度，安全性低。
RAID1：數據在兩個或多個硬碟上有完全相同的副本，提高數據安全性。
RAID10：結合了RAID0和RAID1的特點。
RAID5：數據與Parity信息分散在所有硬碟上，單一硬碟故障可恢復。
RAID6：類似於RAID5，但有兩個Parity，可以容忍兩個硬碟故障。
*/

type Disk struct {
	Data []byte
}

type RAID struct {
	Disks      []*Disk
	Type       string
	StripeSize int
}

func NewRAID(raidType string, numDisks, stripeSize int) *RAID {
	disks := make([]*Disk, numDisks)
	for i := range disks {
		disks[i] = &Disk{Data: make([]byte, stripeSize)}
	}
	return &RAID{Disks: disks, Type: raidType, StripeSize: stripeSize}
}

func (r *RAID) Write(data []byte) {
	switch r.Type {
	case "RAID0":
		for i, b := range data {
			diskIndex := i % len(r.Disks)
			r.Disks[diskIndex].Data = append(r.Disks[diskIndex].Data, b)
		}
	case "RAID1":
		for _, disk := range r.Disks {
			disk.Data = append(disk.Data, data...)
		}
	case "RAID10":
		for i, b := range data {
			pairIndex := (i / r.StripeSize) % (len(r.Disks) / 2)
			r.Disks[2*pairIndex].Data = append(r.Disks[2*pairIndex].Data, b)
			r.Disks[2*pairIndex+1].Data = append(r.Disks[2*pairIndex+1].Data, b)
		}
	}
}

func (r *RAID) Read() string {
	var result []byte
	switch r.Type {
	case "RAID0":
		stripeSize := len(r.Disks[0].Data)
		for i := 0; i < stripeSize; i++ {
			for _, disk := range r.Disks {
				if i < len(disk.Data) {
					result = append(result, disk.Data[i])
				}
			}
		}
	case "RAID1":
		result = r.Disks[0].Data
	case "RAID10":
		numPairs := len(r.Disks) / 2
		stripeSize := len(r.Disks[0].Data)
		for i := 0; i < stripeSize; i++ {
			for j := 0; j < numPairs; j++ {
				pairStartIndex := 2 * j
				if pairStartIndex < len(r.Disks) && i < len(r.Disks[pairStartIndex].Data) {
					result = append(result, r.Disks[pairStartIndex].Data[i])
				}
			}
		}
	}
	return string(result)
}

func main() {
	raid := NewRAID("RAID0", 3, 1024)
	data := []byte("Hello, World")
	raid.Write(data)
	fmt.Println("Data written to RAID0:", raid.Read())

	raid1 := NewRAID("RAID1", 2, 1024)
	raid1.Write(data)
	fmt.Println("Data written to RAID1:", raid1.Read())

	raid10 := NewRAID("RAID10", 4, 1024)
	raid10.Write(data)
	fmt.Println("Data written to RAID10:", raid10.Read())
}
