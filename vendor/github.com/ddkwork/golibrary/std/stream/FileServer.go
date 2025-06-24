package stream

import (
	"embed"
	"net"
	"net/http"

	"github.com/ddkwork/golibrary/std/mylog"
)

var DefaultFileServerPort = ":8080"

func FileServer() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	ps := GetLocalIPs()
	mylog.Info("Listening on", "http://"+ps[0].To4().String()+DefaultFileServerPort)
	mylog.Check(http.ListenAndServe(DefaultFileServerPort, nil))
}

func FileServerFS(fs embed.FS) {
	handler := http.FileServerFS(fs)
	ps := GetLocalIPs()
	mylog.Info("Listening on", "http://"+ps[0].To4().String()+DefaultFileServerPort)
	mylog.Check(http.ListenAndServe(DefaultFileServerPort, handler))
}

func isLocalLink(ip net.IP) bool {
	return ip.IsGlobalUnicast() && !ip.IsLoopback() && !ip.IsLinkLocalUnicast()
}

func GetLocalIPs() []net.IP {
	var ips []net.IP
	interfaces := mylog.Check2(net.Interfaces())
	for _, face := range interfaces {
		adders := mylog.Check2(face.Addrs())
		for _, addr := range adders {
			ipNet, ok := addr.(*net.IPNet)
			if ok && isLocalLink(ipNet.IP) {
				if ipNet.IP.To4() != nil {
					ips = append(ips, ipNet.IP)
				}
			}
		}
	}
	return ips
}
