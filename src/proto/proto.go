/*
   |    2 bytes     |   4 bytes   |
   |  command define |  data length |
*/

package proto

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

type Proto struct {
	conn net.Conn
	br   *bufio.Reader
	bw   *bufio.Writer
	brw  *bufio.ReadWriter
}

func NewProto(conn net.Conn) *Proto {
	cbr := bufio.NewReader(conn)
	cbw := bufio.NewWriter(conn)
	return &Proto{
		conn: conn,
		br:   cbr,
		bw:   cbw,
		brw:  bufio.NewReadWriter(cbr, cbw),
	}
}

func (pr *Proto) ReadPackage() (uint16, []byte, error) {
	var (
		cmd  uint16
		data []byte
		err  error
	)

	cbytes, err := pr.readLength(2)
	if err != nil {
		return cmd, data, err
	}

	cmd = binary.BigEndian.Uint16(cbytes)

	lbytes, err := pr.readLength(4)
	if err != nil {
		return cmd, data, err
	}

	length := binary.BigEndian.Uint32(lbytes)

	data, err = pr.readLength(int(length))
	if err != nil {
		return cmd, data, err
	}

	return cmd, data, err
}

func (pr *Proto) readLength(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(pr.br, buf)
	return buf, err
}

func (pr *Proto) WritePackage(cmd uint16, data []byte) error {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf[:], cmd)
	_, err := pr.bw.Write(buf)
	if err != nil {
		return err
	}
	length := uint32(len(data))
	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf[:], length)
	_, err = pr.bw.Write(buf)
	if err != nil {
		return err
	}
	_, err = pr.bw.Write(data)
	if err != nil {
		return err
	}
	pr.bw.Flush()
	return nil
}
