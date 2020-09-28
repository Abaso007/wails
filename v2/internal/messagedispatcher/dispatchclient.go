package messagedispatcher

import (
	"fmt"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// Client defines what a frontend client can do
type Client interface {
	Quit()
	NotifyEvent(message string)
	CallResult(message string)
	OpenDialog(dialogOptions *options.OpenDialog, callbackID string)
	SaveDialog(dialogOptions *options.SaveDialog, callbackID string)
	WindowSetTitle(title string)
	WindowShow()
	WindowHide()
	WindowCenter()
	WindowMaximise()
	WindowUnmaximise()
	WindowMinimise()
	WindowUnminimise()
	WindowPosition(x int, y int)
	WindowSize(width int, height int)
	WindowFullscreen()
	WindowUnFullscreen()
	WindowSetColour(colour int)
}

// DispatchClient is what the frontends use to interface with the
// dispatcher
type DispatchClient struct {
	id     string
	logger logger.CustomLogger

	bus *servicebus.ServiceBus

	// Client
	frontend Client
}

func newDispatchClient(id string, frontend Client, logger logger.CustomLogger, bus *servicebus.ServiceBus) *DispatchClient {

	return &DispatchClient{
		id:       id,
		frontend: frontend,
		logger:   logger,
		bus:      bus,
	}

}

// DispatchMessage is called by the front ends. It is passed
// an IPC message, translates it to a more concrete message
// type then publishes it on the service bus.
func (d *DispatchClient) DispatchMessage(incomingMessage string) {

	// Parse the message
	d.logger.Trace(fmt.Sprintf("Received message: %+v", incomingMessage))
	parsedMessage, err := message.Parse(incomingMessage)
	if err != nil {
		d.logger.Trace("Error: " + err.Error())
		return
	}

	// Save this client id
	parsedMessage.ClientID = d.id

	d.logger.Trace("I got a parsedMessage: %+v", parsedMessage)

	// Check error
	if err != nil {
		d.logger.Trace("Error: " + err.Error())
		// Hrm... what do we do with this?
		d.bus.PublishForTarget("generic:message", incomingMessage, d.id)
		return
	}

	// Publish the parsed message
	d.bus.PublishForTarget(parsedMessage.Topic, parsedMessage.Data, d.id)

}
