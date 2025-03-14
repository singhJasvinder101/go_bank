-- order matters alot

-- Tx1
begin;
update accounts set balance = balance - 100 where id = 13 returning *;
update accounts set balance = balance + 100 where id = 16 returning *;
commit;

-- Tx2
begin;
-- cause deadlock
-- update accounts set balance = balance - 100 where id = 16 returning *;
-- update accounts set balance = balance + 100 where id = 13 returning *;

-- avoid deadlock
update accounts set balance = balance + 100 where id = 13 returning *;
update accounts set balance = balance - 100 where id = 16 returning *;
commit;


