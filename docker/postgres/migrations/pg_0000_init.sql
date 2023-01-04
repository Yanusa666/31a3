create sequence public.authors_id_seq;
create table if not exists public.authors
(
    id   bigint default nextval('public.authors_id_seq') not null,
    name varchar                                         not null,
    constraint authors_pk
        primary key (id)
);


create sequence public.posts_id_seq;
create table if not exists public.posts
(
    id         bigint default nextval('public.posts_id_seq') not null,
    author_id  bigint                                        not null,
    title      varchar                                       not null,
    content    varchar                                       not null,
    created_at timestamptz                                   not null,
    constraint posts_pk
        primary key (id),
    constraint posts_author_id_fk
        foreign key (author_id) references public.authors
);