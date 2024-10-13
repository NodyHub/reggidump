#!/bin/bash

DUMP_DIR=$1

if [ -z $DUMP_DIR ]; then
    echo "Usage: $0 <dump_dir>"
    exit 1
fi

if [ ! -d $DUMP_DIR ]; then
    echo "Dump directory does not exist"
    exit 1
fi

# convert truffle output
for f in $(find $DUMP_DIR -name truffle.out); do echo "[" > $f.json ; awk 'NF{print $0 ","}' $f >> $f.json ; echo "{}" >> $f.json; echo "]" >> $f.json ; done

# print results
cat $DUMP_DIR/*/truffle.out.json | jq '.[] | .DetectorName + " " + .Raw ' | sort -u | grep -v '" "'

# Summary
echo "Lines: $(cat $DUMP_DIR/*/truffle.out.json | jq '.[] | .DetectorName + " " + .Raw ' | sort -u | grep -v '" "'| wc -l)"


exit 0
