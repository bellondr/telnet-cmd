package telnet_cmd

import (
	"bytes"
	"net"
	"fmt"
	"time"
)

type Client struct {
	Conn
	buf     [4096]byte
}

type Conn struct {
	conn interface {
		Read(b []byte) (n int, err error)
		Write(b []byte) (n int, err error)
		Close() error
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
	}
	dataReader *internalDataReader
	dataWriter *internalDataWriter
}

func NewClient(addr string) (*Client, error) {
	if addr == "" {
		addr = "127.0.0.1:telnet"
	}
	const proxy = "tcp"

	conn, err := net.Dial(proxy, addr)
	if nil != err {
		fmt.Println("net dial err: ", err.Error())
		return nil, err
	}

	dataReader := newDataReader(conn)
	dataWriter := newDataWriter(conn)

	c := &Client{}
	c.Conn = Conn{
		conn:conn,
		dataReader:dataReader,
		dataWriter:dataWriter,
	}
	return c, nil
}

func (c *Client) Close() {
	c.Conn.conn.Close()
}

func (c *Client) Login(user, pwd string) error {
	out := bytes.NewBuffer(nil)
	var buffer bytes.Buffer
	var p []byte
	var crlfBuffer [2]byte = [2]byte{'\r','\n'}
	crlf := crlfBuffer[:]
	time.Sleep(1 * time.Second)
	var readBuffer [2]byte
	readP := readBuffer[1:]
	for {
		n, err := c.dataReader.Read(readP)
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			break
		}
		LongWrite(out, readP)
		if (readBuffer[0] == '#' || readBuffer[0] == ':') && readBuffer[1] == ' ' {
			break
		}
		readBuffer[0] = readBuffer[1]
	}

	for _, cmd := range []string{user, pwd} {
		buffer.Write([]byte(cmd))
		buffer.Write(crlf)

		p = buffer.Bytes()

		n, err := LongWrite(c.dataWriter, p)
		if nil != err {
			break
		}
		if expected, actual := int64(len(p)), n; expected != actual {
			fmt.Printf("Transmission problem: tried sending %d bytes, but actually only sent %d bytes. \n ", expected, actual)
			return fmt.Errorf("Transmission problem: tried sending %d bytes, but actually only sent %d bytes. \n ", expected, actual)
		}

		buffer.Reset()

		for {
			n, err := c.dataReader.Read(readP)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}
			LongWrite(out, readP)
			if (readBuffer[0] == '#' || readBuffer[0] == ':') && readBuffer[1] == ' ' {
				break
			}
			readBuffer[0] = readBuffer[1]
		}

	}

	time.Sleep(3 * time.Millisecond)
	return nil
}

func (c *Client) RunCmdWithOutput(cmd string) (string, error) {
	time.Sleep(3 * time.Millisecond)
	out := bytes.NewBuffer(nil)
	var buffer bytes.Buffer
	var p []byte
	var crlfBuffer [2]byte = [2]byte{'\r','\n'}
	crlf := crlfBuffer[:]
	var readBuffer [2]byte
	readP := readBuffer[1:]
	buffer.Write([]byte(cmd))
	buffer.Write(crlf)
	p = buffer.Bytes()

	n, err := LongWrite(c.dataWriter, p)
	if nil != err {
		return out.String(), err
	}
	if expected, actual := int64(len(p)), n; expected != actual {
		fmt.Printf("Transmission problem: tried sending %d bytes, but actually only sent %d bytes. \n ", expected, actual)
		return out.String(), fmt.Errorf("Transmission problem: tried sending %d bytes, but actually only sent %d bytes. \n ", expected, actual)
	}

	buffer.Reset()

	for {
		n, err := c.dataReader.Read(readP)
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			break
		}
		LongWrite(out, readP)
		if (readBuffer[0] == '#' || readBuffer[0] == ':') && readBuffer[1] == ' ' {
			break
		}
		readBuffer[0] = readBuffer[1]
	}

	time.Sleep(3 * time.Millisecond)
	return out.String(), nil
}

func (c *Client) RunCmd(cmd string) (error) {
	time.Sleep(3 * time.Millisecond)
	out := bytes.NewBuffer(nil)
	var buffer bytes.Buffer
	var p []byte
	var crlfBuffer [2]byte = [2]byte{'\r','\n'}
	crlf := crlfBuffer[:]
	var readBuffer [2]byte
	readP := readBuffer[1:]
	buffer.Write([]byte(cmd))
	buffer.Write(crlf)
	p = buffer.Bytes()

	n, err := LongWrite(c.dataWriter, p)
	if nil != err {
		return fmt.Errorf("err: %s. output %s", err.Error(), out.String())
	}
	if expected, actual := int64(len(p)), n; expected != actual {
		fmt.Printf("Transmission problem: tried sending %d bytes, but actually only sent %d bytes. \n ", expected, actual)
		return fmt.Errorf("err: %s. output %s", err.Error(), out.String())
	}

	buffer.Reset()

	for {
		n, err := c.dataReader.Read(readP)
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			break
		}
		LongWrite(out, readP)
		if (readBuffer[0] == '#' || readBuffer[0] == ':') && readBuffer[1] == ' ' {
			break
		}
		readBuffer[0] = readBuffer[1]
	}

	time.Sleep(3 * time.Millisecond)
	return nil
}