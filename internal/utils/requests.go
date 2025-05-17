package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/totorialman/kr-net-6-sem/internal/consts"
)

type Segment struct {
	SegmentNumber  int       `json:"sequence_number"`
	TotalSegments  int       `json:"total_parts"`
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
	// Логируем начало отправки
	fmt.Printf("Sending segment: %+v\n", body)

	// Маршалим тело запроса
	reqBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", consts.CodeUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус-код
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Successfully sent segment to %s (status: %d)\n", consts.CodeUrl, resp.StatusCode)
	} else {
		fmt.Printf("Failed to send segment. Server responded with status: %d\n", resp.StatusCode)
	}
}
func SendReceiveRequest(body ReceiveRequest) {
	// Логируем начало отправки
	log.Printf("Sending request to %s: %+v", consts.ReceiveUrl, body)

	// Маршалим в JSON
	reqBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		return
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", consts.ReceiveUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	// Логируем успешный ответ
	log.Printf("Request sent successfully. Status code: %d", resp.StatusCode)
}
