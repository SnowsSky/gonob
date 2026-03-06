# Welcome to a simple AUR helper written in go named gonob.
## gonob is a replacement for nob.
gonob is disponible on : 
    - fr_FR
    - us_US
# 1 - Downloading 🛜.
    - `git clone https://github.com/SnowsSky/gonob.git
# 2 - Installing
    I made a script to simplify the installatin process, just run : `chmod +x install.sh` and then `./install.sh`
# 3 - Documentation
    - 3.1 - How to use it ?
        - gonob is pretty simple, using `gonob --help` should be enough.
        - To install a package from the aur, run `gonob install / -S --aur <packages>`
        - To install a package from the official repos, run `gonob install / -S <packages>`
        - To install a package from unknown source, run `gonob local_install / -U <packagepath>`
        - To remove a package, run `gonob remove / -R <packages>`
        - To find a package from the aur, run `gonob search / -Ss --aur <package>`
        - To find a package from the localDB, run `gonob list / -Q | grep <package>`
        - To find a package from the official DBs, run `gonob search / -Ss <package>`
        - To check how many aur packages you have installed, run `gonob list / -Q --aur`