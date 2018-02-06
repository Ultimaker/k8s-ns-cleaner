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
	fmt.Println("start ns-cleaner")

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
		fmt.Print("check if ", ns.Name, " has correct labels... ")
		labels := ns.Annotations

		if val, ok := labels["expiresAt"]; ok {
			fmt.Print("yes ")
			expiresAt, err := time.Parse(format, val)

			if err != nil {
				fmt.Println(err)
				continue
			}

			if time.Now().After(expiresAt) {
				fmt.Println("expired, remove namespace")
				err := cli.CoreV1().Namespaces().Delete(ns.Name, &v1.DeleteOptions{})

				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("no")
		}
	}
}
