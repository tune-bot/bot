#!/bin/bash

sed -i "/#\$nrconf{restart} = 'i';/s/.*/\$nrconf{restart} = 'a';/" /etc/needrestart/needrestart.conf

apt update
apt upgrade -y
apt install -y golang

mkdir -p bin

echo "#!/bin/bash" > bin/discord
echo "source vars/discord.env" >> bin/discord
echo "cd discord && git stash && git pull && go run ." >> bin/discord
chmod a+rx bin/discord