#!/usr/bin/env bash

rm -rf ../web/* | exit

mv ./dist/css ../web/ | exit
mv ./dist/fonts ../web/ | exit
mv ./dist/img ../web/ | exit
mv ./dist/js ../web/ | exit
mv ./dist/favicon.ico ../web/ | exit
mv ./dist/robots.txt ../web/ | exit
mv ./dist/index.html ../web/ | exit
