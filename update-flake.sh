#!/usr/bin/env bash
go mod vendor
nix hash path vendor
rm -Rf vendor
