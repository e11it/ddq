package main

import (
	"fmt"

	"github.com/e11it/ddq/pool"
	"github.com/e11it/ddq/ssh"
)

func main() {
	fmt.Println("test")

	p1 := pool.New("test_pool", "im")
	p1.AddHost("192.168.56.1")
	p1.AddHost("192.168.56.1")
	p1.AddHost("192.168.56.1")
	p1.AddHost("192.168.56.1")
	p1.AddHost("192.168.56.1")
	p1.Connect()
	fmt.Println(p1.Name)
	p1.Exec("date")

	if client, err := ssh.ConnectWithPassword("192.168.56.1", "im", "inferion"); err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("---------------")
		if out, err := client.Exec("uptime"); err == nil {
			fmt.Printf("Exec: %s\n", out)
		} else {
			fmt.Println("Err:", err.Error())
		}

		fmt.Println("---------------")
		if out, err := client.Exec("df -m|grep '^/dev/'"); err == nil {
			fmt.Printf("Exec: %s\n", out)
		} else {
			fmt.Println("Err:", err.Error())
		}
		fmt.Println("---------------")
		if out, err := client.Exec("uname -r"); err == nil {
			fmt.Printf("Exec: %s\n", out)
		} else {
			fmt.Println("Err:", err.Error())
		}
		client.Close()
	}

}
