This is a learning project. (SQL, Locks, Transactions) It implements a coupon giveaway API that can support concurrent users atomically.

The use case of this API is a situation where thousands of users are waiting to collect coupons. What happens if thousands of users try to collect the same limited coupon at the same time?

API has 3 endpoints to give coupons:

- /give/:couponId: First comes to mind solution. Wrong way.

- /give-transaction/:couponId: Gives coupons with a transaction created with database/sql package.

- give-transaction-pgfunction/:couponId: Gives coupons by calling a PostgreSQL function. Check the "give_coupon_to_user" function below.

Swagger Url: http://localhost:1323/swagger/index.html

You can use a simple load tester for testing, i used load-tester:

npm install -g load-tester

load-tester 5000

http://localhost:5000

Example Test Setup:

````
{
    "baseUrl": "http://localhost:1323",
    "duration": 15000,
    "connections": 100,
    "sequence": [
    { "method": "POST", "path": "/coupon/api/v1/give-transaction/d7513edf-511c-4515-a7e7-2920d50f1237"}
    ]
}
````

SQL For Tables:

`````
create table coupon
(
    id       uuid    not null
        constraint coupon_pk
            primary key,
    name     text,
    type     text    not null,
    quantity integer not null
);

create unique index coupon_id_uindex
    on coupon (id);
    
create table coupon_given_events
(
    couponid    text    not null,
    userid      text    not null,
    newquantity integer not null
);
`````

PostgreSQL Function Used By Api(give-transaction-pgfunction):
`````
CREATE OR REPLACE FUNCTION give_coupon_to_user(couponId uuid, userId text) RETURNS void AS
$$
declare
    quantity_of_coupon integer;
begin
    select quantity
    into quantity_of_coupon
    from coupon
    WHERE id = couponId
    FOR UPDATE;

    IF quantity_of_coupon >= 1 THEN
        UPDATE coupon
        SET quantity = quantity - 1
        WHERE id = couponId;

        INSERT INTO coupon_given_events(CouponId, UserId, NewQuantity)
        VALUES (couponId, userId, quantity_of_coupon - 1);
    ELSE
        RAISE EXCEPTION 'Coupon run out!';
    end if;


end;
$$ LANGUAGE plpgsql;
`````
