CREATE TABLE IF NOT EXISTS tasks (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    Description TEXT,
    Completed BOOLEAN,
    AddedFrom TEXT
);
