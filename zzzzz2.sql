-- database: /Users/samuele/Documents/wasa/WASA-1/WASA2DB.db
-- database: /Users/samuele/Documents/wasa/WASA-1/WASADB2.db

-- database: /Users/samuele/Documents/wasa/WASA-1/WASADB2.db
CREATE TABLE conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for the conversation
    name TEXT NOT NULL,                         -- Conversation name (e.g., group chat name)
    createdAt INTEGER NOT NULL,                 -- Timestamp of when the conversation was created (Unix time)
    isGroup BOOLEAN NOT NULL,                   -- Indicates if it's a group chat
    photo BLOB,                                 -- Photo for the conversation (e.g., group avatar)
    description TEXT                            -- Optional description field
    CHECK (
        ( isGroup = TRUE and description is not null and photo is not null )
         or 
        ( isGroup = FALSE and description is null and photo is null)
       
    )
);

CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for the message
    content TEXT,                               -- Text content of the message
    photoContent BLOB,                          -- Binary content (e.g., photo) of the message
    -- senderId INTEGER NOT NULL,                  -- Foreign key to the user table (who sent the message)
    sentAt INTEGER NOT NULL,                    -- Timestamp of when the message was sent (Unix time)
    conversationId INTEGER NOT NULL,            -- Foreign key to the conversation table
    -- answerTo INTEGER,
    -- CHECK (
        -- answerTo <> id
    -- )
    -- FOREIGN KEY (answerTo) REFERENCES messages(id)
    -- FOREIGN KEY (senderId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id)
);

CREATE TABLE participate (
    userId INTEGER NOT NULL,                -- Foreign key to the user table
    conversationId INTEGER NOT NULL,        -- Foreign key to the conversation table
    PRIMARY KEY (userId, conversationId),   -- Composite primary key
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (conversationId) REFERENCES conversations(id)
);

CREATE TABLE sent (
    userId INTEGER NOT NULL,                -- Foreign key to the user table
    messageId INTEGER NOT NULL,             -- Foreign key to the message table
    status TEXT NOT NULL,
    answerTo INTEGER,
    CHECK (
        answerTo <> messageId
    )
    CHECK (status IN ('unread', 'read', 'delivered')),
    FOREIGN KEY (answerTo) REFERENCES messages(id)
    PRIMARY KEY (userId, messageId),        -- Composite primary key
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (messageId) REFERENCES messages(id)
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,       -- Unique identifier for the user
    username TEXT NOT NULL,                     -- Username, must be unique (consider adding UNIQUE constraint)
    name TEXT,                                  -- First name of the user
    surname TEXT,                               -- Last name of the user
    photo BLOB                                  -- Photo stored as a binary large object
);