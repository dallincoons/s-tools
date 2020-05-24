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

type DBCloner struct {
	file *os.File
	cloneFrom string
	cloneTo string
	dumpDir string
	dumpName string
}

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
		dumpDir, _ := cmd.Flags().GetString("dump_dir")
		dumpName, _ := cmd.Flags().GetString("dump_name")
		cloneFrom, _ := cmd.Flags().GetString("from")
		cloneTo, _ := cmd.Flags().GetString("to")
		shouldSwitchDb, _ := cmd.Flags().GetBool("switch")

		dumpName = strings.TrimRight(dumpName, ".sql")

		file, ferr := os.OpenFile(fmt.Sprintf("%s/%s.sql", dumpDir, dumpName), os.O_RDWR|os.O_CREATE, 0644)

		if (ferr != nil) {
			log.Fatalf("error opening file: %s", ferr)
		}

		cloner := createCloner(file, cloneFrom, cloneTo, dumpDir, dumpName)

		switchDbIfRequired(shouldSwitchDb, cloner)

		file.Close()
	},
}

func createCloner(file *os.File, cloneFrom string, cloneTo string, dumpDir string, dumpName string) (*DBCloner) {
	return &DBCloner{
		file:      file,
		cloneFrom: cloneFrom,
		cloneTo:   cloneTo,
		dumpDir:   dumpDir,
		dumpName:  dumpName,
	}
}

func switchDbIfRequired(shouldSwitchDb bool, cloner *DBCloner) {
	if shouldSwitchDb {
		path, _ := os.Getwd()
		cloner.CloneAndSwitch(path)
	} else {
		cloner.Clone()
	}
}

func (this *DBCloner) CloneAndSwitch(environmentPath string) {
	this.Clone()

	RunSwitch(environmentPath, this.cloneTo)
}

func (this *DBCloner) Clone() {
	fmt.Println("Dumping database")
	dumpDatabase(this.file, this.cloneFrom)
	fmt.Println("Dumped")

	fmt.Println("Importing database")
	this.importDatabase()
	fmt.Println("Imported")

	fmt.Println("Cleaning up")
	removeDumpFiles(this.dumpDir, this.dumpName)

	fmt.Println(fmt.Sprintf("%s successfully cloned from %s", this.cloneTo, this.cloneFrom))
}

func removeDumpFiles(dump_dir string, dump_name string) {
	cleanup := exec.Command("rm", fmt.Sprintf("%s/%s.sql", dump_dir, dump_name), fmt.Sprintf("%s/%s.sql.bak", dump_dir, dump_name))

	if (Verbose) {
		fmt.Println(cleanup.String())
	}

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

func (this *DBCloner) importDatabase() {
	this.replaceDatabaseName()

	importCmd := exec.Command(
		"mysql",
		fmt.Sprintf("-h%s", viper.GetString("database.host")),
		fmt.Sprintf("-u%s", viper.GetString("database.username")),
		fmt.Sprintf("-p%s", viper.GetString("database.password")),
		fmt.Sprintf("-P%s", viper.GetString("database.port")),
	)

	addDumpToStdin(importCmd, this.file)

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

func (this *DBCloner) replaceDatabaseName() {
	replace := exec.Command(
		"sed",
		"-i.bak",
		fmt.Sprintf("s/%s/%s/g", this.cloneFrom, this.cloneTo),
		fmt.Sprintf("%s/%s.sql", this.dumpDir, this.dumpName),
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

	cloneCmd.MarkFlagRequired("to")
	cloneCmd.MarkFlagRequired("from")
}
