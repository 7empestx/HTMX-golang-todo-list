-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS tasks (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Description VARCHAR(255),
    Completed BOOLEAN, 
    AddedFrom VARCHAR(255) 
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE tasks;
