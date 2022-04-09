# Ethereum Network Setup
This file contains instructions on how to get a local ethereum network running on your computer.<br />
**NB:- The instructions assume that the user is running a linux distro**

## Requirements 
1. Latest `git` version installed.
2. `tar` for extracting the geth zip.


## Getting Started
All required files and keys have been pre-generated as the setup process takes time. If you would like to generate your own nodes and ethereum addresses, follow the instructions below. If not, follow steps 1 and 2. Extract the zip file inside `geth/` directory using `tar -xzvf  "geth/${TAR_FILE_NAME}" --directory geth --strip-components=1` command. Run `start_bnode.sh` in one terminal, and run all `start_node_x.sh` commands, each, in separate terminals.

### Setup
1.  `git clone https://github.com/victorkingi/year3_project_code.git` in a folder.

2. Go to https://geth.ethereum.org/downloads/https://geth.ethereum.org/downloads/ and download the latest linux version of `geth & tools`. Move the `tar.gz` file downloaded to `year3_project_code/geth/` folder so that the next script can easily access it. You will need to run `mkdir geth` inside `year3_project_code/` incase `geth` folder doesn't exist. Make sure `geth` folder contains ONLY 1 file which is the zip file.

3. Run `./setup_env.sh`. If the error `bash: ./setup_env.sh: Permission denied` shows up, first run `chmod ugo+x setup_env.sh` then `./setup_env.sh`. Depending on your system priveleges, you might be required to prefix the `chmod` command with `sudo`.

4. During the script execution, you will be required to create a password for each node instance. Use `helloworld` as the password. If you would like to change it, you will have to change the script at `line 9`.<br /> `echo "MY_NEW_PASSWORD" >> password.txt;`. A new file `accounts.txt` should be created that contains all 5 node addresses in order.

5. Run `geth/puppeth` to create genesis file. type `year3project` as network name. Choose option 2. Configure new genesis. Choose option 1 Create new genesis from scratch. Choose option 2 Clique - proof-of-authority. Type 5 as block time. For accounts allowed to seal, copy the 5 nodes addresses in `accounts.txt` here pressing enter after each. For which accounts should be prefunded, enter the first 2 nodes addresses. Press enter under "Should the precompile-addresses be pre-funded". For the chain id, type `9984`. At this point, we now have a genesis file but we still need to export it. Under "What would you like to do?", choose 2 Manage existing genesis. Choose 2 Export genesis configuration. Type `ethereum` under which folder to save the genesis specs into. After this you can exit the terminal running `geth/puppeth`.


6. Run `initialise.sh` which will tell `geth` to use the genesis file  `year3project.json` and will also generate a boot node key.

7. To start the boot node, run `start_bnode.sh`. This will also generate the enode which will be stored in `enode.txt` file. The file will be used in the execution of the next script.

8. Open up a new terminal while the boot node is running and execute `create_start_scripts.sh` which will create `start_node_x.sh` files needed to start the nodes.

9. While in root directory, `year3_project_code`, run each `start_node_x.sh` script in a separate terminal to start all geth nodes. At this point, the nodes should start looking for peers, while submitting blocks and communicating with each other.