package main

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/golang/glog"

	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apimachinery/pkg/util/clock"
)

// #include <stdint.h>
// typedef void (*k8s_client_watch_callback_fn)(uintptr_t watchKey, int watchType, void* objBytes, int objSize);
// void k8s_client_watch_callback_wrapper(uintptr_t callbackFn, uintptr_t callbackContext, int watchType, void* objBytes, int objSize) {
// 	((k8s_client_watch_callback_fn)callbackFn)(callbackContext, watchType, objBytes, objSize);
// }
import "C"

const (
	Closed = 0
	Added = 1
	Modified = 2
	Deleted = 3
	Error = -1
)
type WatchEventHandler interface {
	HandleEvent(eventType int, obj interface{}) error
}

var watchMu sync.Mutex
var watchMap = map[C.uintptr_t]chan struct{}{}


// watchHandler watches w and keeps *resourceVersion up to date.
func watchHandler (w watch.Interface, expectedType interface{}, handler WatchEventHandler, stopCh <-chan struct{}) (string, error) {
	start := clock.RealClock{}.Now()
	eventCount := 0
	resourceVersion := ""

	// Stopping the watcher should be idempotent and if we return from this function there's no way
	// we're coming back in with the same watch interface.
	defer w.Stop()

	loop:
	for {
		select {
		case <-stopCh:
			glog.Infof("Stop requested")
			handler.HandleEvent(Closed, nil)
			return "", errors.New("Stop requested")
		case event, ok := <-w.ResultChan():
			if !ok {
				break loop
			}
			if event.Type == watch.Error {
				glog.Errorf("watchHandler: %v", apierrs.FromObject(event.Object))
				handler.HandleEvent(Error, -1)
				return "", apierrs.FromObject(event.Object)
			}
			if e, a := expectedType, reflect.TypeOf(event.Object); e != nil && e != a {
				glog.Infof("watchHandler: expected type %v, but watch event object had type %v", e, a)
				continue
			}
			meta, err := meta.Accessor(event.Object)
			if err != nil {
				glog.Infof("watchHandler: unable to understand watch event %#v", event)
				continue
			}
			resourceVersion = meta.GetResourceVersion()
			switch event.Type {
			case watch.Added:
				handler.HandleEvent(Added, event.Object)
				glog.Infof("watch.Added: %v, %v/%v\n", expectedType, meta.GetNamespace(), meta.GetName())
			case watch.Modified:
				handler.HandleEvent(Modified, event.Object)
				glog.Infof("watch.Modified: %v, %v/%v\n", expectedType, meta.GetNamespace(), meta.GetName())
			case watch.Deleted:
				handler.HandleEvent(Deleted, event.Object)
				glog.Infof("watch.Deleted: %v, %v/%v\n", expectedType, meta.GetNamespace(), meta.GetName())
				// TODO: Will any consumers need access to the "last known
				// state", which is passed in event.Object? If so, may need
				// to change this.
			default:
				glog.Errorf("watchHandler: unable to understand watch event %#v", event)
			}
			eventCount++
		}
	}

	watchDuration := clock.RealClock{}.Now().Sub(start)
	if watchDuration < 1*time.Second && eventCount == 0 {
		glog.Errorf("watchHandler: Unexpected watch close - watch lasted less than a second and no items received")
		handler.HandleEvent(Error, nil)
		return "", errors.New("very short watch")
	}
	glog.Infof("watchHandler: Watch close - %v total %v items received", expectedType, eventCount)
	handler.HandleEvent(Closed, nil)
	return resourceVersion, nil
}
