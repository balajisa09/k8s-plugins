package cmd

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"strings"
	"regexp"
	"os"
)

var podsonCmd = &cobra.Command{
	Use:   "podscreatedOn",
	Short: "Gets you the pods created within the given time",
	Args:  cobra.ExactArgs(1),
	Example: "podscreatedOn 2022-08-06",
	Run: func(cmd *cobra.Command, args []string) {
		//validate date string 
		match,_ := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}",args[0])
		if(!match){
			fmt.Println("Date format is YYYY-MM-DD")
			os.Exit(1)
		}
		var kubeconfig *string
		//set kube config path
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		//build config
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
		//create clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
		//Get all pods in default namespace
		pods, err := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		foundPods := false
		for _, pod := range pods.Items {
			timestampStr := pod.CreationTimestamp.String()
			timestamp := strings.Split(timestampStr," ")
			if(timestamp[0] == args[0]){
				fmt.Println(pod.Name)
				foundPods = true
			}

		}
		if(!foundPods){
			fmt.Println("No pods found on the given date")
		}

	},
}

func init() {
	rootCmd.AddCommand(podsonCmd)
}
