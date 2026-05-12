# ansible-role

This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.

## Requirements
* Ansible has to be installed and added to `$PATH`

## Install

### Using Nix
```
nix profile install github:czerwonk/ansible-role
```

### From source
```
cargo install --path .
```

## Use
Executing (privileged) a role named `foo` on each host in group `servers`:
```
ansible-role foo servers -b --ask-become-pass
```

## Notes
The playbook is passed to `ansible-playbook` via stdin (`/dev/stdin`), so no temporary file is written to disk. As a side effect, interactive password prompts (`--ask-pass`) are not supported — use SSH keys or `--vault-password-file` instead.

## Ansible
see https://github.com/ansible/ansible
