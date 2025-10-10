package truenas

import (
	"context"
	"log"

	"github.com/puzpuzpuz/xsync/v3"
)

type ClientSubscribe struct {
	client      *Client
	subscribeId *xsync.MapOf[string, string]
	subscribeCh *xsync.MapOf[string, chan struct{}]
}

func NewClientSubscribe(client *Client) *ClientSubscribe {
	return &ClientSubscribe{
		client:      client,
		subscribeId: xsync.NewMapOf[string, string](),
		subscribeCh: xsync.NewMapOf[string, chan struct{}](),
	}
}

// Subscribe allows subscribing to TrueNAS events via WebSocket
func (cs *ClientSubscribe) Subscribe(ctx context.Context, collection string, collectionUpdate func(Message) error) error {
	var subscribeId string
	ch := make(chan Message, 10)
	cs.client.pending.Store(collection, ch) // add pending channels

	if err := cs.client.Call(ctx, "core.subscribe", []any{collection}, &subscribeId); err != nil {
		log.Printf("error subscribing to collection %s: %v", collection, err)
		return err
	}

	log.Printf("subscribed to collection %s with subscribeId %s", collection, subscribeId)

	unsubscribe := make(chan struct{}, 1)
	cs.subscribeId.Store(collection, subscribeId)
	cs.subscribeCh.Store(collection, unsubscribe)

	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Panicf("channel closed unexpectedly for collection %s", collection)
					return
				}

				if err := collectionUpdate(msg); err != nil {
					// Handle error (e.g., log it)
					log.Printf("error handling collection update: %v", err)
				}
			case <-unsubscribe:
				log.Printf("received unsubscribe signal for collection %s", collection)
				close(ch)
				// stop the goroutine
				return
			}
		}
		log.Printf("stopped listening to collection %s", collection)
	}()

	return nil
}

func (cs *ClientSubscribe) Unsubscribe(ctx context.Context, collection string) error {
	subscribeId, ok := cs.subscribeId.LoadAndDelete(collection)
	if !ok {
		log.Printf("no active subscription found for collection %s", collection)
		return nil
	}

	if err := cs.client.Call(ctx, "core.unsubscribe", []any{subscribeId}, nil); err != nil {
		log.Printf("error unsubscribing from collection %s: %v", collection, err)
		return err
	}
	cs.client.pending.Delete(collection)

	ch, ok := cs.subscribeCh.LoadAndDelete(collection)
	if ok {
		ch <- struct{}{}
		close(ch)
	}
	log.Printf("unsubscribed from collection %s with subscribeId %s", collection, subscribeId)
	return nil
}
