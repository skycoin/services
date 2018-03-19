#!/bin/bash
set -e # exit on first error

eval "$(ssh-agent -s)" # start ssh-agent cache
# id_rsa is decrypted as the first step of Travis build, see .travis.yml
chmod 600 ../.travis/id_rsa.deploy # allow read access to the private key
ssh-add ../.travis/id_rsa.deploy # add the private key to SSH

# prevent authenticity confirmations 
ssh-keyscan $IP >> ~/.ssh/known_hosts

# prepeare deployment
mkdir deploy
cp otc ./deploy
tar -czvf otc.tar.gz deploy
scp otc.tar.gz $RUN_USER@$IP:/home/apps/deploy-otc

# start updated services
ssh $RUN_USER@$IP <<EOF
  cd /home/apps/deploy-otc
  tar -zxvf otc.tar.gz
  sudo systemctl stop otc
  cp -r ./deploy/** /usr/share/otc
  sudo systemctl start otc
EOF
