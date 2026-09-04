# Welcome to a simple AUR helper written in go named gonob.
## gonob is a replacement for nob.
gonob is available on : <br>
    - fr_FR <br>
    - en_US <br>
# 1 - Downloading 🛜.
`git clone https://github.com/SnowsSky/gonob.git`
# 2 - Installing
# 2.1 - Building from source : 
I made a script to simplify the installatin process, just run : `chmod +x install.sh` and then `./install.sh`
# 2.2 - Use my own pacman repository (recommended):
read : https://github.com/SnowsSky/SNrep
# 3 - Documentation
## 3.1 - How to use it ?
gonob is pretty simple, using `gonob --help` should be enough. <br>
To install a package from the aur, run `gonob install / -S --aur <packages>` <br>
To install a package from the official repos, run `gonob install / -S <packages>` <br>
To install a package from unknown source, run `gonob local_install / -U <packagepath>` <br>
To remove a package, run `gonob remove / -R <packages>` <br>
To find a package from the aur, run `gonob search / -Ss --aur <package>` <br>
To find a package from the localDB, run `gonob list / -Q | grep <package>` <br>
To find a package from the aur, run `gonob list / -Q --aur | grep <package>` <br>
To find a package from the official DBs, run `gonob search / -Ss <package>` <br>
To check how many aur packages you have installed, run `gonob list / -Q --aur` <br>
To sync official databases, run `gonob update / -Sy` <br>
To upgrade official packages, run `gonob upgrade / -Su`. Please note that you'll need to run `gonob update / -Sy` before <br>
To only upgrade aur packages, run `gonob upgrade / -Su --aur` <br>
To full upgrade the system (aur + databases), run `gonob -Syu` <br>
