# Welcome to a simple AUR helper written in go named gonob.
## gonob is a replacement for nob.
gonob is disponible on : <br>
    - fr_FR <br>
    - us_US <br>
# 1 - Downloading 🛜.
`git clone https://github.com/SnowsSky/gonob.git`
# 2 - Installing
I made a script to simplify the installatin process, just run : `chmod +x install.sh` and then `./install.sh`
# 3 - Documentation
## 3.1 - How to use it ?
gonob is pretty simple, using `gonob --help` should be enough. <br>
To install a package from the aur, run `gonob install / -S --aur <packages>` <br>
To install a package from the official repos, run `gonob install / -S <packages>` <br>
To install a package from unknown source, run `gonob local_install / -U <packagepath>` <br>
To remove a package, run `gonob remove / -R <packages>` <br>
To find a package from the aur, run `gonob search / -Ss --aur <package>` <br>
To find a package from the localDB, run `gonob list / -Q | grep <package>` <br>
To find a package from the official DBs, run `gonob search / -Ss <package>` <br>
To check how many aur packages you have installed, run `gonob list / -Q --aur` <br>