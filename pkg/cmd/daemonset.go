package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/controlplaneio/kubectl-kubesec/pkg/kubesec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

var daemonsetCmd = &cobra.Command{
	Use:     `daemonset [name]`,
	Short:   "Scans daemonset object",
	Example: `  daemonset podinfo --namespace=default`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("daemonset name is required")
		}
		name := args[0]

		var buffer bytes.Buffer
		writer := bufio.NewWriter(&buffer)

		fmt.Println("scanning daemonset", name, "in namespace", namespace)
		ds, err := kubeClient.AppsV1beta2().DaemonSets(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ds.TypeMeta = metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		}
		err = serializer.Encode(ds, writer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		writer.Flush()

		result, err := kubesec.NewClient().ScanDefinition(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if result.Error != "" {
			fmt.Println(result.Error)
			os.Exit(1)
		}

		result.Dump(os.Stdout)

		return nil
	},
}
