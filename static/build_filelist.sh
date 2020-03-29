#!/bin/bash
cat report_template.html | grep -E -o 'https://[^"]*(css|js)' > filelist.txt
