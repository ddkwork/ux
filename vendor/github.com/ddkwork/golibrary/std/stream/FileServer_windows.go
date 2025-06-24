package stream

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"syscall"

	"github.com/ddkwork/golibrary/std/mylog"
)

func FileServerWindowsDisk() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if !strings.HasSuffix(path, "/") && strings.Count(path, "/") == 1 {
			http.Redirect(w, req, path+"/", http.StatusFound)
			return
		}
		if path != "/" {

			h := path[1:2]
			t := path[2:]
			name := h + ":" + t

			http.ServeFile(w, req, name)
			return
		}
		mylog.Check2(io.WriteString(w, "<p>共享硬盘：</p>"))
		for s := range GetWindowsLogicalDrives() {
			s = s[0:1]
			mylog.Check2(io.WriteString(w, "<a href='"+s+"/'>"+s+"</a>\n"))
		}
	})
	ps := GetLocalIPs()
	mylog.Info("Listening on", "http://"+ps[0].To4().String()+DefaultFileServerPort)
	mylog.Check(http.ListenAndServe(":8080", nil))
}

func GetLogicalDrives() []string {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	GetLogicalDrives := kernel32.MustFindProc("GetLogicalDrives")
	n, _ := mylog.Check3(GetLogicalDrives.Call())
	s := strconv.FormatInt(int64(n), 2)

	drivesAll := []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:", "N:", "O:", "P：", "Q：", "R：", "S：", "T：", "U：", "V：", "W：", "X：", "Y：", "Z："}
	temp := drivesAll[0:len(s)]

	var d []string
	for i, v := range s {
		if v == 49 {
			l := len(s) - i - 1
			d = append(d, temp[l])
		}
	}

	var drives []string
	for i, v := range d {
		drives = append(drives[i:], append([]string{v}, drives[:i]...)...)
	}
	return drives
}
