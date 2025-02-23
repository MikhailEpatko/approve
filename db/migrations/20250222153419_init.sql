-- +goose Up
-- +goose StatementBegin


create table if not exists route
(
  id          bigserial primary key,
  name        text not null,
  description text,
  status      text        default 'STARTED'         not null, -- template, started, finished
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  deleted boolean default false
);

create table if not exists step_group
(
  id          bigserial primary key,
  route_id    bigint                                not null
    constraint fk_route_id
      references route,
  name        text                                  not null,
  number      int                                   not null,
  step_type   text                                  not null, --parallel, sequential
  status      text        default 'STARTED'         not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  deleted boolean default false,
  unique (route_id, number)
);

create table if not exists step
(
  id              bigserial primary key,
  step_group_id   bigint                                not null
    constraint fk_step_group_id
      references step_group,
  approve_type   text                                   not null, -- parallel_any_of, parallel_all_of, sequential
  name            text                                  not null,
  number          int                                   not null,
  status          text        default 'STARTED'         not null,
  created_at      timestamptz default current_timestamp not null,
  modified_at     timestamptz default current_timestamp not null,
  deleted boolean default false,
  unique (step_group_id, number)
);

create table if not exists approver
(
  id          bigserial primary key,
  guid text                                not null,
  step_id     bigint                                not null
    constraint fk_step_id
    references step,
  full_name text                            not null,
  short_name text                            not null,
  email text                                not null,
  phone_number text                         not null,
  number      int                                   not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  deleted boolean default false,
  unique (step_id, number)
);

create table if not exists resolution
(
  id          bigserial primary key,
  approver_id     bigint                                not null
    constraint fk_approver_id
      references approver,
  type        text                                  not null,
  comment     text,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  deleted boolean default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists resolution;
drop table if exists approver;
drop table if exists step;
drop table if exists step_group;
drop table if exists route;
-- +goose StatementEnd
