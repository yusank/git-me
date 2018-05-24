#!/bin/bash
##Author by yusank

function batch_chmod() {
	for file in `ls $1`
	do
		if [ -d $1"/"$file ]
		then 
			batch_chmod $1"/"$file
		else
            if [ "${file##*.}"x = "go"x ];
            then
                echo $1"/"$file
                (cat tempL | cat - $1"/"$file > file0) && mv file0 $1"/"$file
            fi
		fi
	done
}
batch_chmod $1
