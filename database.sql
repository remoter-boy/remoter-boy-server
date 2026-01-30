create table tb_client
(
    client_id    varchar(60)             not null
        constraint tb_client_pk
            primary key,
    connected_at timestamp default now() not null,
    updated_at   timestamp               not null
);

alter table tb_client
    owner to root;

create table tb_cpu
(
    cpu_uuid        varchar(36)      not null
        constraint tb_cpu_pk
            primary key,
    cpu_core_num    integer          not null,
    cpu_use_percent double precision not null,
    client_id       varchar(36)      not null
        constraint tb_cpu_tb_client_client_id_fk
            references tb_client
);

alter table tb_cpu
    owner to root;

create table tb_ram
(
    ram_uuid     varchar(36)      not null
        constraint tb_ram_pk
            primary key,
    used_percent double precision not null,
    used_gb      double precision not null,
    total_gb     double precision not null,
    client_id    varchar(36)      not null
        constraint tb_ram_tb_client_client_id_fk
            references tb_client
);

alter table tb_ram
    owner to root;

create table tb_disk
(
    disk_uuid    varchar(36)      not null
        constraint tb_disk_pk
            primary key,
    device       varchar(100)     not null,
    mount_point  varchar(100)     not null,
    used_percent double precision not null,
    used_gb      double precision not null,
    total_gb     double precision not null,
    client_id    varchar(36)      not null
        constraint tb_disk_tb_client_client_id_fk
            references tb_client
);

alter table tb_disk
    owner to root;

create table tb_internet_protocol
(
    internet_protocol_uuid varchar(36) not null
        constraint tb_internet_protocol_pk
            primary key,
    internet_protocol      varchar(25) not null,
    client_id              varchar(36) not null
        constraint tb_internet_protocol_tb_client_client_id_fk
            references tb_client
);

alter table tb_internet_protocol
    owner to root;

