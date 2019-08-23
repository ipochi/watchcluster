package watchcluster

import (
	log "github.com/golang/glog"
	"github.com/ipochi/watchcluster/config"
	"github.com/nlopes/slack"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var ResourceInformerMap map[string]cache.SharedIndexInformer

func Start(watchclusterConfig *config.Config) {

	log.Info("Starting kubewatch controller")

	ResourceInformerMap = make(map[string]cache.SharedIndexInformer)

	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	SharedInformerFactory := informers.NewSharedInformerFactory(clientset, 0)

	log.Info("watchclusterconfig resources list --- ", watchclusterConfig.Resources)
	for _, resource := range watchclusterConfig.Resources {
		log.Info("Resource name --- ", resource.Name)
		ResourceInformerMap[resource.Name] = SharedInformerFactory.Core().V1().Pods().Informer()
	}

	log.Info("Length of resourceinformermap --- ", len(ResourceInformerMap))
	for _, informer := range ResourceInformerMap {
		log.Info("Inside event handler ... ")
		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				_, err := cache.MetaNamespaceKeyFunc(obj)
				if err != nil {
					log.Errorf("Failed to get MetaNamespaceKey from event resource")
					return
				}
				podObj, ok := obj.(*v1.Pod)
				if !ok {
					return
				}

				log.Infof("PodObj image--- %v", podObj.Spec.Containers[0].Image)

				// Kind of object
				kind := strings.ToLower(podObj.Kind)
				log.Infof("Object --- %v", obj)
				log.Infof("Kind of object --- %s", kind)
				sendMessageToSlack(obj, watchclusterConfig, kind)

			},
		})
	}

	// factory := informers.NewSharedInformerFactory(clientset, 0)

	// podInformer := factory.Core().V1().Pods().Informer()

	// podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	DeleteFunc: onDelete,
	// })

	stopCh := make(chan struct{})
	defer close(stopCh)
	defer runtime.HandleCrash()

	go SharedInformerFactory.Start(stopCh)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm

	log.Info("Controller terminated ....")
}

func sendMessageToSlack(obj interface{}, watchclusterConfig *config.Config, kind string) {
	log.Infof("Sending mesage to slack ::: ")

	token := watchclusterConfig.EventHandler.Slack.Token
	channel := watchclusterConfig.EventHandler.Slack.Channel

	api := slack.New(token)
	params := slack.PostMessageParameters{
		AsUser: true,
	}

	message := "Hello I am coming from the kubernetes cluster ... >> "
	channelID, timestamp, err := api.PostMessage(channel, message, params)
	if err != nil {
		log.Errorf("Error in sending slack message %s", err.Error())
		return
	}

	log.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
}
