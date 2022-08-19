#!/bin/sh
api_url="http://127.0.0.1:8000/api/v1"

admin_user="admin"
admin_pass="admin"

user1_user="user1"
user1_pass="p4sS!"

perform() {
  action=$(printf "%s" "$1" | tr '[:lower:]' '[:upper:]')
  # "on"=$2
  endpoint="$3"
  # "as"=$4
  user="$5"
  pass="$6"
  # "with"=$7
  payload="$8"

  if [ "$payload" = "" ]
  then
    json=$(curl -s -u "$user:$pass" -H 'Content-Type: application/json; charset=utf-8' -X "$action" "$api_url/$endpoint")
  else
    json=$(curl -s -u "$user:$pass" -H 'Content-Type: application/json; charset=utf-8' -X "$action" "$api_url/$endpoint" -d "$payload")
  fi

  if [ $? -ne 0 ]
  then
    printf "{}"
    return 1
  fi

  printf "%s" "$json"
  if printf "%s" "$json" | jq '.success' | grep "false" > /dev/null
  then
    return 2
  fi

  return 0
}

failfast() {
  if [ "$1" -ne "0" ]
  then
    exit "$1"
  fi
}

out=$(
  perform get \
       on users \
       as $admin_user $admin_pass
)
failfast $?
out=$(
  perform post \
       on users \
       as $admin_user $admin_pass \
     with "{
       \"username\": \"$user1_user\",
       \"password\": \"$user1_pass\",
       \"role\": \"user\"
     }"
)
failfast $?
out=$(
  perform get \
       on users \
       as $admin_user $admin_pass
)
failfast $?

user1_id="$(printf "%s" "$out" | jq ".users[] | select(.username == \"$user1_user\") | .id")"
echo "$user1_id"

