CREATE TABLE friends (
    user_id1 INTEGER NOT NULL,
    user_id2 INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    
    PRIMARY KEY (user_id1, user_id2),
    
    FOREIGN KEY (user_id1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id2) REFERENCES users(id) ON DELETE CASCADE,
    
    CHECK (user_id1 < user_id2) 
);