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
	"strings"
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
		dump_dir, _ := cmd.Flags().GetString("dump_dir")
		dump_name, _ := cmd.Flags().GetString("dump_name")
		databaseToCloneName, _ := cmd.Flags().GetString("from")
		importDatabaseName, _ := cmd.Flags().GetString("to")
		shouldSwitchDb, _ := cmd.Flags().GetBool("switch")

		dump_name = strings.TrimRight(dump_name, ".sql")

		file, ferr := os.OpenFile(fmt.Sprintf("%s/%s.sql", dump_dir, dump_name), os.O_RDWR|os.O_CREATE, 0644)

		if (ferr != nil) {
			log.Fatalf("error opening file: %s", ferr)
		}

		cloner := &DBCloner{
			file: file,
			cloneFrom: databaseToCloneName,
			cloneTo: importDatabaseName,
			dumpDir: dump_dir,
			dumpName: dump_name,
		}

		if (shouldSwitchDb) {
			path, _ := os.Getwd()
			cloner.CloneAndSwitch(path)
		} else {
			cloner.Clone()
		}

		file.Close()
	},
}

type DBCloner struct {
	file *os.File
	cloneFrom string
	cloneTo string
	dumpDir string
	dumpName string
}

func (this *DBCloner) CloneAndSwitch(environmentPath string) {
	this.Clone()

	RunSwitch(environmentPath, this.cloneTo)
}

func (this *DBCloner) Clone() {
	fmt.Println("Cloning database")
	dumpDatabase(this.file, this.cloneFrom)

	importDatabase(this.cloneTo, this.dumpDir, this.dumpName, this.file)

	fmt.Println("Cleaning up")
	removeDumpFiles(this.dumpDir, this.dumpName)

	fmt.Println(fmt.Sprintf("%s successfully cloned from %s", this.cloneTo, this.cloneFrom))
}

func removeDumpFiles(dump_dir string, dump_name string) {
	cleanup := exec.Command("rm", fmt.Sprintf("%s/%s.sql", dump_dir, dump_name), fmt.Sprintf("%s/%s.sql.bak", dump_dir, dump_name))

	cleanup.Run()
}

func dumpDatabase(file *os.File, databaseToCloneName string) {
	dump := exec.Command(
		"mysqldump",
		fmt.Sprintf("%s", databaseToCloneName),
		fmt.Sprintf("-h%s", viper.GetString("database.host")),
		fmt.Sprintf("-u%s", viper.GetString("database.username")),
		fmt.Sprintf("-p%s", viper.GetString("database.password")),
		fmt.Sprintf("-P%s", viper.GetString("database.port")),
		"--databases",
	)

	if Verbose {
		fmt.Println(dump.String())
	}

	dump.Stdout = file

	err := dump.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func importDatabase(importDatabase string, to string, name string, file *os.File) {
	replaceDatabaseName(importDatabase, to, name)

	importCmd := exec.Command(
		"mysql",
		fmt.Sprintf("-h%s", viper.GetString("database.host")),
		fmt.Sprintf("-u%s", viper.GetString("database.username")),
		fmt.Sprintf("-p%s", viper.GetString("database.password")),
		fmt.Sprintf("-P%s", viper.GetString("database.port")),
	)

	addDumpToStdin(importCmd, file)

	importCmd.Wait()

	if (Verbose) {
		fmt.Println(importCmd.String())
	}
}

func addDumpToStdin(importCmd *exec.Cmd, file *os.File) {
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

	stdin.Close()
}

func replaceDatabaseName(importDatabase string, to string, name string) {
	replace := exec.Command(
		"sed",
		"-i.bak",
		fmt.Sprintf("s/%s/%s/g", viper.GetString("database.database"), importDatabase),
		fmt.Sprintf("%s/%s.sql", to, name),
	)

	if Verbose {
		fmt.Println(replace.String())
	}

	replace.Run()
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	home, _ := os.UserHomeDir()

	cloneCmd.Flags().String("dump_dir", home, "Specify a directory to output the dump file")
	cloneCmd.Flags().String("dump_name", "dump", "Specify the name of the dump file")
	cloneCmd.Flags().String("to", "cloned", "Specify name of new database")
	cloneCmd.Flags().String("from", viper.GetString("database.database"), "Specify name of database to clone")
	cloneCmd.Flags().Bool("switch", true, "Specify whether to switch databases in environment")
}
