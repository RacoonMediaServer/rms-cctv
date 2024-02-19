package camera

import (
	"context"
	"fmt"
	dac "github.com/xinsnake/go-http-digest-auth-client"
	"go-micro.dev/v4/logger"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/RacoonMediaServer/rms-cctv/internal/iva"
	goonvif "github.com/neirolis/onvif-go"
	"github.com/neirolis/onvif-go/event"
	"github.com/neirolis/onvif-go/media"
	"github.com/neirolis/onvif-go/xsd"
	"github.com/neirolis/onvif-go/xsd/onvif"
)

type onvifController struct {
	u              *url.URL
	l              logger.Logger
	dev            *goonvif.Device
	ctx            context.Context
	snapshotUrl    *url.URL
	eventsEndpoint string
}

type pullMessagesResponse struct {
	CurrentTime         event.CurrentTime
	TerminationTime     event.TerminationTime
	NotificationMessage []event.NotificationMessage
}

const maxMessageLimit = 100

func newOnvifController(ctx context.Context, u *url.URL) *onvifController {
	return &onvifController{
		u:   u,
		ctx: ctx,
		l: logger.Fields(map[string]interface{}{
			"device": u.Redacted(),
			"from":   "onvif",
		}),
	}
}

func (c *onvifController) GetDeviceName() (string, error) {
	if err := c.connect(); err != nil {
		return "", err
	}
	info := c.dev.GetDeviceInfo()
	return info.Name(), nil
}

func (c *onvifController) GetEvents() ([]*iva.Event, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}
	if err := c.subscribe(); err != nil {
		return nil, err
	}

	c.l.Logf(logger.TraceLevel, "GetEvents...")

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

	c.l.Logf(logger.TraceLevel, "GetEvents READY")

	var events []*iva.Event
	for _, ev := range response.NotificationMessage {
		for _, msg := range ev.Message.Messages {
			c.l.Logf(logger.DebugLevel, "Got event topic=%s, operation=%s, source=%v, data=%v", ev.Topic.TopicKinds, msg.PropertyOperation, msg.Source.SimpleItem, msg.Data.SimpleItem)
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

	c.l.Logf(logger.DebugLevel, "GetSnapshot...")

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

	if pw, ok := c.snapshotUrl.User.Password(); ok {
		transport := dac.NewTransport(c.u.User.Username(), pw)
		httpClient.Transport = &transport
	}

	resp, err := httpClient.Get(c.snapshotUrl.String())
	if err != nil {
		c.clearCache()
		return nil, fmt.Errorf("screenshot request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request snapshot failed: invalid code %d != 200", resp.StatusCode)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("download screenshot failed: %w", err)
	}

	c.l.Logf(logger.DebugLevel, "GetSnapshot READY")

	return result, nil
}

func (c *onvifController) GetStreamUri(profileToken string) (*url.URL, error) {
	var err error
	if err = c.connect(); err != nil {
		return nil, err
	}

	c.l.Logf(logger.DebugLevel, "GetStreamUri...")

	ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
	defer cancel()

	request := media.GetStreamUri{ProfileToken: onvif.ReferenceToken(profileToken)}
	var response media.GetStreamUriResponse
	if err := c.dev.CreateRequest(request).WithContext(ctx).Do().Unmarshal(&response); err != nil {
		return nil, fmt.Errorf("method GetStreamUri failed: %w", err)
	}
	streamUrl, err := url.Parse(string(response.MediaUri.Uri))
	if err != nil {
		return nil, fmt.Errorf("parse stream URL failed: %w", err)
	}

	c.l.Logf(logger.DebugLevel, "GetStreamUri READY: %s", streamUrl)

	return streamUrl, nil
}

func (c *onvifController) GetProfiles() ([]string, error) {
	if err := c.connect(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
	defer cancel()

	c.l.Logf(logger.DebugLevel, "GetProfiles...")

	request := media.GetProfiles{}
	var response media.GetProfilesResponse
	if err := c.dev.CreateRequest(request).WithContext(ctx).Do().Unmarshal(&response); err != nil {
		return nil, fmt.Errorf("method GetProfiles failed: %w", err)
	}

	var result []string
	for i := range response.Profiles {
		result = append(result, string(response.Profiles[i].Token))
	}

	c.l.Logf(logger.DebugLevel, "GetProfiles READY: %+v", result)
	return result, nil
}

func (c *onvifController) connect() error {
	if c.dev != nil {
		return nil
	}

	c.l.Logf(logger.DebugLevel, "Connecting...")

	params := goonvif.DeviceParams{
		// TODO: адекватное определение endpoint
		Xaddr:      fmt.Sprintf("%s:80", c.u.Host),
		Username:   c.u.User.Username(),
		HttpClient: new(http.Client),
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

	c.l.Logf(logger.DebugLevel, "Connected")
	c.dev = dev
	return nil
}

func (c *onvifController) subscribe() error {
	if c.eventsEndpoint != "" {
		return nil
	}

	c.l.Logf(logger.DebugLevel, "Subscribing...")

	ctx, cancel := context.WithTimeout(c.ctx, maxNetworkTimeout)
	defer cancel()

	var response event.CreatePullPointSubscriptionResponse
	if err := c.dev.CreateRequest(event.CreatePullPointSubscription{}).WithContext(ctx).Do().Unmarshal(&response); err != nil {
		return fmt.Errorf("method CreatePullPointSubscription failed: %w", err)
	}

	c.l.Logf(logger.DebugLevel, "Subscribed")

	c.eventsEndpoint = string(response.SubscriptionReference.Address)
	return nil
}

func (c *onvifController) clearCache() {
	c.dev = nil
	c.eventsEndpoint = ""
	c.snapshotUrl = nil
}
