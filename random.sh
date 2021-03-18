#!/bin/bash

URL=ac8d637fcb2b54eb49ea69420aa799ae-2044562349.us-east-1.elb.amazonaws.com:$PORT

MAXCOUNT=10
count=1

while [ "$count" -le $MAXCOUNT ]
do
  let "cellname = $RANDOM % 256"
  let "username = $RANDOM % 1000"
  let "number = $RANDOM % 100" 

  let "operation=$RANDOM % 2"
  echo $operation
  if [ "$operation" -eq 0 ]
  then
    curl -s -H "Content-type: application/json" -d '{"cellname": "'$cellname'", "username": "'$username'", "rsrp": '$number'}' $URL/rsrp
  else
    curl -s -N -H "Content-type: application/json" $URL/rsrp/$cellname | head -c 40
  fi

  let "count += 1"
done
