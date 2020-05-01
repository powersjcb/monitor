-- up

create table metrics (
    source varchar(253) not null, -- max length of a valid hostname
    ts timestamp,
    inserted_at timestamp not null,
    name varchar(100) not null,
    value float
);

create index source_ts_index on metrics (source, ts);

