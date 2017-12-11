package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const version string = "0.4.1"

var (
	showVersion = flag.Bool("version", false, "Show version")
	force       = flag.Bool("force", false, "Disables safty checks")
	debug       = flag.Bool("debug", false, "Prints debug output")
	args        []string
)

func init() {
	flag.Usage = printUsage
	args = make([]string, 0)
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	args = flag.Args()
	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	roleName := args[0]
	hosts := args[1]

	if hosts == "all" && !*force {
		fmt.Println("Executing roles for group all might be dangerous. Please add -force parameter to allow this action.")
		os.Exit(1)
	}

	err := executeRole(roleName, hosts)

	if err != nil {
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("Usage: ansible-role [ -force ] $rolename $servers [ any further ansible-playbook parameters ]\n")
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Println("ansible-role")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author: Daniel Czerwonk")
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

	var w io.Writer = f
	if *debug {
		w = &debugWriter{writer: w}
	}

	writeFileContent(roleName, hosts, w)
}

func deleteFile(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("Deleting %s\n", fileName)
		os.Remove(fileName)
	}
}

func writeFileContent(roleName string, hosts string, w io.Writer) {
	fmt.Println("Generating playbook content")

	fmt.Fprintln(w, "---")
	fmt.Fprintf(w, "- hosts: %s\n", hosts)
	fmt.Fprintln(w, "  roles:")
	fmt.Fprintln(w, "  - ", roleName)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func executeAnsible(fileName string) error {
	fmt.Println("Starting ansible playbook")

	ansibleArgs := args[2:]
	ansibleArgs = append(ansibleArgs, fileName)

	cmd := exec.Command("ansible-playbook", ansibleArgs...)

	if *debug {
		fmt.Printf("%v %s\n", cmd.Path, cmd.Args)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
