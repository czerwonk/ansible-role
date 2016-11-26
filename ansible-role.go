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
	err := executeRole(roleName)

	if err != nil {
		os.Exit(2)
	}
}

func executeRole(roleName string) error {
	fmt.Println("Role: ", roleName)

	fileName := "/tmp/ansible-role-" + roleName + ".yml"
	fmt.Println("Creating temporary playbook file in ", fileName)
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
