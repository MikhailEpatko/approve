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
  deleted     boolean     default false
);

create or replace rule soft_delete_route as
  on delete to route
  where old.status in ('STARTED', 'FINISHED')
  do instead
  update route set deleted = true where id = old.id;

create table if not exists step_group
(
  id          bigint primary key generated always as identity,
  route_id    bigint                         not null
    constraint fk_route_id
      references route,
  name        text                           not null,
  number      int                            not null,
  step_order  order_type_enum                not null,
  status      status_enum default 'TEMPLATE' not null,
  is_approved boolean     default false      not null,
  deleted     boolean     default false      not null,
  unique (route_id, number)
);

create or replace rule soft_delete_step_group as
  on delete to step_group
  where old.status in ('STARTED', 'FINISHED')
  do instead
  update step_group set deleted = true where id = old.id;

create table if not exists step
(
  id             bigint primary key generated always as identity,
  step_group_id  bigint                         not null
    constraint fk_step_group_id
      references step_group,
  name           text                           not null,
  number         int                            not null,
  status         status_enum default 'TEMPLATE' not null,
  approver_order order_type_enum                not null,
  is_approved    boolean     default false      not null,
  deleted        boolean     default false      not null,
  unique (step_group_id, number)
);

create or replace rule soft_delete_step as
  on delete to step
  where old.status in ('STARTED', 'FINISHED')
  do instead
  update step set deleted = true where id = old.id;

create table if not exists approver
(
  id       bigint primary key generated always as identity,
  step_id  bigint                         not null
    constraint fk_step_id
      references step,
  guid     text                           not null,
  name     text                           not null,
  position text                           not null,
  email    text                           not null,
  number   int                            not null,
  status   status_enum default 'TEMPLATE' not null,
  deleted  boolean     default false      not null,
  unique (step_id, number)
);

create or replace rule soft_delete_approver as
  on delete to approver
  where old.status in ('STARTED', 'FINISHED')
  do instead
  update approver set deleted = true where id = old.id;

create table if not exists resolution
(
  id          bigint primary key generated always as identity,
  approver_id bigint                not null
    constraint fk_approver_id
      references approver,
  is_approved boolean               not null,
  comment     text    default ''    not null,
  deleted     boolean default false not null
);

create or replace rule soft_delete_resolution as
  on delete to resolution do instead
  update resolution
  set deleted = true
  where id = old.id;
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
