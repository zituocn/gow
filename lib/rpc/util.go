package rpc

import (
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// byte2int  byte -> int
func byte2int(b []byte) int {
	s := string(b)
	i, _ := strconv.ParseInt(s, 10, 32)
	return int(i)
}

// GetLocalIP get local ip address
//	returns address and error
func GetLocalIP() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addr {
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("not get local ip")
}

// serverNameKey 处理 serviceName 为 serverNameKey提供给注册服务使用
func serverNameKey(serviceName string) string {
	if string(serviceName[0]) != "/" {
		serviceName = "/" + serviceName
	}
	if string(serviceName[len(serviceName)-1]) != "/" {
		serviceName = serviceName + "/"
	}
	return serviceName
}

// GetRpcFuncName
func getRpcFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	if len(pc) < 1 {
		return ""
	}
	f := runtime.FuncForPC(pc[0])
	fName := f.Name()
	fList := strings.Split(fName, ".")
	if len(fList) < 1 {
		return ""
	}
	return fList[len(fList)-1]
}

// str2int64 string -> int64
func str2int64(str string) int64 {
	if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		return i
	}
	return 0
}

// int642str int64 -> string
func int642str(i int64) string {
	return strconv.FormatInt(i, 10)
}

// GetIp return ip address
func GetIp() (string, error) {
	ip, err := externalIP()
	return ip.String(), err
}

// externalIP
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loop back interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, fmt.Errorf("get ip address faild")
}

// getIpFromAddr
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

// LocalIPv4s  return all non-loop back IPv4 addresses
func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

const (
	workerBits  uint8 = 10                      // 节点数
	seqBits     uint8 = 12                      // 1毫秒内可生成的id序号的二进制位数
	workerMax   int64 = -1 ^ (-1 << workerBits) // 节点ID的最大值，用于防止溢出
	seqMax      int64 = -1 ^ (-1 << seqBits)    // 同上，用来表示生成id序号的最大值
	timeShift   uint8 = workerBits + seqBits    // 时间戳向左的偏移量
	workerShift uint8 = seqBits                 // 节点ID向左的偏移量
	epoch       int64 = 1567906170596           // 开始运行时间
)

type worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	seq       int64
}

func ID() string {
	w := newWorker(1)
	return int642str(w.next())
}

// NewWorker return a new worker
func newWorker(workerId int64) *worker {
	if workerId < 0 || workerId > workerMax {
		return nil
	}
	return &worker{
		timestamp: 0,
		workerId:  workerId,
		seq:       0,
	}
}

// Next 获取一个新IDutils.NewWorker(1)
func (w *worker) next() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()                // 生成完成后记得 解锁 解锁 解锁
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if w.timestamp == now {
		w.seq = (w.seq + 1) & seqMax
		if w.seq == 0 {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.seq = 0
	}
	w.timestamp = now
	ID := int64((now-epoch)<<timeShift | (w.workerId << workerShift) | (w.seq))
	return ID + int64(w.getRankInt(6))
}

// getRankInt
func (w *worker) getRankInt(e int) int {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(e)
	return x
}
