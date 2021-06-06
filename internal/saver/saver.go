package saver

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"ocp-video-api/internal/flusher"
	"ocp-video-api/internal/models"
)

type OverflowPolicy = uint8

const (
	OverflowDropAll = OverflowPolicy(iota)
	OverflowDropFirst
)

func DropFirstPolicy(s *saver) {
	s.policy = OverflowDropFirst
}

type Saver interface {
	Init(ctx context.Context)
	Save(ctx context.Context, v models.Video) error
	Close()
}

func New(
	capacity uint,
	period time.Duration,
	flusher flusher.Flusher,
	cfgrers ...func(*saver),
) Saver {
	vs := make(chan models.Video, capacity)
	done := make(chan struct{})

	rv := &saver{
		closing: 0,
		policy:  OverflowDropAll,
		videos:  vs,
		done:    done,
		period:  period,
		flusher: flusher,
		ticker:  nil,
	}

	// позволяет позднее добавлять в пакет функции-конфигураторы, вызываемые создателем, предполагается, что
	// внутри пакета есть некая политика по умолчанию
	for _, cfgrer := range cfgrers {
		cfgrer(rv)
	}

	return rv
}

type saver struct {
	closing int32
	policy  OverflowPolicy
	videos  chan models.Video
	done    chan struct{}
	period  time.Duration
	flusher flusher.Flusher
	ticker  *time.Ticker
}

func (s *saver) Init(ctx context.Context) {
	s.ticker = time.NewTicker(s.period)
	go s.flushLoop(ctx)
}

func (s *saver) Save(ctx context.Context, v models.Video) error {
	if atomic.LoadInt32(&s.closing) != 0 {
		return errors.New("saver is closing no new videos can be saved")
	}

	if len(s.videos) == cap(s.videos) {
		switch s.policy {
		case OverflowDropAll:
			//TODO: каждый раз при чтении из канала лочится мьютекс, найти способ единожды залочить мьютекс
			// и вычитать cap(s.videos) из канала
			for drainCnt := cap(s.videos); drainCnt > 0; drainCnt-- {
				<-s.videos
			}
		case OverflowDropFirst:
			<-s.videos
		}
	}
	s.videos <- v
	return nil
}

func handleFlushError(vs []models.Video, err error) {
	//TODO: add log
	fmt.Printf("Something must be done with: %v, got error %v on flush, this videos are not saved\n", vs, err)
}

const PanicMultiper = 3

func (s *saver) flushLoop(ctx context.Context) {
	buff := make([]models.Video, 0, cap(s.videos))

	for {
		select {
		case v := <-s.videos:
			buff = append(buff, v)

		case <-ctx.Done():
			s.ticker.Stop()

			// добавление видео осуществляется асинхронно через канал s.videos из других контестов выполнения,
			// должен наступить момент когда данный мини циклсобытий завершается, далее будут вычитаны все уже
			// добавленные ранее видео для сохранения
			atomic.AddInt32(&s.closing, 1)
			for len(s.videos) > 0 {
				buff = append(buff, <-s.videos)
			}
			// пусть лучше будет паника если есть логическая ошибка и будет запись в закрытый канал
			// из которого уже никогда не будет чтения данных, т.к. тут через s.closing планируется не
			// допустить записей в канал, в блоке выше из этого канала считывается все что есть
			close(s.videos)

			if rest, err := s.flusher.Flush(buff); err != nil {
				handleFlushError(rest, err)
			}
			s.done <- struct{}{}
			return

		case <-s.ticker.C:
			if len(buff) == 0 {
				continue
			}

			if rest, err := s.flusher.Flush(buff); err != nil {
				if len(buff) > PanicMultiper*cap(s.videos) && &rest[0] == &buff[0] {
					//TODO: log error
					panic("flusher critical error")
				}
				//TODO: log error
				fmt.Printf("Got error (%v) on periodic flush, unsaved (%v) will be retried to save on periodic tick\n", err, rest)
				buff = buff[:0]
				buff = append(buff, rest...)
			} else {
				buff = buff[:0]
			}
		}
	}
}

func (s *saver) Close() {
	<-s.done
}
