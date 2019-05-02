package nm

import "net"

func delAddr(ips []net.IP, tgt net.IP) []net.IP {
	var out []net.IP
	for _, i := range ips {
		if i.Equal(tgt) {
			continue
		}
		out = append(out, i)
	}
	return out
}
