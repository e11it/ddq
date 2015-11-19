package ssh

import (
	"fmt"
	"net"
	"os/user"
	"time"

	"golang.org/x/crypto/ssh"
)

const DefaultTimeout = 30 * time.Second

type Client struct {
	SSHClient  *ssh.Client
	SSHSession *ssh.Session
}

// Connect with a password. If username is empty simplessh will attempt to get the current user.
func ConnectWithPassword(host, username, pass string) (*Client, error) {
	return ConnectWithPasswordTimeout(host, username, pass, DefaultTimeout)
}

// Same as ConnectWithPassword but allows a custom timeout. If username is empty simplessh will attempt to get the current user.
func ConnectWithPasswordTimeout(host, username, pass string, timeout time.Duration) (*Client, error) {
	authMethod := ssh.Password(pass)

	return connect(host, username, authMethod, timeout)
}

//
func connect(host string, username string, authMethod ssh.AuthMethod, timeout time.Duration) (*Client, error) {
	if username == "" {
		user, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("Username wasn't specified and couldn't get current user: %v", err)
		}

		username = user.Username
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{authMethod},
	}

	host = addPortToHost(host)

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return nil, err
	}
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, host, config)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(sshConn, chans, reqs)

	c := &Client{SSHClient: client}
	return c, nil
}

// Execute cmd on the remote host and return stderr and stdout
func (c *Client) Exec(cmd string) ([]byte, error) {
	if c.SSHSession != nil {
		c.SSHSession.Close()
	}

	session, err := c.SSHClient.NewSession()
	if err != nil {
		return nil, err
	}

	c.SSHSession = session
	return session.CombinedOutput(cmd)
}

func (c *Client) Close() {
	if c.SSHSession != nil {
		c.SSHSession.Close()
	}
}

func addPortToHost(host string) string {
	_, _, err := net.SplitHostPort(host)

	// We got an error so blindly try to add a port number
	if err != nil {
		return net.JoinHostPort(host, "22")
	}

	return host
}
