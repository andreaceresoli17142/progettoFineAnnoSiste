#!/bin/bash

sudo pkill go 2> /dev/null

kitty --class="go-webdev" --hold exec `go run .; read i` &
