# Ethereum Network Setup
This file contains instructions on how to get a local ethereum network running on your computer.<br />
**NB:- The instructions assume that the user is running a linux distro**

## Requirements 
1. Latest `git` version installed.
2. `tar` for extracting the geth zip.


## Getting Started
All required files and keys have been pre-generated as the setup process takes time. If you would like to generate your own nodes and ethereum addresses, follow [Setup a new network](#setup-a-new-network). If not, follow [Run existing network](#run-existing-network).

### Run existing network
**NB:-** All this commands except for git clone should be run inside `year3_project_code/` directory.

1. `git clone https://github.com/victorkingi/year3_project_code.git` in a folder.

2. . Go to [Official Ethereum Download Page](https://geth.ethereum.org/downloads/) and download the latest linux version of `geth & tools`. Move the `tar.gz` file downloaded to `year3_project_code/geth/` folder so that the next script can easily access it. You will need to run `mkdir geth` inside `year3_project_code/` incase `geth` folder doesn't exist. Make sure `geth` folder contains ONLY 1 file which is the zip file.

3. Extract the zip file inside `geth/` directory using `tar -xzvf  "geth/${TAR_FILE_NAME}" --directory geth --strip-components=1` command.

4. Run `./initialise.sh` which will tell `geth` to use the genesis file  `year3project.json` and will also generate a boot node key.

5. Run `start_bnode.sh` in one terminal.

6. While in root directory, `year3_project_code`, run each `./start_node_{x}.sh` script in a separate terminal to start all geth nodes. At this point, the nodes should start looking for peers, while submitting blocks and communicating with each other.

### Setup a new network
**NB:-** All this commands except for git clone should be run inside `year3_project_code/` directory.

1.  `git clone https://github.com/victorkingi/year3_project_code.git` in a folder.

2. Run `./purge_all.sh` which will delete all pre-generated files and folders.

3. Go to [Official Ethereum Download Page](https://geth.ethereum.org/downloads/) and download the latest linux version of `geth & tools`. Move the `tar.gz` file downloaded to `year3_project_code/geth/` folder so that the next script can easily access it. You will need to run `mkdir geth` inside `year3_project_code/` incase `geth` folder doesn't exist. Make sure `geth` folder contains ONLY 1 file which is the zip file.

4. Run `./setup_env.sh`. If the error `bash: ./setup_env.sh: Permission denied` shows up, first run `chmod ugo+x setup_env.sh` then `./setup_env.sh`. Depending on your system priveleges, you might be required to prefix the `chmod` command with `sudo`.

5. During the script execution, you will be required to create a password for each node instance. Use `helloworld` as the password. If you would like to change it, you will have to change the script at `line 9`.<br /> `echo "MY_NEW_PASSWORD" >> password.txt;`. A new file `accounts.txt` should be created that contains all 5 node addresses in order.

6. Run `geth/puppeth` to create genesis file. type `year3project` as network name. Under "What would you like to do?", choose 2 "Manage existing genesis". Choose option 3 "Remove genesis configuration". Choose option 2 "Configure new genesis". Choose option 1 "Create new genesis from scratch". Choose option 2 "Clique - proof-of-authority". Type 5 as block time. For accounts allowed to seal, copy the 5 nodes addresses in `accounts.txt` here pressing enter after each (You will have to press enter twice to go to the next question). For which accounts should be prefunded, enter the first 2 nodes addresses. 
Press enter under "Should the precompile-addresses be pre-funded". For the chain id, type `9984`. At this point, we now have a genesis file but we still need to export it. Under "What would you like to do?", choose 2 "Manage existing genesis". Choose 2 "Export genesis configuration". Type `ethereum` under which folder to save the genesis specs into. After this you can exit the terminal running `geth/puppeth`.


7. Run `./initialise.sh` which will tell `geth` to use the genesis file  `year3project.json` and will also generate a boot node key.

8. To start the boot node, run `./start_bnode.sh`. This will also generate the enode which will be stored in `enode.txt` file. The file will be used in the execution of the next script.

9. Open up a new terminal while the boot node is running and execute `./create_start_scripts.sh` which will create `start_node_{x}.sh` files needed to start the nodes.

10. While in root directory, `year3_project_code`, run each `./start_node_{x}.sh` script in a separate terminal to start all geth nodes. At this point, the nodes should start looking for peers, while submitting blocks and communicating with each other.

### Common errors debugging

**Problem:** After running `./start_node_{x}.sh` you might get the error `Fatal: Error starting protocol stack: listen udp :30303: bind: address already in use` for udp or similar error in tcp.

**Solution:** Install netstat by running the following command depending on your distro:
```
$ sudo apt-get install net-tools    [On Debian/Ubuntu & Mint] 
$ sudo dnf install net-tools        [On CentOS/RHEL/Fedora and Rocky Linux/AlmaLinux]
$ pacman -S netstat-nat             [On Arch Linux]
$ emerge sys-apps/net-tools         [On Gentoo]
$ sudo dnf install net-tools        [On Fedora]
$ sudo zypper install net-tools     [On openSUSE]

```
Once installed, run:
```
$ netstat -ltnp
```
`-ltnp` represent the following options:
| Command option    | Description                   |
| :---              |       :---:                   |
| l                 | Only show listening sockets   |
| t                 | Display tcp connections       |
| u                 | Display udp connections       |
| n                 | Show numerical addresses      |
| p                 | Shows process ID and name     |

This will show all ports currently being used by the system. If they are a lot, you could use the command `netstat -l{t or u depending if a tcp or udp port was reported}np | grep -w ':{PORT_NUMBER_REPORTED}'` which will pipe the output into `grep` and `grep` will filter out the port number specified.
With this information, you could then check if it is an important process and make a decision to free the port or use a different port by changing `start_node_5.sh` code.<br /><br />
**NB:-** Only node 5 will run http and web socket api. 


