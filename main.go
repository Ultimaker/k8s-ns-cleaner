package main

import (
	"os"
	"fmt"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	config, err := rest.InClusterConfig()
	format := "20060102T150405"

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cli, err := kubernetes.NewForConfig(config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	nsList, err := cli.CoreV1().Namespaces().List(v1.ListOptions{})

	for _, ns := range nsList.Items {
		labels := ns.Annotations

		if val, ok := labels["expiresAt"]; ok {
			expiresAt, err := time.Parse(format, val)

			if err != nil {
				fmt.Println(err)
				continue
			}

			if time.Now().After(expiresAt) {
				err := cli.CoreV1().Namespaces().Delete(ns.Name, &v1.DeleteOptions{})

				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
