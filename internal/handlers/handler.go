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

func HandleTransfer(w http.ResponseWriter, r *http.Request) {
	// читаем тело запроса - сегмент
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// парсим сегмент в структуру
	segment := utils.Segment{}
	if err = json.Unmarshal(body, &segment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// пишем сегмент в Kafka
	if err = kafka.WriteToKafka(segment); err != nil {
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
