#!/bin/sh
debug="$1"

export JOURNALIST_API_URL="http://127.0.0.1:8000/api/v1"

admin_user="admin"
admin_pass="admin"

user1_id=""
user1_user="user1"
user1_pass="p4sS!123456"

user2_id=""
user2_user="user2"
user2_pass="p4sS!123456"

feed1_url="http://lorem-rss.herokuapp.com/feed"
feed2_url="https://xn--gckvb8fzb.com"

failfast() {
  if [ "$1" -ne "0" ]
  then
    printf "   FAILED: %s\n" "$2"
    exit "$1"
  else
    printf "   SUCCESS\n"
    if [ "$debug" = "true" ]
    then
      printf "   DEBUG: %s\n" "$2"
    fi
  fi
}


#------------------------------------------------------------------------------#
printf "\
## Listing all users as admin \
\n"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform get \
       on users \
       as $admin_user $admin_pass
)
failfast $? "$out"



#------------------------------------------------------------------------------#
printf "\
## Creating user as admin with username %s \
\n" "$user1_user"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform post \
       on users \
       as $admin_user $admin_pass \
     with "{
       \"username\": \"$user1_user\",
       \"password\": \"$user1_pass\",
       \"role\": \"user\"
     }"
)
failfast $? "$out"
#------------------------------------------------------------------------------#



#------------------------------------------------------------------------------#
printf "\
## Listing all users as admin \
\n"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform get \
       on users \
       as $admin_user $admin_pass
)
failfast $? "$out"
user1_id="$(printf "%s" "$out" | jq --raw-output ".users[] | select(.username == \"$user1_user\") | .id")"
#------------------------------------------------------------------------------#



#------------------------------------------------------------------------------#
printf "\
## Updating %s as admin with 'admin' role \
\n" "$user1_user"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform put \
       on users/$user1_id \
       as $admin_user $admin_pass \
     with "{
       \"role\": \"admin\"
     }"
)
failfast $? "$out"
#------------------------------------------------------------------------------#



#------------------------------------------------------------------------------#
printf "\
## Updating %s as admin with 'user' role \
\n" "$user1_user"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform put \
       on users/$user1_id \
       as $admin_user $admin_pass \
     with "{
       \"role\": \"user\"
     }"
)
failfast $? "$out"
#------------------------------------------------------------------------------#



#------------------------------------------------------------------------------#
printf "\
## Creating user as admin with username %s \
\n" "$user2_user"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform post \
       on users \
       as $admin_user $admin_pass \
     with "{
       \"username\": \"$user2_user\",
       \"password\": \"$user2_pass\",
       \"role\": \"user\"
     }"
)
failfast $? "$out"
#------------------------------------------------------------------------------#



#------------------------------------------------------------------------------#
printf "\
## Creating token as %s with name 'mytoken' \
\n" "$user1_user"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform post \
       on tokens \
       as $user1_user $user1_pass \
     with "{
       \"name\": \"mytoken\"
     }"
)
failfast $? "$out"
#------------------------------------------------------------------------------#



#------------------------------------------------------------------------------#
printf "\
## Creating feed as %s with URL %s \
\n" "$user1_user" "$feed1_url"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform post \
       on feeds \
       as $user1_user $user1_pass \
     with "{
       \"name\": \"xn--gckvb8fzb.com\",
       \"url\": \"$feed1_url\",
       \"group\": \"Journals\"
     }"
)
failfast $? "$out"
#------------------------------------------------------------------------------#



#------------------------------------------------------------------------------#
printf "\
## Creating feed as %s with URL %s \
\n" "$user2_user" "$feed1_url"
# -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  - #
out=$(
  ./redacteur perform post \
       on feeds \
       as $user2_user $user2_pass \
     with "{
       \"name\": \"xn--gckvb8fzb.com\",
       \"url\": \"$feed1_url\",
       \"group\": \"Journals\"
     }"
)
failfast $? "$out"
#------------------------------------------------------------------------------#

