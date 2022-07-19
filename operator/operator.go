package operator

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

	masterListOptions k8smetav1.ListOptions
	infraListOptions  k8smetav1.ListOptions
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

	t.masterListOptions = k8smetav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/master=",
	}

	t.infraListOptions = k8smetav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/infra=",
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

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
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

func (t *Operator) getPrismaConfig(ctx context.Context) (*prisma_types.PrismaConfig, error) {

	masterPods, err := t.kubeClient.CoreV1().Pods("").List(ctx, t.masterListOptions)

	if err != nil {
		return nil, err
	}

	prismaConfig := prisma_types.NewPrismaConfig(t.prismaLabel)

	for _, pod := range masterPods.Items {

		hostname := pod.Spec.Hostname
		podIP := pod.Status.PodIP

		if podIP == "" {
			continue
		}

		var tags []string
		tags = append(tags, "externalnetwork:name="+hostname)
		tags = append(tags, "externalnetwork:name=masterPods")

		var entries []string
		entries = append(entries, podIP+"/32")

		extNetwork := prisma_types.NewExternalnetwork(hostname).
			SetDescription("Auto generated").
			SetEntries(entries).SetAssociatedTags(tags)

		prismaConfig.AddExternalnetwork(extNetwork)
	}

	// infraPods, err := t.kubeClient.CoreV1().Pods("").List(ctx, t.infraListOptions)

	// if err != nil {
	// 	return nil, err
	// }

	return prismaConfig, nil
}

// masterNodeList, err := clientset.CoreV1().Nodes().List(ctx, opts)
