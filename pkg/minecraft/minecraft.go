package minecraft

import (
	"github.com/Tnze/go-mc/net"
	"strings"
)

type Connection struct {
	Client net.RCONClientConn
	config Config
}

type Config struct {
	Address  string
	Password string
}

func NewConnection(c Config) (*Connection, error) {
	client, err := net.DialRCON(c.Address, c.Password)
	if err != nil {
		return &Connection{
			Client: nil,
			config: c,
		}, err
	}

	return &Connection{Client: client, config: c}, nil
}

func (c *Connection) CloseConnection() error {
	if c.Client == nil {
		return &ConnectionError{Reason: "Client not available"}
	}

	return c.Client.Close()
}

func (c *Connection) reconnect() (e error) {
	client, e := net.DialRCON(c.config.Address, c.config.Password)
	if e != nil {
		return
	}
	c.Client = client
	return
}

func (c *Connection) ExecuteCommand(cmd string) (string, error) {
	if c.Client == nil {
		err := c.reconnect()
		if err != nil {
			return "", &ConnectionError{Reason: "No connection to server was ever established and could not be now"}
		}
	}

	err := c.Client.Cmd(cmd)
	if err != nil {
		return c.handleConnectionError(err, cmd)
	}

	resp, err := c.Client.Resp()
	if err != nil {
		return c.handleConnectionError(err, cmd)
	}

	return resp, nil
}

func (c *Connection) handleConnectionError(err error, cmd string) (string, error) {
	str := err.Error()
	if strings.HasPrefix(str, "connect fail") || err.Error() == "read packet length fail: EOF" || strings.HasSuffix(err.Error(), "broken pipe") {
		e := c.reconnect()
		if e != nil {
			return "", &ConnectionError{Reason: "Connection to server could not be established."}
		}

		return c.ExecuteCommand(cmd)
	}

	return "", err
}

type ConnectionError struct {
	Reason string
}

func (c *ConnectionError) Error() string {
	return c.Reason
}
