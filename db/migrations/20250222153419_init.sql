-- +goose Up
-- +goose StatementBegin
create table idm_user
(
  id           bigserial primary key,
  guid         text,
  full_name    text                                  not null,
  short_name   text                                  not null,
  email        text,
  phone_number text,
  birth_date   date                                  not null,
  created_at   timestamptz default current_timestamp not null,
  modified_at  timestamptz default current_timestamp not null
);

create table service
(
  id          bigserial primary key,
  code        text                                  not null unique,
  name        text                                  not null unique,
  description text,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null
);

create table role_type
(
  id          bigserial primary key,
  service_id  bigint
    constraint fk_service_id
      references service
      on delete cascade,
  code        text                                  not null unique,
  name        text                                  not null unique,
  description text,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  unique (service_id, code),
  unique (service_id, name)
);

create table role
(
  id          bigserial primary key,
  type_id     bigint                                not null
    constraint fk_role_type_id
      references role_type
      on delete cascade,
  code        text                                  not null,
  name        text                                  not null,
  description text,
  parent_id   bigint
    constraint fk_parent_id
      references role,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  unique (type_id, code),
  unique (type_id, name)
);

create table role_owner
(
  id          bigserial primary key,
  role_id     bigint                                not null
    constraint fk_role_id
      references role
      on delete cascade,
  idm_user_id bigint                                not null
    constraint fk_user_id
      references idm_user
      on delete cascade,
  number      int                                   not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null,
  unique (role_id, number)
);

create table role_to_user
(
  id          bigserial primary key,
  role_id     bigint                                not null
    constraint fk_role_id
      references role
      on delete cascade,
  idm_user_id bigint                                not null
    constraint fk_idm_user_id
      references idm_user
      on delete cascade,
  date_start  date                                  not null,
  date_finish date                                  not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null
);

-- approving
create table route
(
  id          bigserial primary key,
  name        text,
  status      text        default 'STARTED'         not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null
);

create table step_group
(
  id          bigserial primary key,
  route_id    bigint                                not null
    constraint fk_route_id
      references route
      on delete cascade,
  name        text                                  not null,
  number      int                                   not null,
  step_type   text                                  not null, -- parallel_any_of, parallel_all_of, sequential
  status      text        default 'STARTED'         not null,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null
);

create table step
(
  id              bigserial primary key,
  step_group_id   bigint                                not null
    constraint fk_step_group_id
      references step_group
      on delete cascade,
  role_owner_id   bigint                                not null
    constraint fk_role_owner_id
      references role_owner
      on delete cascade,
  name            text                                  not null,
  number          int                                   not null,
  status          text        default 'STARTED'         not null,
  resolution_type text                                  not null,
  created_at      timestamptz default current_timestamp not null,
  modified_at     timestamptz default current_timestamp not null
);

create table decision
(
  id          bigserial primary key,
  step_id     bigint                                not null
    constraint fk_step_id
      references step
      on delete cascade,
  type        text                                  not null,
  comment     text,
  created_at  timestamptz default current_timestamp not null,
  modified_at timestamptz default current_timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists idm_user;
drop table if exists service;
drop table if exists role_type;
drop table if exists role;
drop table if exists role_owner;
drop table if exists role_user;
drop table if exists route;
drop table if exists step_group;
drop table if exists step;
drop table if exists decision;
-- +goose StatementEnd
