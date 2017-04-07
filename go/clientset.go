package main

import (
	"sync"
	"unsafe"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// #include <stdint.h>
import "C"

var clientsetMu sync.Mutex
var clientsetMap = map[C.uintptr_t]kubernetes.Interface{}

//export k8s_client_new_from_kubeconfig
func k8s_client_new_from_kubeconfig(masterUrl, kubeconfigPath *C.char, oClientsetH *C.uintptr_t) *C.char {
	config, err := clientcmd.BuildConfigFromFlags(C.GoString(masterUrl), C.GoString(kubeconfigPath))
	if err != nil {
		return C.CString(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return C.CString(err.Error())
	}

	clientsetMu.Lock()
	defer clientsetMu.Unlock()
	*oClientsetH = C.uintptr_t(uintptr((unsafe.Pointer(clientset))))
	clientsetMap[*oClientsetH] = clientset
	return nil
}

//export k8s_client_delete
func k8s_client_delete(clientsetH C.uintptr_t) {
	clientsetMu.Lock()
	defer clientsetMu.Unlock()
	delete(clientsetMap, clientsetH)
}


