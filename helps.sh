#!/bin/bash
telepresence() {
	echo "@$*"
	go run ./cmd/telepresence "$@"
}

{
telepresence --help
telepresence status --help
telepresence version --help
telepresence quit --help
telepresence dashboard --help
telepresence intercept --help
telepresence leave --help
telepresence list --help
telepresence preview --help
telepresence preview create --help
telepresence remove create --help
telepresence connect --help
telepresence uninstall --help
telepresence login --help
telepresence logout --help
} | grep -e '^  telepresence' -e ^@
