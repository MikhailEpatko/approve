-- +goose Up
-- +goose StatementBegin

create type status_enum as enum ('TEMPLATE', 'NEW', 'STARTED', 'FINISHED');
create type approve_type_enum as enum ('PARALLEL_ANY_OF', 'PARALLEL_ALL_OF', 'SEQUENTIAL_ALL_OFF');
create type decision_enum as enum ('UNKNOWN', 'ACCEPT', 'REJECT', 'REVISION');

create table if not exists route
(
  id          bigserial primary key,
  name        text                                  not null,
  description text,
  status      status_enum default 'TEMPLATE'        not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  deleted     boolean     default false
);

create table if not exists step_group
(
  id          bigserial primary key,
  route_id    bigint                                not null
    constraint fk_route_id
      references route,
  name        text                                  not null,
  number      int                                   not null,
  step_type   text                                  not null,
  status      status_enum default 'NEW'             not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  deleted     boolean     default false,
  unique (route_id, number)
);

create table if not exists step
(
  id            bigserial primary key,
  step_group_id bigint                                not null
    constraint fk_step_group_id
      references step_group,
  name          text                                  not null,
  number        int                                   not null,
  status        status_enum default 'NEW'             not null,
  approve_type  approve_type_enum                     not null,
  created_at    timestamptz default current_timestamp not null,
  modified_at   timestamptz default current_timestamp not null,
  deleted       boolean     default false,
  unique (step_group_id, number)
);

create table if not exists approver
(
  id          bigserial primary key,
  step_id     bigint                                not null
    constraint fk_step_id
      references step,
  guid        text                                  not null,
  name        text                                  not null,
  email       text                                  not null,
  number      int                                   not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  deleted     boolean     default false,
  unique (step_id, number)
);

create table if not exists resolution
(
  id          bigserial primary key,
  approver_id bigint                                  not null
    constraint fk_approver_id
      references approver,
  decision    decision_enum default 'UNKNOWN'         not null,
  comment     text,
  created_at  timestamptz   default current_timestamp not null,
  modified_at timestamptz   default current_timestamp not null,
  deleted     boolean       default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists resolution;
drop table if exists approver;
drop table if exists step;
drop table if exists step_group;
drop table if exists route;
drop type status_enum;
drop type approve_type_enum;
drop type decision_enum;
-- +goose StatementEnd
