#!/bin/bash

INPUT_FILE=$1
DUMP_DIR=$2

if [ -z $INPUT_FILE ]; then
    echo "Usage: $0 <input_file> <dump_dir>"
    exit 1
fi

if [ -z $DUMP_DIR ]; then
    echo "Usage: $0 <input_file> <dump_dir>"
    exit 1
fi

# iterate over each line in the input file
while read line; do

    server=$(echo $line | cut -d ":" -f 1)

    # dump registry
    reggidump -v -t 30 -d $DUMP_DIR $line

    # search with truffle
    sudo trufflehog filesystem --only-verified --directory $DUMP_DIR/layer --archive-timeout=30s --concurrency=8 -j --no-update | tee $DUMP_DIR/$server/truffle.out

    # search for files with size >0 in dumped/layer
    for f in $(find $DUMP_DIR/layer -type f -size +0M);
    do
        rm $f
        touch $f
    done

done < $INPUT_FILE


exit 0