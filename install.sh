#!/bin/bash


echo "gonob - Install Script"

echo "1 - Install / Update / Reinstall gonob"
echo "2 - Build only"
echo "3 - Remove gonob"
read -p "Make your choice: " opt

build() {
    echo "Building gonob..."
    start=$(date +%s)
    go build .
    if [ $? -ne 0 ]; then
        exit 1
    fi
    end=$(date +%s)
    echo "Build finished in $((end - start)) seconds ⏱️"
}

install() {
    echo "Installing gonob..."
    sudo install -Dm755 ./gonob /usr/bin/gonob
    if [ $? -ne 0 ]; then
        exit 1
    fi

    sudo mkdir -p /etc/gonob/translations
    if [ $? -ne 0 ]; then
        exit 1
    fi

    install_translations
    # Install license
    echo "Installing license..."
    sudo install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
    if [ $? -ne 0 ]; then
       exit 1
    fi
}

install_translations() {
    echo "Installing translations..."
    for file in ./translations/*.json; do
        sudo cp "$file" /etc/gonob/translations
        if [ $? -ne 0 ]; then
           exit 1
        fi
    done
}

if [ "$opt" == "1" ]; then
    echo "Updating local clone..."
    git pull -f -n
    echo "Installing Build dependencies..."
    sudo pacman -S go --noconfirm --needed
    if [ $? -ne 0 ]; then
        exit 1
    fi

    echo "Installing runtime dependencies..."
    sudo pacman -S git base-devel fakeroot debugedit --noconfirm --needed
    if [ $? -ne 0 ]; then
        exit 1
    fi

    build

    install
    gonob --version
    exit 0

elif [ "$opt" == "2" ]; then 
    echo "Updating local clone..."
    git pull -f -n
    build
    echo "Do you want to install translations ? (NEEDED to execute gonob) y/n "
    read -p "Make your choice: " opt
    if [ "$opt" == "y" ]; then
        install_translations
    fi 
    ./gonob --version
elif [ "$opt" == "3" ]; then
    echo "Removing executable..."
    sudo rm /usr/bin/gonob
    if [ $? -ne 0 ]; then
        exit 1
    fi

    echo "Removing gonob directory..."
    sudo rm -rf /etc/gonob
    if [ $? -ne 0 ]; then
        exit 1
    fi

    echo "gonob has been successfully removed."
    exit 0
else
    echo "Invalid choice. Please select 1 or 2."
    exit 1
fi
