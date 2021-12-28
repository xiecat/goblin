package qqwry

import (
	"os"
	"sync"

	"github.com/sinlov/qqwry-golang/qqwry"

	log "unknwon.dev/clog/v2"
)

type Database struct {
	wry *qqwry.QQwry
}

// Area returns IpArea according to ipctl
func (db *Database) Area(ip string) string {
	defer func() {
		_ = recover()
	}()
	if db.wry == nil {
		return ""
	}
	ipData := db.wry.SearchByIPv4(ip)
	if ipData.Area == " CZ88.NET" {
		return ipData.Country
	}
	return ipData.Country + " " + ipData.Area
}

var wry *qqwry.QQwry
var once sync.Once

func checkUpdate() {
	_, err := os.Stat("qqwry.dat")
	if err != nil {
		if os.IsNotExist(err) {
			err := download()
			if err != nil {
				log.Warn("Download qqwry.dat failed, caused by:%v, recommend to download it by yourself otherwise the `IpArea` will be null", err)
			}
		}
	}
}

func New() *Database {
	once.Do(func() {
		checkUpdate()
		qqwry.DatData.FilePath = "qqwry.dat"
		init := qqwry.DatData.InitDatFile()
		if v, ok := init.(error); ok {
			if v != nil {
				log.Warn("qqwry init failed")
				wry = nil
			}
		}
		wry = qqwry.NewQQwry()
	})
	return &Database{wry: wry}
}
