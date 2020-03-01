#ifndef COREV1_H_
#define COREV1_H_


#include <string>
#include <memory>

#include "k8s.io/api/core/v1/generated.pb.h"

#include "resource.h"
#include "../go/go.h"

namespace apiv1 = k8s::io::api::core::v1;

class CoreV1 {
public:
	CoreV1(uintptr_t clientsetH): clientsetH(clientsetH) {}

	auto Pods(const std::string& ns) {
		return Resource<apiv1::Pod, apiv1::PodList,
			k8s_client_corev1_pods_list,
			k8s_client_corev1_pods_watch, k8s_client_corev1_pods_stop_watch>(clientsetH, ns);
	}
	auto ReplicationControllers(const std::string& ns) {
		return Resource<apiv1::ReplicationController, apiv1::ReplicationControllerList,
			k8s_client_corev1_rcs_list,
			k8s_client_corev1_rcs_watch, k8s_client_corev1_rcs_stop_watch>(clientsetH, ns);
	}
	auto Nodes() {
		return Resource<apiv1::Node, apiv1::NodeList,
			k8s_client_corev1_nodes_list,
			k8s_client_corev1_nodes_watch, k8s_client_corev1_nodes_stop_watch>(clientsetH, "");
	}

private:
	uintptr_t clientsetH;
};


#endif // COREV1_H_