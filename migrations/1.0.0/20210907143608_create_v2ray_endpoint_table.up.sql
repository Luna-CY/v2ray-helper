create table if not exists v2ray_endpoint
(
    id             integer                not null
        constraint v2ray_endpoint_pk
        primary key autoincrement,
    cloud          integer                not null,
    endpoint       integer                not null,
    rate           text       default ''  not null,
    remark         text       default ''  not null,
    host           text                   not null,
    port           integer    default 443 not null,
    user_id        text                   not null,
    alter_id       integer    default 64  not null,
    level          integer    default 0   not null,
    transport_type integer    default 0   not null,
    web_socket     text       default ''  not null,
    deleted        integer(1) default 0   not null,
    create_time    integer                not null
);

create index idx_cloud_endpoint
    on v2ray_endpoint (cloud, endpoint);

create index idx_endpoint_cloud
    on v2ray_endpoint (endpoint, cloud);

