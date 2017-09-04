package libjmy

import (
	_ "bytes"
	_ "fmt"
	"time"
)

const (
	JMY_PACKET_NONE           = iota
	JMY_PACKET_USER_DATA      = 5
	JMY_PACKET_ACK            = 6
	JMY_PACKET_HEARTBEAT      = 7
	JMY_PACKET_DISCONNECT     = 8
	JMY_PACKET_DISCONNECT_ACK = 9
	JMY_PACKET_USER_ID_DATA   = 10
)

const JMY_PACKET_LEN_HEAD = 2
const JMY_PACKET_LEN_TYPE = 1
const JMY_PACKET_LEN_USER_ID = 4
const JMY_PACKET_LEN_MSG_ID = 2
const JMY_PACKET_LEN_TIMESTAMP = 4

func PACK_UINT16_TO_BUFF(id uint16, buf []byte) int {
	buf[0] = byte((id >> 8) & 0xff)
	buf[1] = byte(id & 0xff)
	return 2
}

func UNPACK_UINT16_FROM_BUFF(buf []byte) uint16 {
	return (uint16(buf[0]<<8) & 0xff00) + uint16(buf[1]&0xff)
}

func PACK_UINT32_TO_BUFF(id uint32, buf []byte) int {
	buf[0] = byte((id >> 24) & 0xff)
	buf[1] = byte((id >> 16) & 0xff)
	buf[2] = byte((id >> 8) & 0xff)
	buf[3] = byte(id & 0xff)
	return 4
}

func UNPACK_UINT32_FROM_BUFF(buf []byte) uint32 {
	var value uint32 = uint32(buf[0]<<24) & 0xff000000
	value += (uint32(buf[1]<<16) & 0xff0000)
	value += (uint32(buf[2]<<8) & 0xff00)
	value += uint32(buf[3] & 0xff)
	return value
}

func PACK_UINT64_TO_BUFF(value uint64, buf []byte) int {
	buf[0] = byte((value >> 56) & 0xff)
	buf[1] = byte((value >> 48) & 0xff)
	buf[2] = byte((value >> 40) & 0xff)
	buf[3] = byte((value >> 32) & 0xff)
	buf[4] = byte((value >> 24) & 0xff)
	buf[5] = byte((value >> 16) & 0xff)
	buf[6] = byte((value >> 8) & 0xff)
	buf[7] = byte(value & 0xff)
	return 8
}

func UNPACK_UINT64_FROM_BUFF(buf []byte) uint64 {
	var value uint64 = uint64(buf[0]<<56) & 0xff00000000000000
	value += (uint64(buf[1]<<48) & 0xff000000000000)
	value += (uint64(buf[2]<<40) & 0xff0000000000)
	value += (uint64(buf[3]<<32) & 0xff00000000)
	value += (uint64(buf[4]<<24) & 0xff000000)
	value += (uint64(buf[5]<<16) & 0xff0000)
	value += (uint64(buf[6]<<8) & 0xff00)
	value += uint64(buf[7] & 0xff)
	return value
}

func jmy_net_proto_get_packet_len_type(buf []byte) (res bool, pack_type int, pack_len uint16) {
	if len(buf) < int(JMY_PACKET_LEN_HEAD+JMY_PACKET_LEN_TYPE) {
		return false, 0, 0
	}
	pack_len = UNPACK_UINT16_FROM_BUFF(buf)
	pack_type = int(buf[2])
	res = true
	return
}

func jmy_net_proto_user_data_pack_len() int {
	return JMY_PACKET_LEN_HEAD + JMY_PACKET_LEN_TYPE + JMY_PACKET_LEN_MSG_ID
}

func jmy_net_proto_user_data_full_pack_len(data_len uint16) int {
	return jmy_net_proto_user_data_pack_len() + int(data_len)
}

func jmy_net_proto_disconnect_pack_len() int {
	return JMY_PACKET_LEN_HEAD + JMY_PACKET_LEN_TYPE
}

func jmy_net_proto_disconnect_ack_pack_len() int {
	return JMY_PACKET_LEN_HEAD + JMY_PACKET_LEN_TYPE
}

func jmy_net_proto_user_id_data_pack_len() int {
	return JMY_PACKET_LEN_HEAD + JMY_PACKET_LEN_TYPE + JMY_PACKET_LEN_USER_ID + JMY_PACKET_LEN_MSG_ID
}

func jmy_net_proto_user_id_data_full_pack_len(data_len uint16) int {
	return jmy_net_proto_user_id_data_pack_len() + int(data_len)
}

func jmy_net_proto_heartbeat_pack_len() int {
	return JMY_PACKET_LEN_HEAD + JMY_PACKET_LEN_TYPE + JMY_PACKET_LEN_TIMESTAMP
}

// pack user data
func jmy_net_proto_pack_user_data_head(buf []byte, msgid uint16, data_len uint16) int {
	var pack_len int = jmy_net_proto_user_data_full_pack_len(data_len)
	if pack_len > len(buf)+int(data_len) {
		return -1
	}
	// head
	PACK_UINT16_TO_BUFF(uint16(pack_len-JMY_PACKET_LEN_HEAD), buf)
	// type
	buf[2] = byte(JMY_PACKET_USER_DATA)
	// msg id
	PACK_UINT16_TO_BUFF(msgid, buf[3:])
	return jmy_net_proto_user_data_pack_len()
}

// pack user id data
func jmy_net_proto_pack_user_id_data_head(buf []byte, user_id uint32, msg_id uint16, data_len uint16) int {
	var pack_len int = jmy_net_proto_user_id_data_full_pack_len(data_len)
	if pack_len > len(buf)+int(data_len) {
		return -1
	}
	// head
	PACK_UINT16_TO_BUFF(uint16(pack_len-JMY_PACKET_LEN_HEAD), buf)
	// type
	buf[2] = byte(JMY_PACKET_USER_ID_DATA)
	// user id
	PACK_UINT32_TO_BUFF(user_id, buf[3:])
	// msg id
	PACK_UINT16_TO_BUFF(msg_id, buf[7:])
	return jmy_net_proto_user_id_data_pack_len()
}

// pack disconnect
func jmy_net_proto_pack_disconnect(buf []byte) int {
	var pack_len int = jmy_net_proto_disconnect_pack_len()
	if pack_len > len(buf) {
		return -1
	}
	// head
	PACK_UINT16_TO_BUFF(JMY_PACKET_LEN_TYPE, buf)
	// type
	buf[2] = byte(JMY_PACKET_DISCONNECT)
	return pack_len
}

// pack disconnect ack
func jmy_net_proto_pack_disconnect_ack(buf []byte) int {
	var pack_len int = jmy_net_proto_disconnect_ack_pack_len()
	if pack_len > len(buf) {
		return -1
	}
	// head
	PACK_UINT16_TO_BUFF(JMY_PACKET_LEN_TYPE, buf)
	// type
	buf[2] = byte(JMY_PACKET_DISCONNECT_ACK)
	return pack_len
}

// pack heart beat
func jmy_net_proto_pack_heartbeat(buf []byte) int {
	var pack_len int = jmy_net_proto_heartbeat_pack_len()
	if pack_len > len(buf) {
		return -1
	}
	// head
	PACK_UINT16_TO_BUFF(uint16(pack_len-JMY_PACKET_LEN_HEAD), buf)
	// type
	buf[2] = byte(JMY_PACKET_HEARTBEAT)
	// timestamp
	var t uint32 = uint32(time.Now().Unix())
	PACK_UINT32_TO_BUFF(t, buf[3:])
	return pack_len
}
