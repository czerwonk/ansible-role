package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: You have to enter a role name!")
		os.Exit(1)
	}

	roleName := os.Args[1]
	executeRole(roleName)
}

func executeRole(name string) {
	fmt.Println("Role: ", name)

	fileName := "/tmp/ansible-role-" + name + ".yml"
	fmt.Println("Creating temporary playbook file in ", fileName)
	createFile(name, fileName)
	defer deleteFile(fileName)

	executeAnsible(fileName)
}

func createFile(roleName string, fileName string) {
	f, err := os.Create(fileName)
	checkError(err)
	defer f.Close()

	writeFileContent(roleName, f)
}

func deleteFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		fmt.Println("Deleting ", fileName)
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

func executeAnsible(fileName string) {
	fmt.Println("Starting ansible playbook")

	ansibleArgs := os.Args[2:]
	ansibleArgs = append(ansibleArgs, fileName)
	cmd := exec.Command("ansible-playbook", ansibleArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
