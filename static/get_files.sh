#!/bin/bash
while IFS= read -r line
do
	(cd js_local_files && curl -s -O "$line")
done < filelist.txt
