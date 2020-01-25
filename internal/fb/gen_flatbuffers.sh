#!/bin/bash

flatc -g --gen-object-api  --gen-onefile --gen-all --go-namespace fb  msg.fbs
flatc -T --no-fb-import --gen-all --short-names msg.fbs

cp msg_generated.ts ~/projects/idolscape/src/fb/
