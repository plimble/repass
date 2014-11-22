#! /bin/bash -e

NAME=repass
INTERFACE=Store,Interface
PKG=github.com/plimble

mockgen \
  -destination=mock.go \
  --self_package=$PKG/$NAME \
  -package=$NAME \
  $PKG/$NAME \
  $INTERFACE

echo "OK"
