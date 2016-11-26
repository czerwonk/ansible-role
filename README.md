# go-ansible-role

This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.

## Requirements
* ansible must be installed and added to $PATH

## Use
Executing the foo role using ansible arguments for become:
```
> ansible-role foo -b --ask-user-pass
```
