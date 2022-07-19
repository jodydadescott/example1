package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // register GCP auth provider
	"k8s.io/client-go/tools/clientcmd"
)

// k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// "k8s.io/client-go/kubernetes"
// "k8s.io/client-go/tools/clientcmd"

// prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"
// token "github.com/aporeto-se/prisma-sdk-go-v2/token/env"
// prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"

// kubectl get nodes -l node-role.kubernetes.io/master=

func main() {

	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	ctx := context.Background()

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	opts := k8smetav1.ListOptions{
		//LabelSelector: "node-role.kubernetes.io/master=",
	}

	pods, err := clientset.CoreV1().Pods("").List(ctx, opts)

	if err != nil {
		return err
	}

	// IP address of the host to which the pod is assigned. Empty if not yet scheduled.
	// HostIP string `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
	// IP address allocated to the pod. Routable at least within the cluster.
	// Empty if not yet allocated.
	// PodIP string `json:"podIP,omitempty" protobuf:"bytes,6,opt,name=podIP"`
	//
	// podIPs holds the IP addresses allocated to the pod. If this field is specified, the 0th entry must
	// match the podIP field. Pods may be allocated at most 1 value for each of IPv4 and IPv6. This list
	// is empty if no IPs have been allocated yet.
	// PodIPs []PodIP `json:"podIPs,omitempty" protobuf:"bytes,12,rep,name=podIPs" patchStrategy:"merge" patchMergeKey:"ip"`

	for _, pod := range pods.Items {
		hostname := pod.Spec.Hostname
		podIP := pod.Status.PodIP
		fmt.Println("hostname: " + hostname)
		fmt.Println("podIP: " + podIP)
	}

	// for _, node := range masterNodeList.Items {
	// 	node.
	// }

	//service, err := t.KubernetesClientset.CoreV1().Services("kube-system").Get(ctx, "kube-dns", k8smetav1.GetOptions{})

	// nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})

	// tokenProvider, err := token.NewClient(token.NewConfig())
	// if err != nil {
	// 	return err
	// }

	// httpClient := &http.Client{}

	// prismaClient, err := prisma_api.NewConfig().
	// 	SetNamespace(namespace).
	// 	SetAPI(api).
	// 	SetTokenProvider(tokenProvider).
	// 	SetHTTPClient(httpClient).Build(ctx)

	// if err != nil {
	// 	return err
	// }

	// var entries []string
	// entries = append(entries, "1.1.1.1/32")
	// entries = append(entries, "2.2.2.2/32")

	// extNetwork := prisma_types.NewExternalnetwork("MyNetwork").SetDescription("My Description").SetEntries(entries)
	// prismaConfig.AddExternalnetwork(extNetwork)

	// err = prismaClient.ImportPrismaConfig(ctx, prismaConfig)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// masterNodeList, err := clientset.CoreV1().Nodes().List(ctx, opts)
