-- test here:
-- docker exec -it postgres psql -U postgres -d go_bank

begin;


-- step 1
insert into transfers (from_account_id, to_account_id, amount)
values (13, 16, 10)
returning *;
-- 


-- step 2
insert into entries (account_id, amount) 
values (13, -10)
returning *;

insert into entries (account_id, amount)
values (16, 10)
returning *;
-- 


-- step 3
select * from accounts where id = 13 for update;
update accounts set balance = balance - 10 where id = 13 returning *;
-- 



-- step 4
select * from accounts where id = 16 for update;
update accounts set balance = balance + 10 where id = 16 returning *;
-- 

rollback;


-- deadlocks caused by foreign key constraints