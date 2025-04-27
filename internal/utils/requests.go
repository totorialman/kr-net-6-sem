package utils

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"time"
	"github.com/totorialman/kr-net-6-sem/internal/consts"
)

type Segment struct {
	SegmentNumber  int       `json:"sequence_number"`
	TotalSegments  int       `json:"total_part"`
	Username       string    `json:"username"`
	SendTime       time.Time `json:"timestamp"`
	SegmentPayload string    `json:"message_part"`
}


type SendRequest struct {
    Id       int       `json:"id,omitempty"`
    Username string    `json:"username"`
    Text     string    `json:"message"`
    SendTime time.Time `json:"timestamp"`
}

type Message struct {
	Received int
	Total    int
	Last     time.Time
	Username string
	Segments []string
}

type ReceiveRequest struct {
    Username string    `json:"username"`
    Text     string    `json:"message"`
    SendTime time.Time `json:"timestamp"`
    Error    string    `json:"error_flag"`
}

func SplitMessage(payload string, segmentSize int) []string {
	result := make([]string, 0)

	length := len(payload) // длина в байтах
	segmentCount := int(math.Ceil(float64(length) / float64(segmentSize)))

	for i := 0; i < segmentCount; i++ {
		result = append(result, payload[i*segmentSize:min((i+1)*segmentSize, length)]) // срез делается также по байтам
	}

	return result
}

func SendSegment(body Segment) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", consts.CodeUrl, bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

func SendReceiveRequest(body ReceiveRequest) {
   reqBody, _ := json.Marshal(body)
   
   req, _ := http.NewRequest("POST", consts.ReceiveUrl, bytes.NewBuffer(reqBody))
   req.Header.Add("Content-Type", "application/json")
   
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
     return
   }
   
   defer resp.Body.Close()
}