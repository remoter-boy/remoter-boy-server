create table tb_client
(
    client_id    varchar(60)             not null
        constraint tb_client_pk
            primary key,
    connected_at timestamp default now() not null,
    updated_at   timestamp
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
            on delete cascade
);

alter table tb_cpu
    owner to root;

create table tb_ram
(
    ram_uuid        varchar(36)      not null
        constraint tb_ram_pk
            primary key,
    ram_use_percent double precision not null,
    ram_use_gb      double precision not null,
    ram_total_gb    double precision not null,
    client_id       varchar(36)      not null
        constraint tb_ram_tb_client_client_id_fk
            references tb_client
            on delete cascade
);

alter table tb_ram
    owner to root;

create table tb_disk
(
    disk_uuid        varchar(36)      not null
        constraint tb_disk_pk
            primary key,
    disk_device      varchar(100)     not null,
    disk_mount_point varchar(100)     not null,
    disk_use_percent double precision not null,
    disk_use_gb      double precision not null,
    disk_total_gb    double precision not null,
    client_id        varchar(36)      not null
        constraint tb_disk_tb_client_client_id_fk
            references tb_client
            on delete cascade
);

alter table tb_disk
    owner to root;

create table tb_internet_protocol
(
    internet_protocol_uuid varchar(36) not null
        constraint tb_internet_protocol_pk
            primary key,
    v4                     varchar(25) not null,
    client_id              varchar(36) not null
        constraint tb_internet_protocol_tb_client_client_id_fk
            references tb_client
            on delete cascade
);

alter table tb_internet_protocol
    owner to root;

