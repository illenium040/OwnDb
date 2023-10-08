create user owndb password 'owndb_password' createdb;
create database own owner owndb;

create schema main authorization owndb;

select 1/0;

create table main.file_meta (
    id integer primary key generated by default as identity,
    file_data_id integer not null,
    name text not null,
    extension char(32) not null,
    original_path text not null,
    size bigint not null,
    dt_created date not null,
    dt_changed date,
    foreign key (file_data_id) references main.file_data(id),
    constraint check_input check (
        extension != ''
        and original_path != ''
        and size > 0
    )
)
;

create table main.file_data (
    id integer primary key generated by default as identity,
    hash text not null,
    data_oid oid not null
)
;

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
