package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"k8s.io/klog"
)

type Resolver interface {
	Resolve(string) string
}

type Server struct {
	resolver Resolver
	binded   bool
}

func New(resolver Resolver) *Server {
	return &Server{
		resolver: resolver,
	}
}

func (s *Server) bind() {
	if s.binded {
		return
	}
	dns.Handle(fmt.Sprintf("."), s)
	s.binded = true
}

func (s *Server) StartUDP() error {
	s.bind()
	dnsServer := &dns.Server{Addr: fmt.Sprintf(":53"), Net: "udp"}
	return dnsServer.ListenAndServe()
}

func (s *Server) StartTCP() error {
	s.bind()
	dnsServer := &dns.Server{Addr: fmt.Sprintf(":53"), Net: "tcp"}
	return dnsServer.ListenAndServe()
}

func (s *Server) ServeDNS(res dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Compress = false

	if req.Opcode == dns.OpcodeQuery {
		for _, q := range m.Question {
			if q.Qtype != dns.TypeA {
				continue
			}
			ip := s.resolver.Resolve(q.Name[:len(q.Name)-1])
			if ip == "" {
				continue
			}
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip));
			if err != nil {
				fmt.Println("unable to dns.NewRR(): ", err.Error())
				continue
			}
			m.Answer = append(m.Answer, rr)
		}
	}
	if err := res.WriteMsg(m); err != nil {
		klog.Errorf("unable to write dns response: %s", err.Error())
	}
}
