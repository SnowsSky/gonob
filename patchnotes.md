- gonob 1.0.0-dev-16 <br>
-> Bug Fixes, added --noconfirm. <br>
-> added remove.go <br>
-> Added alpm logfile & more. <br>

- gonob 1.0.0  <br>
-> Bug fixes <br>
-> Added install.go <br>
-> Better remove func + Added --noconfirm opt <br>
-> Added gonob list for non aur packages <br>
-> Added gonob local_install (gonob -U) for local install <br>
-> Fixed lock issues <br>
-> gonob install --aur now use gonob -U instead of pacman -U <br>
-> Added install.sh script <br>

- gonob 1.1.0-dev-1 <br>
-> by now, You can't launch AUR commands with sudo. <br>
-> Fixed + Improvements to the download bar <br>
-> the packages with "-debug" won't shows up anymore on `gonob list --aur` as a unknown package. <br>
-> Added `gonob -Ss` for official packages. <br>
-> Fixed the colors issues on local-install <br>
-> Added `gonob release_notes` to see release note of the version <br>
-> Improved translate.go (Loading files in mem instead of reading every time) <br>

- gonob 1.1.0 <br>
-> by now, You can't launch AUR commands with sudo. <br>
-> Fixed + Improvements to the download bar <br>
-> the packages with "-debug" won't shows up anymore on `gonob list --aur` as a unknown package. <br>
-> Added `gonob -Ss` for official packages. <br>
-> Fixed the colors issues on local-install <br>
-> Added `gonob release_notes` to see release note of the version <br>
-> Improved translate.go (Loading files in mem instead of reading every time) <br>
-> added release link on `gonob release_notes`. <br>
-> go-git for git usage. <br>

- gonob 2.0.0 <br>
-> New parser (so better parsing, better args usage, less bugs...) <br>
-> Added noconfirm for wrapper.install. <br>
-> Fixed len package to remove. <br>
-> Better help message <br>
-> bug fixes (especially the lock one) <br>

- gonob 2.1.0 <br>
-> Minors change + bug fixes  <br>
-> scolor usage for color <br>

- gonob 2.2.0 <br>
-> Better remove command, you can now delete any package (only if no conflicts) <br>
-> ALPM funcs like progresscallback are now in alpm.go <br>
-> New command : clean_cache, -Sc, removes the cache files.<br>

- gonob 2.3.0
-> Bug fixes (Like translation one)
-> New command : check_cache, -Sch, Checks the cache folder<br>
-> New command : update, -Sy, Syncs dabatases<br>
-> New command : upgrade, -Syu, Upgrade your entire system.<br>
-> Fix : Fixed cache folder.