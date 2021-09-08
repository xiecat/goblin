package geoip

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	log "unknwon.dev/clog/v2"
)

type Database struct {
	geo *geoip2.Reader
}

var geo *geoip2.Reader
var once sync.Once

// Area returns IpArea according to ipctl
func (db *Database) Area(ip string) string {
	defer func() {
		_ = recover()
	}()
	record, err := db.geo.City(net.ParseIP(ip))
	if err != nil {
		return ""
	}

	country := record.Country.Names["en"]
	city := record.City.Names["en"]
	if city == "" {
		city = record.Location.TimeZone
	}
	return fmt.Sprintf("%s %s", country, city)
}

func checkUpdate(licenseKey string) {
	info, err := os.Stat("GeoLite2-City.mmdb")
	if err != nil {
		if os.IsNotExist(err) {
			err := download(licenseKey)
			if err != nil {
				log.Warn("Download GeoLite2-City.mmdb failed, caused by:%v, recommend to download it by yourself otherwise the `IpArea` will be null", err)
			}
		}
	} else if -time.Until(info.ModTime()) > 7*24*time.Hour {
		log.Info("Updating GeoLite2-City.mmdb...")
		err := download(licenseKey)
		if err != nil {
			log.Warn("Update GeoLite2-City.mmdb failed, please download GeoLite2-City.mmdb by yourself")
		}
	}
}

func New(licenseKey string) *Database {
	once.Do(func() {
		var err error
		checkUpdate(licenseKey)
		geo, err = geoip2.Open("GeoLite2-City.mmdb")
		if err != nil {
			log.Error("Load GeoLite2-City.mmdb failed, `IpArea` will be null")
		}
	})
	return &Database{geo: geo}
}
