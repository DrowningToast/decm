-- name: CreateEventContract :one
INSERT INTO event_contracts (
    event_id,
    access_manager_contract_address,
    event_contract_address,
    ticket_contract_address,
    certificate_contract_address
) VALUES (
    sqlc.arg('event_id'),
    sqlc.arg('access_manager_contract_address'),
    sqlc.arg('event_contract_address'),
    sqlc.arg('ticket_contract_address'),
    sqlc.arg('certificate_contract_address')
) RETURNING *;

-- name: GetEventContractByEventID :one
SELECT * FROM event_contracts WHERE event_id = sqlc.arg('event_id');

-- name: UpdateEventContract :one
UPDATE event_contracts
SET 
    access_manager_contract_address = sqlc.arg('access_manager_contract_address'),
    event_contract_address = sqlc.arg('event_contract_address'),
    ticket_contract_address = sqlc.arg('ticket_contract_address'),
    certificate_contract_address = sqlc.arg('certificate_contract_address'),
    updated_at = NOW()
WHERE event_id = sqlc.arg('event_id')
RETURNING *;

-- name: DeleteEventContract :exec
DELETE FROM event_contracts WHERE event_id = sqlc.arg('event_id');
