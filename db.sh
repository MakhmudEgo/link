#!/bin/bash
if [[ "$1" == "redis" || "$1" == "postgres" ]]; then
cp Dockerfile Dockerfile.tmp
sed '$ d' Dockerfile.tmp > Dockerfile
echo "CMD [\"./main\", \"$1\"]" >> Dockerfile
rm -f Dockerfile.tmp
else
echo "No valid arg, write 'redis' or 'postgres'"
fi


