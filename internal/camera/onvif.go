package camera

import (
	"context"
	"fmt"
	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	goonvif "github.com/neirolis/onvif-go"
	"github.com/neirolis/onvif-go/event"
	"github.com/neirolis/onvif-go/media"
	"github.com/neirolis/onvif-go/xsd"
	"github.com/neirolis/onvif-go/xsd/onvif"
	"io"
	"net/http"
	"net/url"
	"time"
)

type onvifController struct {
	u              *url.URL
	dev            *goonvif.Device
	ctx            context.Context
	snapshotUrl    *url.URL
	streamUrl      string
	eventsEndpoint string
}

type pullMessagesResponse struct {
	CurrentTime         event.CurrentTime
	TerminationTime     event.TerminationTime
	NotificationMessage []event.NotificationMessage
}

const maxMessageLimit = 100

func newOnvifController(ctx context.Context, u *url.URL) *onvifController {
	return &onvifController{u: u, ctx: ctx}
}

func (c *onvifController) GetEvents() ([]*iva.Event, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	if err := c.subscribe(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
	defer cancel()

	var response pullMessagesResponse
	request := event.PullMessages{
		MessageLimit: maxMessageLimit,
		Timeout:      xsd.Duration(fmt.Sprintf("PT%dS", maxNetworkTimeout/(2*time.Second))),
	}
	if err := c.dev.CreateRequest(request).WithEndpoint(c.eventsEndpoint).WithContext(ctx).Do().Unmarshal(&response); err != nil {
		c.clearCache()
		return nil, fmt.Errorf("method PullMessages failed: %w", err)
	}

	var events []*iva.Event
	for _, ev := range response.NotificationMessage {
		for _, msg := range ev.Message.Messages {
			conv := convertEvent(string(ev.Topic.TopicKinds), msg)
			if conv != nil {
				events = append(events, conv)
			}
		}
	}
	return events, nil
}

func (c *onvifController) GetSnapshot(profileToken string) ([]byte, error) {
	var err error
	if err = c.connect(); err != nil {
		return nil, err
	}
	if c.snapshotUrl == nil {
		ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
		defer cancel()

		request := media.GetSnapshotUri{ProfileToken: onvif.ReferenceToken(profileToken)}
		var response media.GetSnapshotUriResponse
		if err := c.dev.CreateRequest(request).WithContext(ctx).Do().Unmarshal(&response); err != nil {
			return nil, fmt.Errorf("method GetSnapshotUri failed: %w", err)
		}
		c.snapshotUrl, err = url.Parse(string(response.MediaUri.Uri))
		if c.snapshotUrl.Host == "" {
			c.snapshotUrl.Host = c.u.Host
			c.snapshotUrl.Scheme = c.u.Scheme
		}
		c.snapshotUrl.User = c.u.User
	}

	httpClient := http.Client{
		Timeout: maxNetworkTimeout,
	}

	resp, err := httpClient.Get(c.snapshotUrl.String())
	if err != nil {
		c.clearCache()
		return nil, fmt.Errorf("screenshot request failed: %w", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("download screenshot failed: %w", err)
	}
	return result, nil
}

func (c *onvifController) GetStreamUri(profileToken string) (string, error) {
	var err error
	if err = c.connect(); err != nil {
		return "", err
	}
	if c.streamUrl != "" {
		return c.streamUrl, nil
	}

	ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
	defer cancel()

	request := media.GetStreamUri{ProfileToken: onvif.ReferenceToken(profileToken)}
	var response media.GetStreamUriResponse
	if err := c.dev.CreateRequest(request).WithContext(ctx).Do().Unmarshal(&response); err != nil {
		return "", fmt.Errorf("method GetStreamUri failed: %w", err)
	}
	c.streamUrl = string(response.MediaUri.Uri)
	return c.streamUrl, nil
}

func (c *onvifController) connect() error {
	if c.dev != nil {
		return nil
	}

	params := goonvif.DeviceParams{
		// TODO: адекватное определение endpoint
		Xaddr:    fmt.Sprintf("%s:80", c.u.Host),
		Username: c.u.User.Username(),
	}
	if password, set := c.u.User.Password(); set {
		params.Password = password
	}
	dev := goonvif.NewDevice(params)

	ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
	defer cancel()

	_, err := dev.InspectWithCtx(ctx)
	if err != nil {
		return fmt.Errorf("camera is unaccessible: %w", err)
	}

	c.dev = dev
	return nil
}

func (c *onvifController) subscribe() error {
	if c.eventsEndpoint != "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
	defer cancel()

	var response event.CreatePullPointSubscriptionResponse
	if err := c.dev.CreateRequest(event.CreatePullPointSubscription{}).WithContext(ctx).Do().Unmarshal(&response); err != nil {
		return fmt.Errorf("method CreatePullPointSubscription failed: %w", err)
	}
	c.eventsEndpoint = string(response.SubscriptionReference.Address)
	return nil
}

func (c *onvifController) clearCache() {
	c.dev = nil
	c.streamUrl = ""
	c.eventsEndpoint = ""
	c.snapshotUrl = nil
}
