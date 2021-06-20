// TODO: вытянуть общий каркас с internal/saver вверх по иерархии
// возможно стоит подумать о общем фреймворке-планировщике событий с пулом горутин в которые оффлоадить таски

package producer

import (
	"errors"
	"github.com/rs/zerolog/log"
	"sync/atomic"
)

type EventTyp = uint64

const (
	EventTypDefault EventTyp = iota
	EventTypCreated
	EventTypUpdated
	EventTypRemoved
)

type Event struct {
	Type EventTyp
	Data map[string]interface{}
}

type Producer interface {
	Init() error
	SendEvent(Event) error
	Close()
}

type producer struct {
	closing int32
	events  chan Event
	// в saver идея отменять через контекст похоже не идиоматична для Go,
	// здесь отмена через запись в канал в Close
	close  chan struct{}
	done   chan struct{}
	sender Sender
}

func NewProducer(capacity uint, sender Sender) (Producer, error) {
	return &producer{
		closing: 0,
		events:  make(chan Event, capacity),
		close:   make(chan struct{}),
		done:    make(chan struct{}),
		sender:  sender,
	}, nil
}

func (p *producer) Init() error {
	err := p.sender.Init()
	if err != nil {
		return err
	}
	go p.sendLoop()

	return nil
}

func (p *producer) SendEvent(event Event) error {
	if atomic.LoadInt32(&p.closing) != 0 {
		return errors.New("producer is closing no new events can be sended")
	}

	if len(p.events) == cap(p.events) {
		// TODO: возможно если буфер полон есть проблема и нужно предпринять какие-то действия
		// в случае если буфер канала полон - дропаем наиболее старую запись
		// чтобы не заблокироваться тут, если буфер полон
		<-p.events
		log.Print("Dropped event due to buffered chan full")
	}
	p.events <- event

	return nil
}

func (p *producer) Close() {
	if atomic.LoadInt32(&p.closing) != 0 {
		return
	}

	p.close <- struct{}{}
	<-p.done
}

func (p *producer) sendLoop() {
	for {
		select {
		case event := <-p.events:
			err := p.sender.Send(event)
			if err != nil {
				// TODO: перепланировать позже? проанализировать и устранить причину (например переподключением?)
				log.Print("Error sending event", err)
			}
		case <-p.close:
			atomic.AddInt32(&p.closing, 1)
			for len(p.events) > 0 {
				err := p.sender.Send(<-p.events)
				if err != nil {
					// TODO: перепланировать позже? проанализировать и устранить причину (например переподключением?)
					log.Print("Error sending event during drain", err)
				}
			}

			// записи в этот канал быть не должно, т.к. SendEvent возвращает error
			// если будет паника из-за записи в закрытый канал, то это логическая ошибка
			// и это лучше если это событие просто потеряется
			close(p.events)

			err := p.sender.Close()
			if err != nil {
				// TODO: перепланировать позже? проанализировать и устранить причину (например переподключением?)
				log.Print("Failed to close sarama.Producer", err)
			}
			p.done <- struct{}{}
			return
		}
	}
}
