package consts

import "time"

const (
	SegmentLostError = "lost"
	KafkaReadPeriod  = 2 * time.Second
	CodeUrl          = "http://192.168.123.120:8000/code" // адрес канального уровня
	SegmentSize      = 100
	ReceiveUrl = "http://192.168.123.140:3000/receive" // адрес websocket-сервера прикладного уровня
	KafkaAddr  = "kafka:29092"
	KafkaTopic = "segments"
)
