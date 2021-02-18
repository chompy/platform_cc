#!/usr/bin/env bash

# - vars
SSH_URL="$1"
PCC_PATH=~/.local/bin/pcc
PRIVATE_SSH_KEY_PATH=~/.ssh/id_rsa
RSYNC_PARAMS="-auve 'ssh -i /tmp/id_rsa' --include='*/' --include='*.jpg' --include='*.jpeg' --include='*.gif' --include='*.png' --include='*.webp' --exclude='*'"
if [[ "$OSTYPE" == "darwin"* ]]; then
    PCC_PATH=/usr/local/bin/pcc
fi

# - fetch code
GET_VARIABLES_CODE=$(cat <<END
import sys
import json
vars = json.load(sys.stdin)
for key in vars:
    print("%s;%s" % (key, vars[key]))
END
)
GET_DATABASE_CODE=$(cat <<END
import sys
import json
r = json.load(sys.stdin)
dbs = r.get("database", [])
for k in r:
    for rel in r[k]:
        if rel.get("scheme") != "mysql": continue
        print("%s;%d;%s;%s;%s" % (
            rel.get("host", ""), rel.get("port", 0), rel.get("username", ""), 
            rel.get("password", ""), rel.get("path", "")
        ))
END
)
GET_MOUNTS_CODE=$(cat <<END
import sys
import json
c = json.load(sys.stdin)
m = c.get("applications", [{}])[0].get("configuration", {}).get("mounts", {})
for k in m:
    print("%s;%s" % (k, m[k].get("souce_path", "")))
END
)

# - functions
display_error() {
    printf "\n\e[31mERROR:\e[0m\n$1\n\n"
    exit 1
}
get_variables() {
    echo "$1" | python3 -c "$GET_VARIABLES_CODE"
}
get_databases() {
    echo "$1" | python3 -c "$GET_DATABASE_CODE"
}
get_mounts() {
    echo "$1" | python3 -c "$GET_MOUNTS_CODE"
}

# - sanity checks
if [ -z $1 ]; then
    echo "USAGE: pcc_platform_sh_clone <ssh_url>"
    exit 0
fi
if [ ! -f $PCC_PATH ]; then
    display_error "Platform.CC not found at $PCC_PATH."
fi
command -v python3 >/dev/null 2>&1 || { display_error "Python 3 not found."; }
if [ ! -f $PRIVATE_SSH_KEY_PATH ]; then
    display_error "Private SSH key not found at $PRIVATE_SSH_KEY_PATH."
fi
if [ ! -f ".platform.app.yaml" ]; then
    display_error "Could not find .platform.app.yaml file."
fi
if [ "$(ssh $1 echo 'SUCCESS')" != "SUCCESS" ]; then
    display_error "Unable to establish SSH connection to $1."
fi

# - retreive data
echo "> Fetch data."
PLATFORM_RELATIONSHIPS=`ssh $SSH_URL 'echo $PLATFORM_RELATIONSHIPS | base64 -d'`
PLATFORM_VARIABLES=`ssh $SSH_URL 'echo $PLATFORM_VARIABLES | base64 -d'`
CONFIG_JSON=$(${PCC_PATH} project:configjson)

# - set variables
echo "> Set variables."
VAR_LIST=$(get_variables "$PLATFORM_VARIABLES")
if [[ $VAR_LIST ]]; then
    while IFS= read -r var; do
        IFS=';' read -ra var_split <<< "$var"
        $PCC_PATH var:set "${var_split[0]}" "${var_split[1]}"
    done <<< "$VAR_LIST"
else
    echo "    No variables found"
fi

# - start project
echo "> Start project."
$PCC_PATH project:start

# - get databse dumps (mysql)
echo "> Fetch database dumps."
DATABASE_LIST=$(get_databases "$PLATFORM_RELATIONSHIPS")
if [[ $DATABASE_LIST ]]; then
    while IFS= read -r db; do
        IFS=';' read -ra db_vals <<< "$db"
        IFS=
        if [ ! -z "${db_vals[4]}" ]; then
            echo "> Dump ${db_vals[4]}."
            ssh $SSH_URL "mysqldump --host=${db_vals[0]} --port=${db_vals[1]} --user=${db_vals[2]} --password=${db_vals[3]} ${db_vals[4]}" | gzip > /tmp/dump.sql.gz
            gunzip /tmp/dump.sql.gz
            echo "> Import ${db_vals[4]}."
            echo "drop schema if exists ${db_vals[4]}; create schema ${db_vals[4]}" | $PCC_PATH mysql:sql
            $PCC_PATH mysql:sql -d "${db_vals[4]}" < /tmp/dump.sql
            rm /tmp/dump.sql
        fi
    done <<< "$DATABASE_LIST"
else
    echo "    No databases found"
fi
# - upload ssh key to app
SSH_ID_RSA=`cat $PRIVATE_SSH_KEY_PATH`
$PCC_PATH app:sh "echo '$SSH_ID_RSA' > /tmp/id_rsa && chmod 0600 /tmp/id_rsa"

# - sync mounts
echo "> Fetch mounts."
MOUNT_LIST=$(get_mounts "$CONFIG_JSON")
MOUNT_SYNC_CMD=""
while IFS= read -r mount; do
    IFS=';' read -ra mount_split <<< "$mount"
    IFS=
    dest="${mount_split[0]}"
    src=$(echo "${mount_split[1]}" | sed "s/\:/_/")
    MOUNT_SYNC_CMD+="rsync $RSYNC_PARAMS $SSH_URL:/app$dest/ /mnt/data/$src/ || true && "
done <<< "$MOUNT_LIST"
MOUNT_SYNC_CMD+="true"
$PCC_PATH app:sh --root "$MOUNT_SYNC_CMD"

# - deploy hook
echo "> Deploy hook."
$PCC_PATH project:deploy