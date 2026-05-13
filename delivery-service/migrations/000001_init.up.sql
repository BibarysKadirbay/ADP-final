create extension if not exists "uuid-ossp";

create table couriers (
    id uuid primary key,
    user_id uuid not null unique,
    full_name varchar(180) not null,
    phone varchar(40) not null unique,
    vehicle_type varchar(40) not null,
    rating numeric(3,2) not null default 5 check (rating >= 0 and rating <= 5),
    total_deliveries integer not null default 0 check (total_deliveries >= 0),
    is_available boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create table deliveries (
    id uuid primary key,
    order_id uuid not null,
    courier_id uuid not null references couriers(id),
    restaurant_id uuid not null,
    customer_id uuid not null,
    status varchar(30) not null,
    pickup_address text not null,
    delivery_address text not null,
    estimated_eta_minutes integer not null default 0 check (estimated_eta_minutes >= 0),
    pickup_time timestamptz,
    delivered_time timestamptz,
    route_distance_km numeric(10,2) not null default 0 check (route_distance_km >= 0),
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    constraint deliveries_status_check check (status in ('pending','assigned','picked_up','on_the_way','delivered','cancelled'))
);

create table delivery_status_history (
    id uuid primary key,
    delivery_id uuid not null references deliveries(id) on delete cascade,
    old_status varchar(30),
    new_status varchar(30) not null,
    changed_at timestamptz not null default now(),
    constraint history_status_check check (new_status in ('pending','assigned','picked_up','on_the_way','delivered','cancelled'))
);

create table courier_ratings (
    id uuid primary key,
    courier_id uuid not null references couriers(id) on delete cascade,
    order_id uuid not null,
    customer_id uuid not null,
    rating integer not null check (rating >= 1 and rating <= 5),
    comment text not null default '',
    created_at timestamptz not null default now(),
    unique (courier_id, order_id, customer_id)
);

create index idx_couriers_available_rating on couriers(is_available, rating desc, total_deliveries asc);
create index idx_couriers_vehicle_available on couriers(vehicle_type, is_available);
create index idx_deliveries_order_id on deliveries(order_id);
create index idx_deliveries_courier_id on deliveries(courier_id);
create index idx_deliveries_restaurant_id on deliveries(restaurant_id);
create index idx_deliveries_customer_id on deliveries(customer_id);
create index idx_deliveries_status on deliveries(status);
create index idx_deliveries_created_at on deliveries(created_at desc);
create index idx_delivery_status_history_delivery_id on delivery_status_history(delivery_id, changed_at desc);
create index idx_courier_ratings_courier_id on courier_ratings(courier_id);
