# Ethereum Network Setup
This file contains instructions on how to get a local ethereum network running on your computer.<br />
**NB:- The instructions assume that the user is running a linux distro**

## Requirements 
1. None


## Installation

### Getting started
1.  `git clone https://github.com/victorkingi/year3_project_code.git` in a folder.

2. Go to https://geth.ethereum.org/downloads/https://geth.ethereum.org/downloads/ and download the latest linux version of Geth. Move the `tar.gz` file downloaded to `year3_project_code/geth/` folder so that the next script can easily access it.

3. Run `tar -xzvf  ../geth/geth-linux-amd64-{x.xx.xx-x}.tar.gz --directory ../geth --strip-components=1`. To extract the `geth` executable into `geth/` folder.

4. `cd` into `year3_project_code/setup` and run `./setup_env.sh`. If the error `bash: ./setup_env.sh: Permission denied` shows up, first run `chmod ugo+x setup_env.sh` then `./setup_env.sh`. Depending on your system, you might be required to prefix the `chmod` command with `sudo`.

5. During the script execution, you will be required to create a password for each of the accounts. Use `helloworld` as the password. If you would like to change it, you will have to change the script at `line 9`.<br /> `echo "MY_NEW_PASSWORD" >> password.txt;`.

4. Run `mvn install:install-file -Dfile="src/main/resources/lib/OpenPseudonymiserCryptoLib.jar" -DgroupId="com.open-pseudonymiser" -DartifactId="open-pseudonymiser" -Dversion="1.0.0" -Dpackaging="jar"` from the root source directory to install `OpenPseudonymiser.jar` as a `mvn local repository dependency`.
4. If you do not wish to compile the source code, find the jar file
   in releases.
Run `java -jar healthcare-data-simulators-x.x-SNAPSHOT.jar` from releases.

###### NB:- The jar files in step 2, 3 & 4 will be extracted from `main/resources/lib` folder to a `lib` folder for execution.



## Contributors
Vlad Andrei Bucur - vladbucur2000  
George Edward Nechitoaia - georgeedward2000 <br>
Victor Traistaru - Wyktorrr  <br>
Ena Balatinac - ennaena <br>
Victor Kingi - victorkingi 

## License
This software is being developed under the MIT License.