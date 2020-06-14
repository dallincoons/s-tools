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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"surgio-tools/stool"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		view_root := viper.GetString("views.root")
		view_name, _ := cmd.Flags().GetString("view")

		explainer := &stool.ViewExplainer{}
		finder := &stool.ViewFinder{
			view_root,
		}

		variables, _ := explainer.GetAllVariablesFrom(finder.GetFilePath(view_name))

		fmt.Println(variables)

		indexer := &stool.ViewIndexer{
			explainer,
				finder,
		}

		tree := indexer.Index(view_name)

		fmt.Println(tree.Children)
	},
}

func init() {
	explainCmd.Flags().String("view", "", "specify name of view to explain")
	explainCmd.MarkFlagRequired("view")

	rootCmd.AddCommand(explainCmd)
}
