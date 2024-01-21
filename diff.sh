#!/bin/bash

find . -type f | while read f
do echo `basename $f`,`cat $f | wc -l`,`diff -yw --left-column ../prompts-gpt/$f $f | grep -cv '($'`
done