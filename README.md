# go-ansible-role

This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.

## Requirements
* ansible must be installed and added to $PATH

## Install
```
go install ansible-role.go
```

## Use
Executing a role named foo with ansible arguments for become:
```
ansible-role foo -b --ask-user-pass
```

## Ansible
see https://github.com/ansible/ansible
