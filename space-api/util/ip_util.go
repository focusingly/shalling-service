package util

import (
	"fmt"
	"net"
	"space-api/pack/resource"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var _ipSearcher *xdb.Searcher

func init() {
	if searcher, err := xdb.NewWithBuffer(resource.IpData); err != nil {
		panic(err)
	} else {
		_ipSearcher = searcher
	}
}

func GetIpSearcher() *xdb.Searcher {
	return _ipSearcher
}

// Ipv4Str2U32 将表示 ipv4 地址的字符串转换为 int32 值, 如果解析错误, err 不为 nil
//
//	if i32, err := Ipv4Str2U32("207.67.34.8"); err != nil {
//		return
//	} else {
//		fmt.Println("addr to uint32 is: ", i32)
//	}
func Ipv4Str2U32(ipStr string) (uint32, error) {
	// 将 IPv4 地址字符串解析为 net.IP 类型
	ip := net.ParseIP(ipStr)
	// 检查是否是有效的 IPv4 地址
	if ip == nil || ip.To4() == nil {
		return 0, fmt.Errorf("invalid IPv4 address")
	}

	// 将 IPv4 地址转换为 uint32
	return uint32(ip[12])<<24 | uint32(ip[13])<<16 | uint32(ip[14])<<8 | uint32(ip[15]), nil
}
