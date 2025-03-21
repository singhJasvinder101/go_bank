-- name: CreateAccount :one
insert into accounts(
    owner, balance, currency
) values (
    $1, $2, $3
) returning *;

-- name: GetAccountById :one
select * from accounts
where id = $1 limit 1;

-- name: GetAccountForUpdate :one
select * from accounts
where id = $1 limit 1
for no key update;

-- name: ListAccounts :many
select * from accounts
order by id
limit $1
offset $2;

-- name: UpdateAccountByID :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: UpdateAccountBalanceByID :one
update accounts
set balance = balance + sqlc.arg(amount)
where id = sqlc.arg(account_id)
RETURNING *;


-- name: DeleteAccountByID :exec
delete from accounts
where id = $1;