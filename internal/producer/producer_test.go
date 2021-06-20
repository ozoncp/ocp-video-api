package producer_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ocp-video-api/internal/mock"
	"ocp-video-api/internal/producer"
)

var _ = Describe("Producer", func() {
	var (
		err error

		ctrl *gomock.Controller

		mockSender *mock.MockSender

		p producer.Producer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockSender = mock.NewMockSender(ctrl)

		p, _ = producer.NewProducer(1, mockSender)

		mockSender.EXPECT().Init().Return(nil)
		err = p.Init()
	})

	AfterEach(func() {
		p.Close()
		ctrl.Finish()
	})

	Context("test sending", func() {
		BeforeEach(func() {
			e := producer.Event{
				Type: producer.EventTypCreated,
				Data: map[string]interface{}{"data": "created"},
			}

			mockSender.EXPECT().Send(producer.Event{
				Type: producer.EventTypCreated,
				Data: map[string]interface{}{"data": "created"},
			}).Return(nil)
			mockSender.EXPECT().Close().Return(nil)

			err = p.SendEvent(e)
		})

		It("succesfully sended", func() {
			Expect(err).Should(BeNil())
		})
	})

	Context("test sending to closing is failing", func() {
		BeforeEach(func() {
			e := producer.Event{
				Type: producer.EventTypCreated,
				Data: map[string]interface{}{"data": "created"},
			}

			mockSender.EXPECT().Close().Return(nil)

			p.Close()
			err = p.SendEvent(e)
		})

		It("succesfully sended", func() {
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("test double close do not hangs", func() {
		BeforeEach(func() {
			e := producer.Event{
				Type: producer.EventTypCreated,
				Data: map[string]interface{}{"data": "created"},
			}

			mockSender.EXPECT().Send(producer.Event{
				Type: producer.EventTypCreated,
				Data: map[string]interface{}{"data": "created"},
			}).Return(nil)
			mockSender.EXPECT().Close().Return(nil).Times(1)

			err = p.SendEvent(e)
			p.Close()
			p.Close()
		})

		It("succesfully sended", func() {
			Expect(err).Should(BeNil())
		})
	})
})
