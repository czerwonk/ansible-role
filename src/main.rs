use std::io::Write;
use std::process::{Command, Stdio};

use anyhow::{Context, Result};
use serde::Serialize;

const VERSION: &str = "0.5.0";

struct Config {
    show_version: bool,
    force: bool,
    debug: bool,
    gather_facts: bool,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            show_version: false,
            force: false,
            debug: false,
            gather_facts: true,
        }
    }
}

#[derive(Serialize)]
struct Playbook {
    hosts: Vec<String>,
    gather_facts: bool,
    roles: Vec<String>,
}

fn main() {
    let raw_args: Vec<String> = std::env::args().skip(1).collect();
    let (config, positional) = parse_args(&raw_args);

    if config.show_version {
        print_version();
        std::process::exit(0);
    }

    if positional.len() < 2 {
        print_usage();
        std::process::exit(1);
    }

    let role_name = &positional[0];
    let hosts = &positional[1];
    let extra_args = &positional[2..];

    if hosts == "all" && !config.force {
        println!(
            "Executing roles for group all might be dangerous. Please add -force parameter to allow this action."
        );
        std::process::exit(1);
    }

    if let Err(e) = execute_role(role_name, hosts, extra_args, &config) {
        eprintln!("{e}");
        std::process::exit(2);
    }
}

fn parse_args(args: &[String]) -> (Config, Vec<String>) {
    let mut config = Config::default();
    let mut i = 0;

    while i < args.len() {
        match args[i].as_str() {
            "-version" | "--version" => config.show_version = true,
            "-force" | "--force" => config.force = true,
            "-debug" | "--debug" => config.debug = true,
            "-gather-facts" | "--gather-facts" => config.gather_facts = true,
            "-gather-facts=false" | "--gather-facts=false" => config.gather_facts = false,
            _ => break,
        }
        i += 1;
    }

    (config, args[i..].to_vec())
}

fn print_usage() {
    println!(
        "Usage: ansible-role [ -force ] $rolename $servers [ any further ansible-playbook parameters ]\n"
    );
    println!("  -debug");
    println!("\tPrints debug output");
    println!("  -force");
    println!("\tDisables safety checks");
    println!("  -gather-facts");
    println!("\tGather information of target hosts (default true)");
    println!("  -version");
    println!("\tShow version");
}

fn print_version() {
    println!("ansible-role");
    println!("Version: {VERSION}");
    println!("Author: Daniel Brendgen-Czerwonk");
}

fn execute_role(
    role_name: &str,
    hosts: &str,
    extra_args: &[String],
    config: &Config,
) -> Result<()> {
    println!("Role: {role_name}");
    println!("Generating playbook content");

    let playbook = Playbook {
        hosts: vec![hosts.to_string()],
        gather_facts: config.gather_facts,
        roles: vec![role_name.to_string()],
    };

    let yaml = serde_yaml::to_string(&[&playbook]).context("could not create playbook YAML")?;

    if config.debug {
        print!("{yaml}");
    }

    execute_ansible(&yaml, extra_args, config)
}

fn execute_ansible(playbook_yaml: &str, extra_args: &[String], config: &Config) -> Result<()> {
    println!("Starting ansible playbook");

    let mut ansible_args: Vec<&str> = extra_args.iter().map(String::as_str).collect();
    ansible_args.push("/dev/stdin");

    if config.debug {
        println!("ansible-playbook {ansible_args:?}");
    }

    let mut child = Command::new("ansible-playbook")
        .args(&ansible_args)
        .stdin(Stdio::piped())
        .stdout(Stdio::inherit())
        .stderr(Stdio::inherit())
        .spawn()
        .context("failed to start ansible-playbook")?;

    {
        let mut stdin = child.stdin.take().expect("stdin is piped");
        stdin
            .write_all(playbook_yaml.as_bytes())
            .context("could not write playbook to stdin")?;
    }

    child
        .wait()
        .context("failed to wait for ansible-playbook")?;

    Ok(())
}
