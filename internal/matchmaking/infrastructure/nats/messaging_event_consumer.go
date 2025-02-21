package matchnats

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/xfrr/randomtalk/internal/shared/messaging"
)

// MessagingEventConsumer consumes messaging events from a Nats JetStream stream.
type MessagingEventConsumer struct {
	js              jetstream.JetStream
	messageConsumer jetstream.Consumer
	logger          *zerolog.Logger
}

// SetupMessagingEventConsumer creates a new MessagingEventConsumer.
func SetupMessagingEventConsumer(
	ctx context.Context,
	nc *nats.Conn,
	logger *zerolog.Logger,
	streamName string,
	cfg jetstream.ConsumerConfig,
) (*MessagingEventConsumer, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	consumer, err := js.CreateOrUpdateConsumer(ctx, streamName, cfg)
	if err != nil {
		return nil, err
	}

	return &MessagingEventConsumer{
		js:              js,
		logger:          logger,
		messageConsumer: consumer,
	}, nil
}

func (c *MessagingEventConsumer) Subscribe(ctx context.Context) (<-chan *messaging.Event, error) {
	eventsCh := make(chan *messaging.Event)

	// start listening for msgCtx
	msgCtx, err := c.messageConsumer.Messages()
	if err != nil {
		return eventsCh, err
	}

	// start a goroutine to listen for messages
	c.startListening(ctx, msgCtx, eventsCh)
	return eventsCh, nil
}

func (c *MessagingEventConsumer) startListening(
	ctx context.Context,
	messages jetstream.MessagesContext,
	outCh chan<- *messaging.Event,
) {
	go func() {
		defer close(outCh)

		for {
			msg, err := messages.Next()
			if err != nil &&
				!errors.Is(err, nats.ErrConnectionClosed) &&
				!errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				c.logger.Error().
					Err(err).
					Msg("error getting next nats messaging event")
				return
			}

			// decode the message
			msgEvent := messaging.NewEvent()
			err = msgEvent.UnmarshalJSON(msg.Data())
			if err != nil {
				c.logger.Error().
					Err(err).
					Str("msg", string(msg.Data())).
					Msg("error unmarshalling nats messaging event")
				_ = msg.TermWithReason(jetstream.ErrInvalidDigestFormat.Error())
				continue
			}

			msgEvent.SetHeader(msg.Headers())

			go func(msg jetstream.Msg) {
				// wait for ack
				select {
				case <-ctx.Done():
					err = msg.Nak()
					if err != nil {
						c.logger.Error().
							Err(err).
							Msg("error rejecting nats messaging event")
					}
					return
				case <-msgEvent.WaitAck():
					// acknowledge the message
					err = msg.Ack()
					if err != nil {
						c.logger.Error().
							Err(err).
							Msg("error acknowledging nats messaging event")
						return
					}
				case <-msgEvent.WaitNack():
					err = msg.Nak()
					if err != nil {
						c.logger.Error().
							Err(err).
							Msg("error rejecting nats messaging event")
						return
					}
				case <-msgEvent.WaitReject():
					err = msg.Term()
					if err != nil {
						c.logger.Error().
							Err(err).
							Msg("error rejecting nats messaging event")
						return
					}
				}
			}(msg)

			// send the event to the output channel
			outCh <- msgEvent
		}
	}()
}
