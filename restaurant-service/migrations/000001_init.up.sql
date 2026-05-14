create extension if not exists "uuid-ossp";

create table restaurants (
    id uuid primary key,
    owner_id uuid not null,
    name varchar(180) not null,
    description text not null default '',
    cuisine_type varchar(80) not null,
    address text not null default '',
    city varchar(120) not null,
    rating numeric(3,2) not null default 0 check (rating >= 0 and rating <= 5),
    total_reviews integer not null default 0 check (total_reviews >= 0),
    image_url text not null default '',
    is_open boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table menu_categories (
    id uuid primary key,
    restaurant_id uuid not null references restaurants(id) on delete cascade,
    name varchar(120) not null,
    created_at timestamptz not null default now(),
    unique (restaurant_id, name)
);

create table menu_items (
    id uuid primary key,
    category_id uuid not null references menu_categories(id) on delete cascade,
    restaurant_id uuid not null references restaurants(id) on delete cascade,
    name varchar(180) not null,
    description text not null default '',
    price numeric(12,2) not null check (price >= 0),
    image_url text not null default '',
    is_available boolean not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index idx_restaurants_owner_id on restaurants(owner_id);
create index idx_restaurants_city on restaurants(city);
create index idx_restaurants_cuisine_type on restaurants(cuisine_type);
create index idx_restaurants_rating on restaurants(rating desc, total_reviews desc);
create index idx_restaurants_open_city on restaurants(city, is_open);
create index idx_restaurants_name_search on restaurants using gin (to_tsvector('simple', name || ' ' || description));
create index idx_menu_categories_restaurant_id on menu_categories(restaurant_id);
create index idx_menu_items_restaurant_id on menu_items(restaurant_id);
create index idx_menu_items_category_id on menu_items(category_id);
create index idx_menu_items_available on menu_items(restaurant_id, is_available);
