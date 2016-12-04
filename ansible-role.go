/*
Copyright 2016 Daniel Czerwonk (d.czerwonk@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const version string = "0.1"

var (
	showVersion = flag.Bool("version", false, "Show version")
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("Error: You have to enter a role name!")
		os.Exit(1)
	}

	roleName := os.Args[1]
	err := executeRole(roleName)

	if err != nil {
		os.Exit(2)
	}
}

func printVersion() {
	fmt.Println("ansible-role")
	fmt.Printf("Version: %s\n", version)
}

func executeRole(roleName string) error {
	fmt.Printf("Role: %s\n", roleName)

	fileName := "/tmp/ansible-role-" + roleName + ".yml"
	fmt.Printf("Creating temporary playbook file in %s\n", fileName)
	createFile(roleName, fileName)
	defer deleteFile(fileName)

	return executeAnsible(fileName)
}

func createFile(roleName string, fileName string) {
	f, err := os.Create(fileName)
	checkError(err)
	defer f.Close()

	writeFileContent(roleName, f)
}

func deleteFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("Deleting %s\n", fileName)
		os.Remove(fileName)
	}
}

func writeFileContent(roleName string, f *os.File) {
	fmt.Println("Generating playbook content")

	fmt.Fprintln(f, "---")
	fmt.Fprintln(f, "- hosts: all")
	fmt.Fprintln(f, "  roles:")
	fmt.Fprintln(f, "  - ", roleName)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func executeAnsible(fileName string) error {
	fmt.Println("Starting ansible playbook")

	ansibleArgs := os.Args[2:]
	ansibleArgs = append(ansibleArgs, fileName)

	cmd := exec.Command("ansible-playbook", ansibleArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
