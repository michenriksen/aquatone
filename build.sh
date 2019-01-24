#!/bin/bash

BUILD_FOLDER=build
VERSION=$(cat core/banner.go | grep Version | cut -d '"' -f 2)

bin_dep() {
  BIN=$1
  which $BIN > /dev/null || { echo "[-] Dependency $BIN not found !"; exit 1; }
}

create_exe_archive() {
  bin_dep 'zip'

  OUTPUT=$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" aquatone.exe ../README.md ../LICENSE.txt > /dev/null
  rm -rf aquatone aquatone.exe
}

create_archive() {
  bin_dep 'zip'

  OUTPUT=$1

  echo "[*] Creating archive $OUTPUT ..."
  zip -j "$OUTPUT" aquatone ../README.md ../LICENSE.txt > /dev/null
  rm -rf aquatone aquatone.exe
}

build_linux_amd64() {
  echo "[*] Building linux/amd64 ..."
  GOOS=linux GOARCH=amd64 go build -o aquatone ..
}

build_macos_amd64() {
  echo "[*] Building darwin/amd64 ..."
  GOOS=darwin GOARCH=amd64 go build -o aquatone ..
}

build_windows_amd64() {
  echo "[*] Building windows/amd64 ..."
  GOOS=windows GOARCH=amd64 go build -o aquatone.exe ..
}

rm -rf $BUILD_FOLDER
mkdir $BUILD_FOLDER
cd $BUILD_FOLDER

build_linux_amd64 && create_archive aquatone_linux_amd64_$VERSION.zip
build_macos_amd64 && create_archive aquatone_macos_amd64_$VERSION.zip
build_windows_amd64 && create_exe_archive aquatone_windows_amd64_$VERSION.zip
shasum -a 256 * > checksums.txt

echo
echo
du -sh *

cd --
