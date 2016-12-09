# ansible-role

This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.

## Requirements
* Ansible has to be installed and added to $PATH

## Install
```
go get github.com/czerwonk/ansible-role
```

## Use
Executing (privileged) a role named foo on each host in group servers:
```
ansible-role foo -b --ask-user-pass -l servers
```

## Ansible
see https://github.com/ansible/ansible
