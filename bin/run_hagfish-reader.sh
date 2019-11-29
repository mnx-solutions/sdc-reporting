#!/bin/bash

. /opt/local/etc/.env

if [ ! -e /var/log/usage/archive ]
then
    mkdir -p /var/log/usage/archive
fi

for usage_log in `find -name /var/log/usage/*.gz`
do
    /opt/local/bin/hagfish-reader ${usage_log}
    mv ${usage_log} /var/log/usage/archive/
done


#todo add cleanup of /var/log/usage/archive/

