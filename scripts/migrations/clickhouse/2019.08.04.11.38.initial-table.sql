
CREATE DATABASE IF NOT EXISTS leadpipe ON CLUSTER leadpipe_cluster;

CREATE TABLE IF NOT EXISTS leadpipe.hits_v1 ON CLUSTER leadpipe_cluster (
    hit_time Date,
    org_id Int64,
    visitor_id String,
    version UInt32,
    user_id Int64,
    page String,
    addr String,
    agent String
)
ENGINE ReplicatedReplacingMergeTree('/clickhouse/tables/{shard}/hits_v1', '{replica}', version)
PARTITION BY hit_time
ORDER BY (org_id, hit_time);

CREATE TABLE IF NOT EXISTS leadpipe.dist_hits_v1 ON CLUSTER leadpipe_cluster (
    hit_time Date,
    org_id Int64,
    visitor_id String,
    version UInt32,
    user_id Int64,
    page String,
    addr String,
    agent String
)
ENGINE Distributed(leadpipe_cluster, leadpipe, hits_v1, sipHash64(concat(toString(org_id), '|', toString(hit_time))))