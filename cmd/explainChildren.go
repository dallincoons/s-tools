/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"surgio-tools/stool"

	"github.com/spf13/cobra"
)

// explainChildrenCmd represents the explainChildren command
var explainChildrenCmd = &cobra.Command{
	Use:   "explain_children",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		view_root := viper.GetString("views.root")
		view_name, _ := cmd.Flags().GetString("view")

		explainer := stool.GetExplainer(view_root)

		variables := explainer.CollectVariablesFromParents(view_name)

		for v, count := range variables {
			fmt.Fprintf(os.Stdout, "%-2d %-8s\n", count, v)
		}
	},
}

func init() {
	explainChildrenCmd.Flags().String("view", "", "specify name of view to explain")
	explainChildrenCmd.MarkFlagRequired("view")

	rootCmd.AddCommand(explainChildrenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explainChildrenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explainChildrenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
