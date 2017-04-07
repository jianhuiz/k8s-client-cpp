package main

import (
	"reflect"
	"unsafe"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

// #include <stdlib.h>
// #include <stdint.h>
// typedef void (*k8s_client_watch_callback_fn)(uintptr_t watchKey, int watchType, void* objBytes, int objSize);
// extern void k8s_client_watch_callback_wrapper(uintptr_t callbackFn, uintptr_t callbackContext, int watchType, void* objBytes, int objSize);
import "C"

//export k8s_client_corev1_nodes_list
func k8s_client_corev1_nodes_list(clientsetKey C.uintptr_t, namespace *C.char, optsBytes unsafe.Pointer, optsSize C.int, oBytes *unsafe.Pointer, oSize *C.int) *C.char {
	listOptions := metav1.ListOptions{}
	listOptions.Unmarshal(no_copy_slice_from_c_array(optsBytes, optsSize))
	objList, err := clientsetMap[clientsetKey].CoreV1().Nodes().List(listOptions)
	if err != nil {
		return C.CString(err.Error())
	}
	resultProto, _ := objList.Marshal()
	*oBytes = C.CBytes(resultProto)
	*oSize = C.int(len(resultProto))
	return nil
}

type nodeWatchHandlerFunc struct {
	callbackFn C.uintptr_t
	callbackContext C.uintptr_t
}
func (h *nodeWatchHandlerFunc) HandleEvent(eventType int, o interface{}) error {
	obj, ok := o.(*apiv1.Node)
	if ok {
		objProto, _ := obj.Marshal()
		objBytes := unsafe.Pointer(&objProto[0])
		objSize := C.int(len(objProto))
		C.k8s_client_watch_callback_wrapper(C.uintptr_t(h.callbackFn), C.uintptr_t(h.callbackContext), C.int(eventType), objBytes, objSize)
	} else {
		C.k8s_client_watch_callback_wrapper(C.uintptr_t(h.callbackFn), C.uintptr_t(h.callbackContext), C.int(eventType), nil, 0)
	}
	return nil
}

//export k8s_client_corev1_nodes_watch
func k8s_client_corev1_nodes_watch(clientsetKey C.uintptr_t, namespace *C.char, optsBytes unsafe.Pointer, optsSize C.int,
callbackFn C.uintptr_t, callbackContext C.uintptr_t) *C.char {
	listOptions := metav1.ListOptions{}
	listOptions.Unmarshal(no_copy_slice_from_c_array(optsBytes, optsSize))
	listOptions.ResourceVersion = "0"
	watch, err := clientsetMap[clientsetKey].CoreV1().Nodes().Watch(listOptions);
	if err != nil {
		return C.CString(err.Error())
	}

	watchMu.Lock()
	defer watchMu.Unlock()
	stopCh := make(chan struct{})
	watchMap[callbackContext] = stopCh

	go watchHandler(watch, reflect.TypeOf(&apiv1.Node{}), &nodeWatchHandlerFunc{callbackFn, callbackContext}, stopCh)

	return nil
}

//export k8s_client_corev1_nodes_stop_watch
func k8s_client_corev1_nodes_stop_watch(watchH C.uintptr_t) C.int {
	watchMu.Lock()
	defer watchMu.Unlock()
	if stopCh, ok := watchMap[watchH]; ok {
		close(stopCh)
		delete(watchMap, watchH)
		return 0
	}
	return 1
}
