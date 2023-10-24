create user owndb password 'owndb_password' createdb;
create database own owner owndb;

create schema main authorization owndb;

select 1/0;

create table main.file_data (
    id integer primary key generated by default as identity,
    hash text not null,
    data_oid oid not null
)
;

create table main.folders (
     id integer primary key generated by default as identity,
     parent_folder_id integer,
     name text not null,
     dt_created timestamp not null default current_date,
     dt_changed timestamp
)
;

create table main.file_meta (
    id integer primary key generated by default as identity,
    file_data_id integer not null,
    folder_id integer,
    name text not null,
    extension text not null,
    original_path text not null,
    size bigint not null,
    dt_created timestamp not null default current_date,
    dt_changed timestamp,
    foreign key (file_data_id) references main.file_data(id), -- TODO: on delete cascade?
    foreign key (folder_id) references main.folders(id), -- TODO: on delete cascade?
    constraint check_input check (
        trim(extension) != ''
        and trim(original_path) != ''
        and size > 0
    )
)
;

insert into main.folders(id, parent_folder_id, name, dt_changed)
values (0, null, 'root', null);

CREATE OR REPLACE FUNCTION main.GetLargeObjectSize(oid) RETURNS bigint
    VOLATILE STRICT
    LANGUAGE 'plpgsql'
AS $$
DECLARE
    fd integer;
    sz bigint;
BEGIN
    -- Open the LO; N.B. it needs to be in a transaction otherwise it will close immediately.
    -- Luckily a function invocation makes its own transaction if necessary.
    -- The mode x'40000'::int corresponds to the PostgreSQL LO mode INV_READ = 0x40000.
    fd := lo_open($1, x'40000'::int);
    -- Seek to the end.  2 = SEEK_END.
    PERFORM lo_lseek(fd, 0, 2);
    -- Fetch the current file position; since we're at the end, this is the size.
    sz := lo_tell(fd);
    -- Remember to close it, since the function may be called as part of a larger transaction.
    PERFORM lo_close(fd);
    -- Return the size.
    RETURN sz;
END;
$$;
