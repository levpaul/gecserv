#!/bin/bash

flatc -g --gen-object-api --gen-onefile --gen-all --go-namespace fb server.fbs
flatc -T --no-fb-import --gen-all --gen-onefile --short-names --es6-js-export -o client msg.fbs

cp client/msg_generated.ts ~/projects/idolscape/src/fb/
