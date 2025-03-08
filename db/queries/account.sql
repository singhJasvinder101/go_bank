insert into accounts(
    owner, balance, currency
) values (
    $1, $2, $3
) returning *;
-- return all columns *