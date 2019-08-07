#!/usr/bin/env bash
SCRIPTS_DIR=$(dirname "${BASH_SOURCE[0]}")
set -x

# Make the clickhouse data volume root owned by the clickhouse user:
docker-compose run \
    --no-deps \
    --rm \
    --user root \
    --entrypoint chown \
    clickhouse \
    clickhouse:clickhouse /opt/clickhouse

docker-compose up -d zookeeper clickhouse
sleep 10
# Apply schema & migrations. Migrations are run sequentially in lexicographical order by filename:
for filename in $(ls ${SCRIPTS_DIR}/migrations/clickhouse | grep 'sql$'); do
    echo ${filename}
    docker-compose \
        --file docker-compose.yml \
        --file docker-compose-extra.yml \
        run \
        --rm \
        clickhouse-client \
        clickhouse-client \
        --host clickhouse \
        --multiquery < ${SCRIPTS_DIR}/migrations/clickhouse/${filename}
done
