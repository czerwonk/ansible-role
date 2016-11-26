# go-ansible-role

This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.

## Requirements
* ansible must be installed and added to $PATH

## Install
```
go install ansible-role.go
```

## Use
Executing (privileged) a role named foo on each host in group servers:
```
ansible-role foo -b --ask-user-pass -l servers
```

## Ansible
see https://github.com/ansible/ansible
