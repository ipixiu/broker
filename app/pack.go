package main

import (
	"encoding/binary"
	"strconv"
	"strings"
)

import (
	proto "github.com/gogo/protobuf/proto"
	"paipai.cn/bc_goclientsdk/udpbroadcast/BcSystem"
)

func StringIpToInt(ipstring string) uint32 {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt uint32 = 0
	var pos uint32 = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | uint32(tempInt)
		pos -= 8
	}
	return ipInt
}

func PackMessage(packType, level uint32, idType int32, flags BcSystem.DataBCHeadFlag,
	data []byte, flakeID uint64, channels []uint64) []byte {

	header := &BcSystem.DataBCHead{
		// SourceServerID:     proto.Uint32(0),
		// SourceServerTypeID: proto.Uint32(0),
		// MinProtocol:        proto.Int32(0),
		// MaxProtocol:        proto.Int32(0),
		OrderID: proto.Uint64(flakeID),
		Qos:     proto.Uint32(level),
		Flags:   &flags,
		Type:    &idType,
		Id:      channels,
	}

	headerData, err := proto.Marshal(header)
	if err != nil {
		return nil
	}

	headerlen, bodylen := len(headerData), len(data)

	totallen := len(uint16) + headerlen + len(uint32) + len(uint16)

	message := make([]byte, totallen)

	// header: headerlen + header
	offset := 0
	binary.LittleEndian.PutUint16(message[offset:len(uint16)], uint16(headerlen))
	offset += len(uint16)
	copy(message[offset:], headerData)
	offset += headerlen

	// body: packtype + bodylen + body
	binary.LittleEndian.PutUint32(message[offset:], packType)
	offset += len(uint32)
	binary.LittleEndian.PutUint16(offset:], uint16(bodylen))
	// offset += len(uint16)

	return append(message, data...)
}

