package camera

import (
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	"github.com/RacoonMediaServer/rms-packages/pkg/events"
	"github.com/neirolis/onvif-go/event"
	"time"
)

func simpleItemsToDict(si []event.SimpleItem) map[string]string {
	result := make(map[string]string)
	for _, item := range si {
		result[string(item.Name)] = string(item.Value)
	}
	return result
}

func convertIntervalEvent(ev *event.Message, field string, kind events.Alert_Kind) *iva.Event {
	isInitialized := ev.PropertyOperation == "Initialized"
	dict := simpleItemsToDict(ev.Data.SimpleItem)
	isActive := dict[field]
	interval := iva.End
	if isActive == "true" {
		interval = iva.Begin
	} else if isInitialized {
		return nil
	}

	return &iva.Event{
		Kind:      kind,
		Interval:  interval,
		Timestamp: time.Now(),
	}
}

func convertSingleEvent(ev *event.Message, kind events.Alert_Kind) *iva.Event {
	return &iva.Event{
		Kind:      kind,
		Interval:  iva.Once,
		Timestamp: time.Now(),
	}
}

func convertEvent(topic string, ev event.Message) *iva.Event {
	switch topic {
	case "tns1:RuleEngine/CellMotionDetector/Motion":
		return convertIntervalEvent(&ev, "IsMotion", events.Alert_MotionDetected)
	case "tns1:RuleEngine/TamperDetector/Tamper":
		return convertIntervalEvent(&ev, "IsTamper", events.Alert_TamperDetected)
	case "tns1:RuleEngine/FieldDetector/ObjectsInside":
		return convertIntervalEvent(&ev, "IsInside", events.Alert_IntrusionDetected)
	case "tns1:RuleEngine/LineDetector/Crossed":
		return convertSingleEvent(&ev, events.Alert_CrossLineDetected)
	default:
		return nil
	}
}
