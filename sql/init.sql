drop table if exists users cascade;
drop table if exists forums cascade;
drop table if exists threads cascade;
drop table if exists posts cascade;
drop table if exists votes cascade;
drop table if exists forums_users cascade;

drop function if exists upd_count_of_posts();
drop function if exists upd_count_of_threads();
drop function if exists change_edit_status();
drop function if exists update_vote_count();
drop function if exists update_vote_count_on_upd();

drop trigger if exists upd_post_count on posts;
drop trigger if exists upd_thread_count on threads;
drop trigger if exists change_status on posts;
drop trigger if exists update_vote on votes;
drop trigger if exists update_vote_on_upd on votes;

create table users
(
    nickname  varchar unique,
    full_name varchar not null,
    email     varchar unique,
    about     varchar,
    primary key (nickname, email)
);

create table forums
(
    posts   int     not null default 0,
    slug    varchar primary key,
    threads int     not null default 0,
    title   varchar not null,
    owner   varchar not null references users (nickname)
);

create table threads
(
    author  varchar                  not null references users (nickname),
    created timestamp with time zone not null default current_timestamp,
    forum   varchar references forums (slug),
    id      serial unique,
    message varchar                  not null,
    slug    varchar,
    primary key (id, slug),
    title   varchar                  not null,
    votes   int                               default 0
);

create table posts
(
    author   varchar not null references users (nickname),
    created  timestamp with time zone default current_timestamp,
    forum    varchar references forums (slug),
    id       serial primary key,
    isEdited bool                     default false,
    message  varchar not null,
    parent   int     not null         default 0,
    thread   int     not null references threads (id)
);

create table votes
(
    nickname varchar not null references users (nickname) on delete cascade,
    thread   int     not null references threads (id) on delete cascade,
    voice    int     not null check (voice = -1 or voice = 1),
    primary key (nickname, thread)
);

create table forums_users
(
    forum    varchar not null references forums (slug) on delete cascade,
    nickname varchar not null references users (nickname) on delete cascade,
    unique (nickname, forum)
);

create function upd_count_of_posts() returns trigger as
$$
begin
    update forums
    set posts = posts + 1
    where slug = new.forum;

    insert into forums_users
    values (new.forum, new.author)
    on conflict do nothing;

    return new;
end;
$$ LANGUAGE plpgsql;

create trigger upd_post_count
    after insert
    on posts
    for each row
execute procedure upd_count_of_posts();

create function upd_count_of_threads() returns trigger as
$$
begin
    update forums
    set threads = threads + 1
    where slug = new.forum;

    insert into forums_users
    values (new.forum, new.author)
    on conflict do nothing;

    return new;
end;
$$ LANGUAGE plpgsql;

create trigger upd_thread_count
    after insert
    on threads
    for each row
execute procedure upd_count_of_threads();

create function change_edit_status() returns trigger as
$$
begin
    if new.message <> old.message then
        new.isEdited = true;
    end if;

    return new;
end ;
$$ LANGUAGE plpgsql;

create trigger change_status
    before update
    on posts
    for each row
execute procedure change_edit_status();

create function update_vote_count() returns trigger as
$$
begin
    update threads
    set votes = votes + new.voice
    where id = new.thread;

    return new;
end;
$$ LANGUAGE plpgsql;

create trigger update_vote
    after insert
    on votes
    for each row
execute procedure update_vote_count();

create function update_vote_count_on_upd() returns trigger as
$$
begin
    update threads
    set votes = votes - old.voice + new.voice
    where id = new.thread;

    return new;
end;
$$ LANGUAGE plpgsql;

create trigger update_vote_on_upd
    before update
    on votes
    for each row
execute procedure update_vote_count_on_upd();
