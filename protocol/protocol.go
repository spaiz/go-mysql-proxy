package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

/*
PacketHeader represents packet header
 */
type PacketHeader struct {
	Length uint32
	SequenceId uint8
}

/*
InitialHandshakePacket represents initial handshake packet sent by MySQL Server
 */
type InitialHandshakePacket struct {
	ProtocolVersion   uint8
	ServerVersion     []byte
	ConnectionId      uint32
	AuthPluginData    []byte
	Filler            byte
	CapabilitiesFlags CapabilityFlag
	CharacterSet      uint8
	StatusFlags       uint16
	AuthPluginDataLen uint8
	AuthPluginName    []byte
	header            *PacketHeader
}

// Decode decodes the first packet received from the MySQl Server
// It's a handshake packet
func (r *InitialHandshakePacket) Decode(conn net.Conn) error {
	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		return err
	}

	header := &PacketHeader{}
	ln := []byte{data[0], data[1], data[2], 0x00}
	header.Length = binary.LittleEndian.Uint32(ln)
	// a single byte integer is the same in BigEndian and LittleEndian
	header.SequenceId = data[3]

	r.header = header
	/**
	Assign payload only data to new var just  for convenience
	*/
	payload := data[4:header.Length + 4]
	position := 0
	/**
	As defined in the documentation, this value is alway 10 (0x00 in hex)
	1	[0a] protocol version
	 */
	r.ProtocolVersion = payload[0]
	if r.ProtocolVersion != 0x0a {
		return errors.New("non supported protocol for the proxy. Only version 10 is supported")
	}

	position += 1

	/**
	Extract server version, by finding the terminal character (0x00) index,
	and extracting the data in between
	string[NUL]    server version
	 */
	index := bytes.IndexByte(payload, byte(0x00))
	r.ServerVersion = payload[position: index]
	position = index + 1

	connectionId := payload[position : position + 4]
	id := binary.LittleEndian.Uint32(connectionId)
	r.ConnectionId = id
	position += 4

	/*
	The auth-plugin-data is the concatenation of strings auth-plugin-data-part-1 and auth-plugin-data-part-2.
	 */

	r.AuthPluginData = make([]byte, 8)
	copy(r.AuthPluginData, payload[position: position + 8])

	position += 8

	r.Filler = payload[position]
	if r.Filler != 0x00 {
		return errors.New("failed to decode filler value")
	}

	position += 1

	capabilitiesFlags1 := payload[position: position + 2]
	position += 2

	r.CharacterSet = payload[position]
	position += 1

	r.StatusFlags = binary.LittleEndian.Uint16(payload[position: position + 2])
	position += 2

	capabilityFlags2 := payload[position: position + 2]
	position += 2

	/**
	Reconstruct 32 bit integer from two 16 bit integers.
	Take low 2 bytes and high 2 bytes, ans sum it.
	 */
	capLow := binary.LittleEndian.Uint16(capabilitiesFlags1)
	capHi := binary.LittleEndian.Uint16(capabilityFlags2)
	cap := uint32(capLow) | uint32(capHi) << 16

	r.CapabilitiesFlags = CapabilityFlag(cap)

	if r.CapabilitiesFlags&clientPluginAuth != 0 {
		r.AuthPluginDataLen = payload[position]
		if r.AuthPluginDataLen == 0 {
			return errors.New("wrong auth plugin data len")
		}
	}

	/*
	Skip reserved bytes

	string[10]     reserved (all [00])
	 */

	position += 1 + 10

	/**
	This flag tell us that the client should hash the password using algorithm described here:
	https://dev.mysql.com/doc/internals/en/secure-password-authentication.html#packet-Authentication::Native41
	 */
	if r.CapabilitiesFlags&clientSecureConn != 0 {
		/*
			The auth-plugin-data is the concatenation of strings auth-plugin-data-part-1 and auth-plugin-data-part-2.
		*/
		end := position + Max(13, int(r.AuthPluginDataLen) - 8)
		r.AuthPluginData = append(r.AuthPluginData, payload[position:end]...)
		position = end
	}

	index = bytes.IndexByte(payload[position:], byte(0x00))

	/*
	Due to Bug#59453 the auth-plugin-name is missing the terminating NUL-char in versions prior to 5.5.10 and 5.6.2.
	We know the length of the payload, so if there is no NUL-char, just read all the data until the end
	*/
	if index != -1 {
		r.AuthPluginName = payload[position:position+index]
	} else {
		r.AuthPluginName = payload[position:]
	}

	return nil
}

	// Encode encodes the InitialHandshakePacket to bytes
	func (r InitialHandshakePacket) Encode() ([]byte, error) {
		buf := make([]byte, 0)
		buf = append(buf, r.ProtocolVersion)
		buf = append(buf, r.ServerVersion...)
		buf = append(buf, byte(0x00))

		connectionId := make([]byte, 4)
		binary.LittleEndian.PutUint32(connectionId, r.ConnectionId)
		buf = append(buf, connectionId...)

		//auth1 := make([]byte, 8)
		auth1 := r.AuthPluginData[0:8]
		buf = append(buf, auth1...)
		buf = append(buf, 0x00)

		cap := make([]byte, 4)
		binary.LittleEndian.PutUint32(cap, uint32(r.CapabilitiesFlags))

		cap1 := cap[0:2]
		cap2 := cap[2:]

		buf = append(buf, cap1...)
		buf = append(buf, r.CharacterSet)

		statusFlag := make([]byte, 2)
		binary.LittleEndian.PutUint16(statusFlag, r.StatusFlags)
		buf = append(buf, statusFlag...)
		buf = append(buf, cap2...)
		buf = append(buf, r.AuthPluginDataLen)

		reserved := make([]byte, 10)
		buf = append(buf, reserved...)
		buf = append(buf, r.AuthPluginData[8:]...)
		buf = append(buf, r.AuthPluginName...)
		buf = append(buf, 0x00)

		h := PacketHeader{
			Length:     uint32(len(buf)),
			SequenceId: r.header.SequenceId,
		}

		newBuf := make([]byte, 0, h.Length + 4)

		ln := make([]byte, 4)
		binary.LittleEndian.PutUint32(ln, h.Length)

		newBuf = append(newBuf, ln[:3]...)
		newBuf = append(newBuf, h.SequenceId)
		newBuf = append(newBuf, buf...)

		return newBuf, nil
	}

func (r InitialHandshakePacket) String() string {
	return r.CapabilitiesFlags.String()
}
