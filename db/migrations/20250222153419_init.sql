-- +goose Up
-- +goose StatementBegin

create type status_enum as enum ('TEMPLATE', 'NEW', 'STARTED', 'FINISHED');
create type order_type_enum as enum ('PARALLEL_ANY_OF', 'PARALLEL_ALL_OF', 'SERIAL');

create table if not exists route
(
  id          bigint primary key generated always as identity,
  name        text                           not null,
  description text                           not null,
  status      status_enum default 'TEMPLATE' not null,
  is_approved boolean     default false      not null,
  unique (name)
);

create table if not exists step_group
(
  id          bigint primary key generated always as identity,
  route_id    bigint                         not null
    constraint fk_route_id
      references route
      on delete cascade,
  name        text                           not null,
  number      int                            not null,
  step_order  order_type_enum                not null,
  status      status_enum default 'TEMPLATE' not null,
  is_approved boolean     default false      not null,
  unique (route_id, number)
);

create table if not exists step
(
  id             bigint primary key generated always as identity,
  step_group_id  bigint                         not null
    constraint fk_step_group_id
      references step_group
      on delete cascade,
  name           text                           not null,
  number         int                            not null,
  status         status_enum default 'TEMPLATE' not null,
  approver_order order_type_enum                not null,
  is_approved    boolean     default false      not null,
  unique (step_group_id, number)
);

create table if not exists approver
(
  id       bigint primary key generated always as identity,
  step_id  bigint                         not null
    constraint fk_step_id
      references step
      on delete cascade,
  guid     text                           not null,
  name     text                           not null,
  position text                           not null,
  email    text                           not null,
  number   int                            not null,
  status   status_enum default 'TEMPLATE' not null,
  unique (step_id, number)
);

create table if not exists resolution
(
  id          bigint primary key generated always as identity,
  approver_id bigint          not null
    constraint fk_approver_id
      references approver
      on delete cascade,
  is_approved boolean         not null,
  comment     text default '' not null,
  unique (approver_id)
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
drop type order_type_enum;
-- +goose StatementEnd
