#!/bin/bash

set -e

rm -rf lib bin usr lib64 etc

for i in talk nudoku strace mv clear nnn nano touch mkdir locale tar gzip tree id bash vim grep less more echo cat ls hostname; do
        path=$(which $i)

        mkdir -p .$(dirname $path)
        cp -v $path .$path

        for f in `ldd $path | grep = | awk '{print $3}'`; do
                mkdir -p .$(dirname $f)
                cp $f .$f 
        done
done

for f in /lib/x86_64-linux-gnu/libnss_* /etc/services /etc/hosts /etc/nsswitch.conf /etc/resolv.conf; do
        mkdir -p .$(dirname $f)
        cp $f .$f 
done

for i in /lib/terminfo /usr/lib/locale; do
        mkdir -p .$i
        tar -cf - $i | tar -C . -xf -
done

echo ./lib64/ld-linux-x86-64.so.2
find . -type f | grep -v 'build.sh'
