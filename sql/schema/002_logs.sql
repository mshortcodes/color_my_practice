-- +goose Up
CREATE TABLE logs (
    id UUID PRIMARY KEY,
    date DATE NOT NULL,
    color_depth INT NOT NULL CHECK (color_depth BETWEEN 1 AND 5),
    confirmed BOOLEAN NOT NULL DEFAULT false, 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id, date)
);

-- +goose Down
DROP TABLE logs;