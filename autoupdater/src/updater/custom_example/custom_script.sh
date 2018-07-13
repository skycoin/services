#!/bin/bash

echo "service {$1}"
echo "version {$2}"
shift 2

echo "arguments {$@}"
