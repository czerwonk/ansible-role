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

	if len(os.Args) < 3 {
		fmt.Println("Usage: ansible-role $role $hosts [ $ansible_args ]")
		os.Exit(1)
	}

	roleName := os.Args[1]
	hosts := os.Args[2]
	err := executeRole(roleName, hosts)

	if err != nil {
		os.Exit(2)
	}
}

func printVersion() {
	fmt.Println("ansible-role")
	fmt.Printf("Version: %s\n", version)
}

func executeRole(roleName string, hosts string) error {
	fmt.Printf("Role: %s\n", roleName)

	fileName := "/tmp/ansible-role-" + roleName + ".yml"
	fmt.Printf("Creating temporary playbook file in %s\n", fileName)
	createFile(roleName, hosts, fileName)
	defer deleteFile(fileName)

	return executeAnsible(fileName)
}

func createFile(roleName string, hosts string, fileName string) {
	f, err := os.Create(fileName)
	checkError(err)
	defer f.Close()

	writeFileContent(roleName, hosts, f)
}

func deleteFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("Deleting %s\n", fileName)
		os.Remove(fileName)
	}
}

func writeFileContent(roleName string, hosts string, f *os.File) {
	fmt.Println("Generating playbook content")

	fmt.Fprintln(f, "---")
	fmt.Fprintf(f, "- hosts: %s\n", hosts)
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

	ansibleArgs := os.Args[3:]
	ansibleArgs = append(ansibleArgs, fileName)

	cmd := exec.Command("ansible-playbook", ansibleArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
