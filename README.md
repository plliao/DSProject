# Distributed System Final Project
NYU Tandon 2018 Distributed System Final Project

Instructor: Gustavo Sandoval

## Members
Name | NID | Email
--- | --- | ---
PEI-LUN LIAO | N18410090 | pll273@nyu.edu
I-TING CHEN | N19037964 | itc233@nyu.edu

## Verison

Assignment | Verison
--- | ---:
Part 1: Basic web app | 0.1.0 
Part 2: Separating Front End and Back End | 0.2.0 
Part 3: Replicated Back End | 0.3.1
Part 4: Project Demo | 

## Usage
1. Set up GOPATH

    `export GOPATH=$YOUR_PATH_TO_DSProject/DSProject`

1. Build web applications

    `go build frontEnd.go`
    
    `go build backEnd.go`

1. Serve the applications 

    Serve the backEnd server first.

    `./backEnd -id id_in_config -config config_file`
    
    `./frontEnd -port your_port -config config_file`

    See more instructions by
  
    `./backEnd -h`
    
    `./frontEnd -h`

## API
* Login and Signup

    The entry of the web application is on here with default port number 8811.
    
    `http://frontEnd:port/login/`
