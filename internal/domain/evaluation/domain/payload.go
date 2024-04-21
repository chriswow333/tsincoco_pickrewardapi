package domain

import (
	"errors"
	"math"
	"pickrewardapi/internal/domain/evaluation/domain/event"
	commonM "pickrewardapi/internal/shared/common/model"

	log "github.com/sirupsen/logrus"
)

type PayloadOperator int32

const (
	MaxAndPayloadOperator PayloadOperator = iota
	MaxOrPayloadOperator
	XorPayloadOperator
	AddPayloadOperator
)

type PayloadType int32

const (
	SelfPayloadType PayloadType = iota
	ContainerPayloadType
)

type Payload struct {
	ID              string
	PayloadOperator PayloadOperator
	PayloadType     PayloadType

	/**
	*		如果payload type 是 payload
	*		則 feedback 應為 nil
	*		因為要用operator 去算 nested payloads 的結果
	**/
	Feedback  *Feedback
	Payloads  []*Payload
	Container *Container
}

// 計算每個 payload 細節
func (p *Payload) Judge(e *commonM.Event) (*event.PayloadEventResult, error) {
	logPos := "[evaluation.domain][Payload.Judge]"

	switch p.PayloadType {
	case SelfPayloadType:
		return p.judgePayload(e)
	case ContainerPayloadType:
		return p.judgeContainer(e)
	}

	log.WithFields(log.Fields{
		"pos":       logPos,
		"PayloadID": p.ID,
		"eventID":   e.ID,
	}).Errorf("Not found PayloadType, %d", p.PayloadType)

	return nil, errors.New("Not found PayloadType")
}

func (p *Payload) judgeContainer(e *commonM.Event) (*event.PayloadEventResult, error) {
	logPos := "[evaluation.domain][Payload.judgeContainer]"

	containerEvent := p.Container.Satisfy(e)

	var feedbackEvent *event.FeedbackEventResult

	if containerEvent.Pass {
		var err error
		feedbackEvent, err = p.calculateFeedbackResult(e)

		if err != nil {
			log.WithFields(log.Fields{
				"pos":       logPos,
				"PayloadID": p.ID,
				"eventID":   e.ID,
			}).Error("p.calculateFeedbackResult failed", err)
			return nil, err
		}

	} else {
		feedbackEvent = &event.FeedbackEventResult{
			FeedbackID:                p.Feedback.FeedbackID,
			Cost:                      e.Cost,
			FeedbackEventResultStatus: event.GetNone,
		}
	}

	if p.Feedback != nil {
		feedbackEvent.CalculateType = int32(p.Feedback.CalculateType)
	}

	return &event.PayloadEventResult{
		ID:                   p.ID,
		Pass:                 containerEvent.Pass,
		FeedbackEventResult:  feedbackEvent,
		ContainerEventResult: containerEvent,
	}, nil
}

func (p *Payload) judgePayload(e *commonM.Event) (*event.PayloadEventResult, error) {
	logPos := "[evaluation.domain][Payload.judgePayload]"

	payloadEventResults := []*event.PayloadEventResult{}
	for _, payload := range p.Payloads {
		payloadEventResult, err := payload.Judge(e)
		if err != nil {
			log.WithFields(log.Fields{
				"pos":       logPos,
				"PayloadID": p.ID,
				"eventID":   e.ID,
			}).Errorf("payload.Judge failed, %s", err)
			return nil, err
		}
		payloadEventResults = append(payloadEventResults, payloadEventResult)
	}

	payloadEvent, err := p.operate(e, payloadEventResults)
	if err != nil {
		log.WithFields(log.Fields{
			"pos":       logPos,
			"PayloadID": p.ID,
			"eventID":   e.ID,
		}).Errorf("p.operate failed, %s", err)
		return nil, err
	}

	payloadEvent.PayloadEventResults = payloadEventResults

	return payloadEvent, nil
}

func (p *Payload) operate(e *commonM.Event, payloadEvents []*event.PayloadEventResult) (*event.PayloadEventResult, error) {
	logPos := "[reward.domain][Payload.calculatePayloadEvents]"

	logFields := log.Fields{
		"pos":       logPos,
		"PayloadID": p.ID,
	}
	log.WithFields(logFields).Info("Start Payload calculatePayloadEvents, PayloadOperator type :", p.PayloadOperator)

	switch p.PayloadOperator {
	case MaxAndPayloadOperator:
		pass := true
		maxFeedbackEvent := &event.FeedbackEventResult{
			GetReturn: 0,
		}

		for _, e := range payloadEvents {
			pass = pass && e.Pass
			if pass {
				maxFeedbackEvent = max(maxFeedbackEvent, e.FeedbackEventResult)
			} else {
				break
			}
		}

		if !pass {
			maxFeedbackEvent.GetReturn = 0
		}

		maxFeedbackEvent.Cost = e.Cost

		maxPayloadEvent := &event.PayloadEventResult{
			ID:                  p.ID,
			Pass:                pass,
			PayloadEventResults: payloadEvents,
			FeedbackEventResult: maxFeedbackEvent,
		}

		return maxPayloadEvent, nil

	case MaxOrPayloadOperator:
		pass := false
		maxFeedbackEvent := &event.FeedbackEventResult{
			GetReturn: 0,
		}

		for _, e := range payloadEvents {
			pass = pass || e.Pass
			if pass {
				maxFeedbackEvent = max(maxFeedbackEvent, e.FeedbackEventResult)
			}
		}

		maxFeedbackEvent.Cost = e.Cost
		maxPayloadEvent := &event.PayloadEventResult{
			ID:                  p.ID,
			Pass:                pass,
			PayloadEventResults: payloadEvents,
			FeedbackEventResult: maxFeedbackEvent,
		}

		return maxPayloadEvent, nil

	case XorPayloadOperator:

		// If inner have more than two pass, failed.

		xorPass := false
		innerXorPayloadEvent := &event.PayloadEventResult{}
		for _, pe := range payloadEvents {
			if xorPass && pe.Pass {
				xorPass = false
				break
			}
			xorPass = xorPass || pe.Pass
			if xorPass {
				innerXorPayloadEvent = pe
			}
		}

		xorPayloadEvent := &event.PayloadEventResult{
			ID:                  p.ID,
			PayloadEventResults: payloadEvents,
		}

		xorPayloadEvent.Pass = xorPass
		if xorPass {
			xorPayloadEvent.FeedbackEventResult = innerXorPayloadEvent.FeedbackEventResult
		} else {
			xorPayloadEvent.FeedbackEventResult = &event.FeedbackEventResult{
				GetReturn: 0,
				Cost:      e.Cost,
			}
		}
		return xorPayloadEvent, nil

	case AddPayloadOperator:

		addFeedbackEvent := &event.FeedbackEventResult{
			GetReturn:                 0,
			GetPercentage:             0,
			FeedbackEventResultStatus: event.GetNone,
			Cost:                      e.Cost,
		}

		eventStatus := -1

		for _, pe := range payloadEvents {
			if pe.Pass {
				addFeedbackEvent.GetReturn += pe.FeedbackEventResult.GetReturn
				addFeedbackEvent.GetPercentage += pe.FeedbackEventResult.GetPercentage
			}

			if eventStatus == -1 {
				eventStatus = int(pe.FeedbackEventResult.FeedbackEventResultStatus)
			} else if eventStatus == 1 {
				continue
			} else if eventStatus > int(pe.FeedbackEventResult.FeedbackEventResultStatus) ||
				eventStatus < int(pe.FeedbackEventResult.FeedbackEventResultStatus) {
				eventStatus = 1
			}

		}

		addFeedbackEvent.FeedbackEventResultStatus = event.FeedbackEventResultStatus(eventStatus)

		addPayloadEvent := &event.PayloadEventResult{
			ID:                  p.ID,
			PayloadEventResults: payloadEvents,
			FeedbackEventResult: addFeedbackEvent,
		}

		if addFeedbackEvent.GetReturn > 0 {
			addPayloadEvent.Pass = true
			return addPayloadEvent, nil
		} else {
			addPayloadEvent.Pass = false
			return addPayloadEvent, nil
		}

	}

	log.WithFields(logFields).Errorf("Not found PayloadOperator, %d", p.PayloadOperator)
	return nil, errors.New("Not found PayloadOperator")
}

func max(a, b *event.FeedbackEventResult) *event.FeedbackEventResult {
	if a.GetReturn > b.GetReturn {
		return a
	} else {
		return b
	}
}

func (p *Payload) calculateFeedbackResult(e *commonM.Event) (*event.FeedbackEventResult, error) {

	logPos := "[reward.domain][Payload.calculateFeedbackResult]"

	if p.Feedback == nil {
		log.WithFields(log.Fields{
			"pos":       logPos,
			"PayloadID": p.ID,
			"eventID":   e.ID,
		}).Errorf("p.Feedback is nil")
		return nil, errors.New("p.Feedback is nil")
	}

	if e.Cost < p.Feedback.MinCost {
		return &event.FeedbackEventResult{
			CalculateType:             int32(p.Feedback.CalculateType),
			Cost:                      e.Cost,
			FeedbackEventResultStatus: event.GetNone,
			GetPercentage:             0,
		}, nil
	}

	switch p.Feedback.CalculateType {

	case Multiply:
		feedbackEvent := multiply(p.Feedback, e)
		feedbackEvent.Cost = e.Cost
		feedbackEvent.GetPercentage = p.Feedback.Percentage
		feedbackEvent.CalculateType = int32(p.Feedback.CalculateType)
		return feedbackEvent, nil

	case Fixed:
		feedbackEvent := fixed(p.Feedback)
		feedbackEvent.Cost = e.Cost
		feedbackEvent.CalculateType = int32(p.Feedback.CalculateType)
		return feedbackEvent, nil

	case Area:
		feedbackEvent := area(p.Feedback, e)
		feedbackEvent.Cost = e.Cost
		feedbackEvent.CalculateType = int32(p.Feedback.CalculateType)

		return feedbackEvent, nil
	}

	log.WithFields(log.Fields{
		"pos":       logPos,
		"PayloadID": p.ID,
		"eventID":   e.ID,
	}).Error("Cannot find CalculateType", p.Feedback.CalculateType)
	return nil, errors.New("Cannot find CalculateType")
}

func fixed(feedback *Feedback) *event.FeedbackEventResult {

	return &event.FeedbackEventResult{
		GetReturn:                 float64(feedback.Fixed),
		FeedbackEventResultStatus: event.GetAll,
	}
}

func multiply(feedback *Feedback, e *commonM.Event) *event.FeedbackEventResult {

	feedbackEvent := &event.FeedbackEventResult{
		GetPercentage: feedback.Percentage,
	}

	getReturn := math.Round(feedback.Percentage * float64(e.Cost))
	feedbackEvent.FeedbackEventResultStatus = event.GetAll

	if feedback.ReturnMax != 0 && getReturn > feedback.ReturnMax {
		getReturn = feedback.ReturnMax
		feedbackEvent.FeedbackEventResultStatus = event.GetSome
	}

	feedbackEvent.GetReturn = getReturn
	return feedbackEvent

}

func area(feedback *Feedback, e *commonM.Event) *event.FeedbackEventResult {
	feedbackEvent := &event.FeedbackEventResult{
		GetPercentage: feedback.Percentage,
	}

	if e.Cost <= feedback.MinCost {
		feedbackEvent.GetReturn = 0
		feedbackEvent.FeedbackEventResultStatus = event.GetNone
	} else {
		feedbackEvent.GetReturn = math.Round(feedback.Percentage * float64((e.Cost - feedback.MinCost)))
		feedbackEvent.FeedbackEventResultStatus = event.GetAll
	}

	return feedbackEvent
}
