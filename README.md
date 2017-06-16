# ansible-role [![Build Status](https://travis-ci.org/czerwonk/ansible-role.svg)][travis]
[![Go Report Card](https://goreportcard.com/badge/github.com/czerwonk/ansible-role)][goreportcard]

This is a simple wrapper for Ansible to run a single role without the need to generate a playbook first.

## Requirements
* Ansible has to be installed and added to $PATH

## Install
```
go get -u github.com/czerwonk/ansible-role
```

## Use
Executing (privileged) a role named foo on each host in group servers:
```
ansible-role foo servers -b --ask-pass --ask-become-pass
```

## Ansible
see https://github.com/ansible/ansible

[travis]: https://travis-ci.org/czerwonk/ansible-role
[goreportcard]: https://goreportcard.com/report/github.com/czerwonk/ansible-role
