#!/bin/bash

# go build automatically excludes test files
go build -o bookings cmd/web/*.go && ./bookings