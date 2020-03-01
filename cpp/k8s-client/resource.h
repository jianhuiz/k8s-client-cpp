#ifndef RESOURCE_H_
#define RESOURCE_H_


#include <string>
#include <memory>

#include "k8s.io/api/core/v1/generated.pb.h"
#include "../go/go.h"

namespace apiv1 = k8s::io::api::core::v1;
namespace metav1 = k8s::io::apimachinery::pkg::apis::meta::v1;


template <class T, class TList,
		decltype(k8s_client_corev1_pods_list) listFunc,
		decltype(k8s_client_corev1_pods_watch) watchFunc,
		decltype(k8s_client_corev1_pods_stop_watch) stopWatchFunc>
class Resource {
public:
	Resource(uintptr_t clientsetH, const std::string& ns): clientsetH(clientsetH), ns(ns) {}

	TList List(const metav1::ListOptions& listOptions) {
		std::uint8_t listOptionsBytes[listOptions.ByteSize()];
		listOptions.SerializeToArray(listOptionsBytes, sizeof(listOptionsBytes));
		void *objListBytes = NULL;
		int objListSize = 0;
		auto err = listFunc(clientsetH, (char*)ns.data(), listOptionsBytes, sizeof(listOptionsBytes), &objListBytes, &objListSize);
		if (err != NULL) {
			auto errStr = std::string(err);
			free(err);
			throw errStr;
		}

		TList objList;
		objList.ParseFromArray(objListBytes, objListSize);
		free(objListBytes);
		return objList;
	}

	typedef std::function<void(int watchType, const T*)> WatchCallbackFn;
	uintptr_t Watch(const metav1::ListOptions& listOptions, WatchCallbackFn callback) {
		std::uint8_t listOptionsBytes[listOptions.ByteSize()];
		listOptions.SerializeToArray(listOptionsBytes, sizeof(listOptionsBytes));
		WatchCallbackFn *callbackP = new WatchCallbackFn(callback);
		auto callbackFn = uintptr_t(Resource::WatchCallback);
		auto callbackContext = uintptr_t(callbackP);
		auto err = watchFunc(clientsetH, (char*)ns.data(), listOptionsBytes, sizeof(listOptionsBytes), callbackFn, callbackContext);
		if (err != NULL) {
			auto errStr = std::string(err);
			free(err);
			throw errStr;
		}
		return callbackContext;
	}

	void StopWatch(uintptr_t watchH) {
		if (k8s_client_corev1_pods_stop_watch(watchH)) {
			delete (WatchCallbackFn*)watchH;
		}
	}

private:
	static void WatchCallback(uintptr_t callbackContext, int watchType, void* objBytes, int objSize) {
		T obj;
		obj.ParseFromArray(objBytes, objSize);
		(*(WatchCallbackFn*)callbackContext)(watchType, &obj);
	}

private:
	std::string ns;
	uintptr_t clientsetH;
};


#endif // RESOURCE_H_