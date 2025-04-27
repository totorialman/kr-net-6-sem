package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/totorialman/kr-net-6-sem/internal/utils"
	"github.com/totorialman/kr-net-6-sem/internal/consts"
)

type Storage map[time.Time]utils.Message

var storage = Storage{}

func addMessage(segment utils.Segment) {
	storage[segment.SendTime] = utils.Message{
		Received: 0,
		Total:    segment.TotalSegments,
		Last:     time.Now().UTC(),
		Username: segment.Username,
		Segments: make([]string, segment.TotalSegments), // заранее выделяем память, это важно!
	}
}

func AddSegment(segment utils.Segment) {
	// используем мьютекс, чтобы избежать конкуретного доступа к хранилищу
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	// если это первый сегмент сообщения, создаем пустое сообщение
	sendTime := segment.SendTime
	_, found := storage[sendTime]
	if !found {
		addMessage(segment)
	}

	// добавляем в сообщение информацию о сегменте
	message, _ := storage[sendTime]
	message.Received++
	message.Last = time.Now().UTC()
	message.Segments[segment.SegmentNumber-1] = segment.SegmentPayload // сохраняем правильный порядок сегментов
	storage[sendTime] = message
}

func getMessageText(sendTime time.Time) string {
	result := ""
	message, _ := storage[sendTime]
	for _, segment := range message.Segments {
		result += segment
	}
	return result
}



// структура тела запроса на прикладной уровень

type sendFunc func(body utils.ReceiveRequest)

func ScanStorage(sender sendFunc) {
	mu := &sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	payload := utils.ReceiveRequest{}
	for sendTime, message := range storage {
		if message.Received == message.Total { // если пришли все сегменты
			payload = utils.ReceiveRequest{
				Username: message.Username,
				Text:     getMessageText(sendTime), // склейка сообщения
				SendTime: sendTime,
				Error:    "",
			}
			fmt.Printf("sent message: %+v\n", payload)
			go sender(payload)        // запускаем горутину с отправкой на прикладной уровень, не будем дожидаться результата ее выполнения
			delete(storage, sendTime) // не забываем удалять
		} else if time.Since(message.Last) > consts.KafkaReadPeriod+time.Second { // если канальный уровень потерял сегмент
			payload = utils.ReceiveRequest{
				Username: message.Username,
				Text:     "",
				SendTime: sendTime,
				Error:    consts.SegmentLostError, // ошибка
			}
			fmt.Printf("sent error: %+v\n", payload)
			go sender(payload)        // запускаем горутину с отправкой на прикладной уровень, не будем дожидаться результата ее выполнения
			delete(storage, sendTime) // не забываем удалять
		}
	}
}
