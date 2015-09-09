package main

import(
	"fmt"
	"time"
	"crypto/sha1"
	"strings"
	"log"
	"net"
)

// --- Some Shortcuts for often used output functions ---

var fpf = fmt.Printf
var fpl = fmt.Println
var lpl = log.Println
var lf = log.Fatal

func pd(params ...interface{}) {
	if opts_intf["debug"].(bool) {
		log.Println(params...)
	}
}

func pl(params ...interface{}) {
	if opts_intf["verbose"].(bool) {
		fmt.Println(params...)
	}
}

func pf(msg string, params ...interface{}) {
	if opts_intf["verbose"].(bool) {
		fmt.Printf(msg, params...)
	}
}

// --- Crypto ---

func random_sha1() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte("%$" + time.Now().String() + "e{")))
}

//func hashPassword(pword string, salt string) string {
//	return fmt.Sprintf("%x", sha1.Sum([]byte("[--]" + pword + "e{" + salt)))
//
//}

func short_sha(sha string) string{
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

func trim_whitespace(in_str string) string {
	return strings.Trim(in_str, " \n\r\t")
}

// Return all IPs on Geoforce subnets
func IPs() string {
	var ret_addr []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		addr_str := addr.String()
		if strings.Contains(addr_str, "192.168") {
			ret := strings.Split(addr_str, "/")[0]
			ret_addr = append(ret_addr, ret)
		}
	}
	return strings.Join(ret_addr, ", ")
}
