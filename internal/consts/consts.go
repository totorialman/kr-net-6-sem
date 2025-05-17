package consts

import "time"

const (
	SegmentLostError = "lost"
	KafkaReadPeriod  = 2 * time.Second
	CodeUrl          = "mainc:8020/code" // адрес канального уровня
	SegmentSize      = 2300
	ReceiveUrl = "mainc:8020/receive" // адрес websocket-сервера прикладного уровня
	KafkaAddr  = "kafka:9092"
	KafkaTopic = "segments"
)
