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
	if opts.Debug {
		log.Println(params...)
	}
}

func pl(params ...interface{}) {
	if opts.Verbose {
		fmt.Println(params...)
	}
}

func pf(msg string, params ...interface{}) {
	if opts.Verbose {
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
func IPs(class_c_only bool) string {
	var ret_addr []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		addr_str := addr.String()
		if class_c_only && !strings.Contains(addr_str, "192.") {continue}
		if strings.Contains(addr_str, ".") {
			ret := strings.Split(addr_str, "/")[0]
			ret_addr = append(ret_addr, ret)
		}
	}
	return strings.Join(ret_addr, ", ")
}
