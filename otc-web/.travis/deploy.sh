#!/bin/bash
set -e # exit on first error

eval "$(ssh-agent -s)" # start ssh-agent cache
# id_rsa is decrypted as the first step of Travis build, see .travis.yml
chmod 600 ../.travis/id_rsa.deploy # allow read access to the private key
ssh-add ../.travis/id_rsa.deploy # add the private key to SSH

# prevent authenticity confirmations 
ssh-keyscan $IP >> ~/.ssh/known_hosts

# prepeare deployment
tar -czvf otc-ui.tar.gz build
scp otc-ui.tar.gz $RUN_USER@$IP:/home/apps/deploy-otc-ui

# start updated services
ssh $RUN_USER@$IP <<EOF
  cd /home/apps/deploy-otc-ui
  tar -zxvf otc-ui.tar.gz
  cp -r ./build/** /var/www/otc-ui
EOF
