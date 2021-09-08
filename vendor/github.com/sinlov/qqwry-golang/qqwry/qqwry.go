package qqwry

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

// catch of qqwry.dat in mem
var DatData fileData

// init DatFile to add qqwry.dat
//	qqwry.DatData.FilePath = "path.dat"
// this method must use before qqwry.NewQQwry
// how to confirm qqwry.dat file that the load was successful?
// can use as this
//	init := qqwry.DatData.InitDatFile()
//	if v, ok := init.(error); ok {
//		if v != nil {
//			log.Fatalf("init InitDatFile error %s", v)
//			return
//		}
//	}
func (f *fileData) InitDatFile() (rs interface{}) {
	// check file exist
	_, err := os.Stat(f.FilePath)
	if err != nil && os.IsNotExist(err) {
		rs = errors.New(fmt.Sprintf(ERROR_STR_DAT_FILE_NOT_EXIST, f.FilePath))
		return
	}

	startTime := time.Now().UnixNano()
	// open file
	f.Path, err = os.OpenFile(f.FilePath, os.O_RDONLY, 0400)
	if err != nil {
		rs = err
		return
	}
	defer f.Path.Close()

	// catch
	tmpData, err := ioutil.ReadAll(f.Path)
	if err != nil {
		rs = err
		return
	}

	f.Data = tmpData

	buf := f.Data[0:8]
	start := binary.LittleEndian.Uint32(buf[:4])
	end := binary.LittleEndian.Uint32(buf[4:])

	f.IPCount = int64((end-start)/INDEX_LEN + 1)

	endTime := time.Now().UnixNano()
	f.LoadTimeMs = float64(endTime-startTime) / 1000000
	return true
}

// new QQwry by qqwry.dat
// use after qqwry.InitDatFile
func NewQQwry() (qqwry *QQwry) {
	qqwry = &QQwry{
		Data: &DatData,
	}
	return
}

// search by ipv4 string
//	res := qqwry.NewQQwry().SearchByIPv4("ipv4")
// res see struct qqwry.ResQQwry
func (q *QQwry) SearchByIPv4(ip string) (res ResQQwry) {

	res = ResQQwry{}
	res.IP = ip
	if strings.Count(ip, ".") != 3 {
		res.Err = ERROR_SEARCH_STR_NOT_IP_V4
		return res
	}
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	// log.Println("loc offset:", offset)
	if offset <= 0 {
		res.Err = ERROR_CAN_NOT_FIND_IP_V4_OFFSET
		return
	}

	var country []byte
	var area []byte

	mode := q.readMode(offset + 4)
	if mode == REDIRECT_MODE_1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == REDIRECT_MODE_2 {
			c := q.readUInt24()
			country = q.readString(c)
			countryOffset += 4
		} else {
			country = q.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = q.readArea(countryOffset)
	} else if mode == REDIRECT_MODE_2 {
		countryOffset := q.readUInt24()
		country = q.readString(countryOffset)
		area = q.readArea(offset + 8)
	} else {
		country = q.readString(offset + 4)
		area = q.readArea(offset + uint32(5+len(country)))
	}

	enc := mahonia.NewDecoder("gbk")
	res.Country = enc.ConvertString(string(country))
	res.Area = enc.ConvertString(string(area))
	return
}

// read area
func (q *QQwry) readArea(offset uint32) []byte {
	mode := q.readMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		areaOffset := q.readUInt24()
		if areaOffset == 0 {
			return []byte("")
		} else {
			return q.readString(areaOffset)
		}
	} else {
		return q.readString(offset)
	}
	return []byte("")
}

// read data string byte array
func (q *QQwry) readString(offset uint32) []byte {
	q.setOffset(int64(offset))
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = q.readData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

// qqwry.data read mode
func (q *QQwry) readMode(offset uint32) byte {
	mode := q.readData(1, int64(offset))
	return mode[0]
}

// search index of qqwry.dat
func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.readData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, INDEX_LEN)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = q.getMiddleOffset(start, end)
		buf = q.readData(INDEX_LEN, int64(mid))
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == INDEX_LEN {
			offset := byte3ToUInt32(buf[4:])
			buf = q.readData(INDEX_LEN)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		// greater than ip, so move before
		if _ip > ip {
			end = mid
		} else if _ip < ip { // less than ip, so move after
			start = mid
		} else if _ip == ip {
			return byte3ToUInt32(buf[4:])
		}
	}
}

// private read uint32 like uint24
func (q *QQwry) readUInt24() uint32 {
	buf := q.readData(3) // in qqwry offset use 3
	return byte3ToUInt32(buf)
}

// read data by qqwry.dat define, just like num and data offset
func (q *QQwry) readData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.setOffset(offset[0])
	}
	nums := int64(num)
	end := q.Offset + nums
	dataNum := int64(len(q.Data.Data))
	if q.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = q.Data.Data[q.Offset:end]
	q.Offset = end
	return
}

// set dat Offset in qqwry.dat
func (q *QQwry) setOffset(offset int64) {
	q.Offset = offset
}

// find dat middle offset
func (q *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / INDEX_LEN) >> 1
	return start + records*INDEX_LEN
}

// let byte array to uint32 by each count 3
func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
