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
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dump := exec.Command(
			"mysqldump",
				fmt.Sprintf("%s", viper.GetString("database.database")),
				fmt.Sprintf("-h%s", viper.GetString("database.host")),
				fmt.Sprintf("-u%s", viper.GetString("database.username")),
				fmt.Sprintf("-p%s", viper.GetString("database.password")),
				fmt.Sprintf("-P%s", viper.GetString("database.port")),
				"--databases",
		)

		fmt.Println(dump.String())

		to, _ := cmd.Flags().GetString("to")
		name, _ := cmd.Flags().GetString("name")
		importDatabase, _ := cmd.Flags().GetString("database")

		file, ferr := os.OpenFile(fmt.Sprintf("%s/%s.sql", to, name), os.O_RDWR|os.O_CREATE, 0644)

		dump.Stdout = file

		err := dump.Run(); if err != nil {
			log.Fatalln(err.Error())
		}

		if (ferr != nil) {
			log.Fatalf("error opening file: %s", ferr)
		}

		replace := exec.Command(
			"sed",
			"-i.bak",
			fmt.Sprintf("s/%s/%s/g", viper.GetString("database.database"), importDatabase),
			fmt.Sprintf("%s/%s.sql", to, name),
		)

		fmt.Println(replace.String())

		replace.Run()

		importCmd := exec.Command(
			"mysql",
			fmt.Sprintf("-h%s", viper.GetString("database.host")),
			fmt.Sprintf("-u%s", viper.GetString("database.username")),
			fmt.Sprintf("-p%s", viper.GetString("database.password")),
			fmt.Sprintf("-P%s", viper.GetString("database.port")),
		)

		stdin, err := importCmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		err = importCmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		bytes, _ := ioutil.ReadFile(file.Name())
		_, err = io.WriteString(stdin, string(bytes))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(importCmd.String())

		output, _ := importCmd.Output()

		fmt.Println(string(output))

		file.Close()
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	home, _ := os.UserHomeDir()

	cloneCmd.Flags().String("to", home, "Specify a directory to output the dump file")
	cloneCmd.Flags().String("name", "dump", "Specify the name of the dump file")
	cloneCmd.Flags().String("database", "cloned", "Specify a directory to output the dump file")
}
