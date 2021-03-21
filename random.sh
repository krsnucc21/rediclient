#!/bin/bash

PORT=8080
URL=a3e29c096919a4ee698c53c55ba2bf5a-895395087.us-east-1.elb.amazonaws.com:$PORT

MAXCOUNT=10
count=1

while [ "$count" -le $MAXCOUNT ]
do
  let "cellname = $RANDOM % 256"
  let "username = $RANDOM % 1000"
  let "number = $RANDOM % 100" 

  let "operation=$RANDOM % 5"
  echo $operation
  if [ "$operation" -lt 4 ]
  then
    curl -s -H "Content-type: application/json" -d '{"cellname": "'$cellname'", "username": "'$username'", "rsrp": '$number'}' $URL/rsrp
  else
    curl -s -N -H "Content-type: application/json" $URL/rsrp/$cellname | head -c 40
  fi

  let "count += 1"
done
