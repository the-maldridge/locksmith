package nm

import "net"

func delAddr(ips []net.IPNet, tgt net.IPNet) []net.IPNet {
	var out []net.IPNet
	for _, i := range ips {
		if i.IP.Equal(tgt.IP) {
			continue
		}
		out = append(out, i)
	}
	return out
}
