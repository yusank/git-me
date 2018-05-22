#!/usr/bin/env bash

user=root
prod_hosts=45.76.169.195
prod_app_dir=/data/www/golang/git-me/
prod_conf_file=conf/app.conf
prod_conf_dir=/data/www/golang/git-me/conf/
app=git-me


echo "cross compile"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

    echo "deal with $prod_hosts"
    echo "prepare"
    #ssh $user@$prod_hosts "mkdir -p $prod_app_dir"
    #ssh $user@$prod_hosts "mkdir -p $prod_conf_dir"

    echo "scp"
    scp -P 21212 $app $user@$prod_hosts:$prod_app_dir
    scp -P 21212 $prod_conf_file $user@$prod_hosts:$prod_conf_dir

rm -rf $app
echo "done"
