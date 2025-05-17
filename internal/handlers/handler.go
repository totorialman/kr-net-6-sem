package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/totorialman/kr-net-6-sem/internal/kafka"
	"github.com/totorialman/kr-net-6-sem/internal/utils"
	"github.com/totorialman/kr-net-6-sem/internal/consts"
)

type ErrorResponse struct {
    Error string `json:"error"`
}

func HandleTransfer(w http.ResponseWriter, r *http.Request) {
    // Читаем тело запроса
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Попробуем сначала прочитать как ошибку
    var errResp ErrorResponse
    if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
        // Это сообщение об ошибке

        // Парсим как Segment, чтобы получить username и timestamp
        var segment utils.Segment
        if unmarshalErr := json.Unmarshal(body, &segment); unmarshalErr == nil {
            // Отправляем кастомный ответ с error_flag = true
            response := utils.ReceiveRequest{
                Username: segment.Username,
                SendTime: segment.SendTime,
                Error:    errResp.Error,
                Text:     "", // нет текста из-за ошибки
            }

            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(response)
        } else {
            // Не удалось распарсить даже в Segment
            http.Error(w, "Invalid request body", http.StatusBadRequest)
        }
        return
    }

    // Если это не ошибка — продолжаем как обычно
    var segment utils.Segment
    if err := json.Unmarshal(body, &segment); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := kafka.WriteToKafka(segment); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func HandleSend(w http.ResponseWriter, r *http.Request) {
	// читаем тело запроса - сообщение
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// парсим сообщение в структуру
	message := utils.SendRequest{}
	if err = json.Unmarshal(body, &message); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// сразу отвечаем прикладному уровню 200 ОК - мы приняли работу
	w.WriteHeader(http.StatusOK)

	// разбиваем текст сообщения на сегменты
	segments := utils.SplitMessage(message.Text, consts.SegmentSize)
	total := len(segments)

	// в цикле отправляем сегменты на канальный уровень
	for i, segment := range segments {
		payload := utils.Segment{
			SegmentNumber:  i + 1,
			TotalSegments:  total,
			Username:       message.Username,
			SendTime:       message.SendTime,
			SegmentPayload: segment,
		}
		go utils.SendSegment(payload) // запускаем горутину с отправкой на канальный уровень, не будем дожидаться результата ее выполнения
		fmt.Printf("sent segment: %+v\n", payload)
	}
}
