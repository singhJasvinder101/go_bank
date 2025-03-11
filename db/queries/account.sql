-- name: CreateAccount :one
insert into accounts(
    owner, balance, currency
) values (
    $1, $2, $3
) returning *;

-- name: GetAccountById :one
select * from accounts
where id = $1 limit 1;

-- name: ListAccounts :many
select * from accounts
order by id
limit $1
offset $2;

-- name: UpdateAccountByID :exec
update accounts 
set balance = $2
where id = $1;


-- name: DeleteAccountByID :exec
delete from accounts
where id = $1;