package pool

import (
	"fmt"
	"github.com/e11it/ddq/ssh"
)

type Pool struct {
	Name  string
	Hosts []string
	AUser string
	cl    []*ssh.Client
}

func New(name string, auser string) *Pool {
	return &Pool{Name: name, AUser: auser}
}

func (p *Pool) AddHost(h_str string) {
	p.Hosts = append(p.Hosts, h_str)
}

func (p *Pool) Connect() error {
	for _, h := range p.Hosts {
		if cl, err := ssh.ConnectWithPassword(h, p.AUser, "inferion"); err == nil {
			p.cl = append(p.cl, cl)
		}
	}

	return nil
}

func (p *Pool) Exec(cmd string) {
	for _, cl := range p.cl {
		if out, err := cl.Exec(cmd); err == nil {
			fmt.Printf("Exec: %s\n", out)
		} else {
			fmt.Println("Err:", err.Error())
		}
	}
}
