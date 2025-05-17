package consts

import "time"

const (
	SegmentLostError = "lost"
	KafkaReadPeriod  = 2 * time.Second
	CodeUrl          = "http://mainc:8020/code" // адрес канального уровня
	SegmentSize      = 230000
	ReceiveUrl = "http://server:8011/receive" // адрес websocket-сервера прикладного уровня
	KafkaAddr  = "kafka:9092"
	KafkaTopic = "segments"
)
