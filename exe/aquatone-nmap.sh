#!/bin/bash

## equivalent to --ports huge. nmap is way quicker than aquatone-scan, there's no need to trim down the list
SCAN_PORTS="80,81,300,443,591,593,832,981,1010,1311,2082,2087,2095,2096,2480,3000,3128,3333,4243,4567,4711,4712,4993,5000,5104,5108,5800,6543,6379,7000,7396,7474,8000,8001,8008,8014,8042,8069,8080,8081,8088,8090,8091,8118,8123,8172,8222,8243,8280,8281,8333,8443,8500,8834,8880,8888,8983,9000,9043,9060,9080,9090,9091,9200,9443,9800,9981,1000,12443,16080,18091,18092,27018,20720,28017"
SSL_PORTS=",443,832,981,1010,1311,2083,2087,2096,4712,7000,8172,8243,8333,8443,9443,18091,18092,"

## check arguments
if [[ -z ${1} || -z ${2} || ${1} != '-d' ]]; then
    echo "[?] usage: ${0} -d <domain>"
    exit 1
fi

domain=${2}

## check destination folder
if [[ ! -d ~/aquatone/${domain} ]]; then
    echo "[!] ${domain} doesn't exist in aquatone folder"
    exit 1
fi

## begin
cd ~/aquatone/${domain}

## get uniq ips
echo "[>] sorting ips"
cut -f2 -d, hosts.txt \
    | sort -u \
    > hosts.nmap

## scan through nmap
echo "[+] nmap'ing $(wc -l hosts.nmap)..."
sudo nmap -iL hosts.nmap \
    -oG output.nmap \
    -Pn -g 53 -n -v \
    -sS -p ${SCAN_PORTS} --open \
    | grep 'Completed SYN Stealth Scan against'
sudo chmod 666 output.nmap

# sort through and make a list
echo "[>] building urls"
grep -Po 'Host: ([0-9]{1,3}\.){3}[0-9]{1,3}.*Ports:.*' output.nmap \
    | sed -r 's/Host: />/g' \
    | sed -r 's/\s+Ignored.*//g' \
    | grep -Po '(>([0-9]{1,3}\.){3}[0-9]{1,3}|[0-9]{2,5})' \
    | xargs \
    | tr -cs '[0-9].>' ',' \
    | tr '>' '\n' \
    | rev \
    | cut -f2- -d, \
    | rev \
    | sed '/^$/d' \
    > open_ports.txt

# build urls
while read -r line; do
    ip=${line%%,*}
    ports=${line#*,}
    ports=(${ports//,/ })

    # get hosts for ip
    hosts=($(grep "${ip}" hosts.txt | cut -f1 -d,))

    for host in ${hosts[@]}; do
        for port in ${ports[@]}; do
            if [[ ${port} == 80 ]]; then
                url="http://${host}/"
            elif [[ ${port} == 443 ]]; then
                url="https://${host}/"
            elif [[ ${SSL_PORTS} =~ ",${port}," ]]; then
                url="https://${host}:${port}/"
            else
                url="http://${host}:${port}/"
            fi
            echo ${url} >> urls.txt
        done
    done
done < open_ports.txt

echo "[!] all done"

