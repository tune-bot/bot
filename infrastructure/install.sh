#!/bin/bash

sed -i "/#\$nrconf{restart} = 'i';/s/.*/\$nrconf{restart} = 'a';/" /etc/needrestart/needrestart.conf

apt update
apt upgrade -y
apt install -y golang

mkdir -p bin

echo "#!/bin/bash" > bin/bot
echo "source vars/discord.env" >> bin/bot
echo "bin/download -U" >> bin/bot
echo "cd discord && git stash && git pull && go run ." >> bin/bot
chmod a+rx bin/bot