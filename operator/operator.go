package operator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"
	token "github.com/aporeto-se/prisma-sdk-go-v2/token/env"
	prisma_types "github.com/aporeto-se/prisma-sdk-go-v2/types"
)

// kubectl get nodes -l node-role.kubernetes.io/master=

// Operator ...
type Operator struct {
	prismaAPI       string
	prismaLabel     string
	prismaNamespace string
	httpClient      *http.Client
	kubeClient      *kubernetes.Clientset
	prismaClient    *prisma_api.Client

	labelselectors []k8smetav1.ListOptions
}

// NewOperator ...
func NewOperator(ctx context.Context, config *Config) (*Operator, error) {

	if config.PrismaAPI == "" {
		panic("PrismaAPI is required")
	}

	if config.PrismaLabel == "" {
		panic("PrismaLabel is required")
	}

	if config.PrismaNamespace == "" {
		panic("PrismaNamespace is required")
	}

	if len(config.LabelSelectors) <= 0 {
		panic("LabelSelectors must have one or more entries")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	t := &Operator{
		prismaAPI:       config.PrismaAPI,
		prismaLabel:     config.PrismaLabel,
		prismaNamespace: config.PrismaNamespace,
		httpClient:      httpClient,
	}

	// node-role.kubernetes.io/master=
	// node-role.kubernetes.io/infra=

	for _, labelSelector := range config.LabelSelectors {
		t.labelselectors = append(t.labelselectors, k8smetav1.ListOptions{
			LabelSelector: labelSelector,
		})
	}

	err := t.init(ctx)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Run ...
func (t *Operator) Run(ctx context.Context) error {

	prismaConfig, err := t.getPrismaConfig(ctx)
	if err != nil {
		return err
	}

	err = t.prismaClient.ImportPrismaConfig(ctx, prismaConfig)
	if err != nil {
		return err
	}

	return nil
}

func (t *Operator) init(ctx context.Context) error {

	kubeClient, err := newKubeClientWithHomeConfig()
	if err != nil {
		return err
	}

	t.kubeClient = kubeClient

	tokenProvider, err := token.NewClient(token.NewConfig())
	if err != nil {
		return err
	}

	prismaClient, err := prisma_api.NewConfig().
		SetNamespace(t.prismaNamespace).
		SetAPI(t.prismaAPI).
		SetTokenProvider(tokenProvider).
		SetHTTPClient(t.httpClient).Build(ctx)

	if err != nil {
		return err
	}

	t.prismaClient = prismaClient

	return nil
}

func newKubeClientWithHomeConfig() (*kubernetes.Clientset, error) {

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func newKubeClientWithInCluster(ctx context.Context) (*kubernetes.Clientset, error) {

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func (t *Operator) getPrismaConfig(ctx context.Context) (*prisma_types.PrismaConfig, error) {

	config := prisma_types.NewPrismaConfig(t.prismaLabel)

	for _, labelselector := range t.labelselectors {

		extNetworks, err := t.getExternalNetworksWithSelector(ctx, labelselector)
		if err != nil {
			return nil, err
		}

		for _, extNetwork := range extNetworks {
			config.AddExternalnetwork(extNetwork)
		}

	}

	return config, nil
}

// masterNodeList, err := clientset.CoreV1().Nodes().List(ctx, opts)

func (t *Operator) getExternalNetworksWithSelector(ctx context.Context, labelSelector k8smetav1.ListOptions) ([]*prisma_types.Externalnetwork, error) {

	var externalNetworks []*prisma_types.Externalnetwork

	nodes, err := t.kubeClient.CoreV1().Nodes().List(ctx, labelSelector)

	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {

		log.Println("node.Name:", node.Name)

		filter := k8smetav1.ListOptions{
			FieldSelector: node.Name,
		}

		pods, err := t.kubeClient.CoreV1().Pods("").List(ctx, filter)

		if err != nil {
			return nil, err
		}

		for _, pod := range pods.Items {

			hostname := pod.Spec.Hostname
			podIP := pod.Status.PodIP

			if hostname == "" {
				log.Println("Pod missing hostname")
				continue
			}

			if podIP == "" {
				log.Println(fmt.Sprintf("Pod %s not added because it has no IP", hostname))
				continue
			}

			log.Println("hostname:", hostname)
			log.Println("podIP:", podIP)

			var tags []string
			tags = append(tags, "externalnetwork:name="+hostname)
			tags = append(tags, "externalnetwork:name=masterPods")

			var entries []string
			entries = append(entries, podIP+"/32")

			extNetwork := prisma_types.NewExternalnetwork(hostname).
				SetDescription("Auto generated").
				SetEntries(entries).SetAssociatedTags(tags)

			externalNetworks = append(externalNetworks, extNetwork)
		}

	}

	return externalNetworks, nil
}
