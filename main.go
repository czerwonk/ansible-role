package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

const version string = "0.4.4"

var (
	showVersion = flag.Bool("version", false, "Show version")
	force       = flag.Bool("force", false, "Disables safety checks")
	debug       = flag.Bool("debug", false, "Prints debug output")
	gatherFacts = flag.Bool("gather-facts", true, "Gather information of target hosts")
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
	fmt.Print("Usage: ansible-role [ -force ] $rolename $servers [ any further ansible-playbook parameters ]\n\n")
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Println("ansible-role")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author: Daniel Czerwonk")
}

func executeRole(roleName string, hosts string) error {
	fmt.Printf("Role: %s\n", roleName)

	f, err := os.CreateTemp("", "*.yml")
	if err != nil {
		return fmt.Errorf("could not create playbook file: %w", err)
	}
	defer os.Remove(f.Name())

	fmt.Printf("Creating temporary playbook file in %s\n", f.Name())
	err = writePlaybook(roleName, hosts, f)
	if err != nil {
		return err
	}

	return executeAnsible(f.Name())
}

func writePlaybook(roleName string, hosts string, f *os.File) error {
	defer f.Close()

	var w io.Writer = f
	if *debug {
		w = &debugWriter{writer: w}
	}

	fmt.Println("Generating playbook content")
	p := playbook{
		Hosts:       []string{hosts},
		GatherFacts: *gatherFacts,
		Roles:       []string{roleName},
	}
	b, err := yaml.Marshal([]playbook{p})
	if err != nil {
		return fmt.Errorf("could not create playbook YAML: %w", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("could not write playbook to file", err)
	}

	return nil
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
