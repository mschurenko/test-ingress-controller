package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type VirtualHost struct {
	Hostname  string
	ClusterIp string
	NodePort  int32
	Path      string
}

func main() {
	const annotationKey = "mschurenko/ingress.enabled"
	usage := fmt.Sprintf("Usage: %s -albTag", os.Args[0])

	var albTag string
	flag.StringVar(&albTag, "albTag", "ingressCtrl", "Set the Tag used to track ALB(s)")
	flag.Parse()

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	/* Use Watcher */
	// watchInterface, err := clientSet.ExtensionsV1beta1().Ingresses("").Watch(metav1.ListOptions{})
	// fmt.Println("watch interface type:")
	// fmt.Printf("%T\n", watchInterface)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var wg sync.WaitGroup

	// wg.Add(1)

	// chn := watchInterface.ResultChan()
	// go func(chn <-chan watch.Event) {
	// 	for event := range chn {
	// 		fmt.Println(event)
	// 	}
	// }(chn)
	// wg.Wait()

	for {
		ingresses, err := clientSet.ExtensionsV1beta1().Ingresses("").List(metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("ingresses:")
		for _, ingress := range ingresses.Items {
			usesIngress := ingress.ObjectMeta.Annotations[annotationKey]
			if usesIngress != "true" {
				continue
			}

			path := ingress.Spec.Rules[0].HTTP.Paths[0].Path
			svcName := ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName
			svcPort := ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort.StrVal
			host := ingress.Spec.Rules[0].Host
			fmt.Println(path)
			fmt.Println(svcName, svcPort, host)

			// get clusterIp and nodePort from service
			svc, err := clientSet.CoreV1().Services("default").Get(svcName, metav1.GetOptions{})

			if errors.IsNotFound(err) {
				fmt.Printf("svc %s is not found\n", svcName)
				continue
			} else if err != nil {
				log.Fatal(err)
			}

			clusterIP := svc.Spec.ClusterIP
			nodePort := svc.Spec.Ports[0].NodePort

			v := VirtualHost{
				host,
				clusterIP,
				nodePort,
				path,
			}

			fmt.Println(v)

			// create ALB if not found or existing ALB(s) don't have enough capacity

			// add target group and listener to newest ALB

			// add route53 record for host

			// update ingress status.loadBalancer

		}

		time.Sleep(10 * time.Second)
	}
}
