package ipsave
import "net"
import "net/http"

func RemoteIp(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("X-Real-IP");ip !="" {
		remoteAddr = ip
	}else if ip = r.Header.Get("X-Forwarded-For");ip != ""{
		remoteAddr = ip
	} else {
		remoteAddr,_,_ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}
