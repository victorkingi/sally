# Ethereum Network Setup
This file contains instructions on how to get a local ethereum network running on your computer.<br />
**NB:- The instructions assume that the user is running a linux distro**

## Requirements 
1. None


## Installation

### Getting started
1.  `git clone https://github.com/victorkingi/year3_project_code.git` in a folder.

2. Go to https://geth.ethereum.org/downloads/https://geth.ethereum.org/downloads/ and download the latest linux version of `geth & tools`. Move the `tar.gz` file downloaded to `year3_project_code/geth/` folder so that the next script can easily access it. You will need to run `mkdir geth` inside `year3_project_code/` incase `geth` folder doesn't exist.

3. Run `tar -xzvf  ../geth/geth-alltools-linux-amd64-{x.xx.xx-x}.tar.gz --directory ../geth --strip-components=1`. To extract the `geth` executable into `geth/` folder.

4. `cd` into `year3_project_code/setup` and run `./setup_env.sh`. If the error `bash: ./setup_env.sh: Permission denied` shows up, first run `chmod ugo+x setup_env.sh` then `./setup_env.sh`. Depending on your system, you might be required to prefix the `chmod` command with `sudo`.

5. During the script execution, you will be required to create a password for each node instance. Use `helloworld` as the password. If you would like to change it, you will have to change the script at `line 9`.<br /> `echo "MY_NEW_PASSWORD" >> password.txt;`.

6. When the above script exists successfully, you can now run `initialise.sh` which will tell `geth` to use the genesis file  `year3project.json`.

7. Create folder inside `ethereum` called `bnode` and run `../../geth/bootnode -genkey boot.key`. To start the boot node, run `../../geth/bootnode -nodekey "boot.key" -verbosity 7 -addr "127.0.0.1:30301"` while still inside `ethereum/bnode/` folder. An enode string will be generated in the terminal with this syntax `enode://xx@127.0.0.1:0?discport=30301`. Copy it and update the `start_node_x.sh` scripts with it incase they are not the same.

8. While in root directory, `year3_project_code`, run each `start_node_x.sh` script in a separate terminal to start all geth nodes.