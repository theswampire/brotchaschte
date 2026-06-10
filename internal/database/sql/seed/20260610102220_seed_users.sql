-- +goose Up
-- +goose StatementBegin
-- password: 'password123'
INSERT INTO user VALUES('admin', '$argon2id$v=19$m=262144,t=3,p=1$JnPNwaqWhrPXqijevHUyCA$Wo8j2Ncr9MmNZYGXGYcQolzCw+z2/MlKTu7BQ0iMHEs');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'nothing';
-- +goose StatementEnd
