#!/bin/bash

go build -o Bookings cmd/web/*.go && ./Bookings
./bookings -dbname=Bookings -dbuser=filipeparreiras -cache=false -production=false