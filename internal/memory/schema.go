package memory

const schemaSQL = `
CREATE TABLE IF NOT EXISTS memories (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    agent TEXT DEFAULT '',
    created_at TEXT NOT NULL,
    embedding BLOB
);
CREATE INDEX IF NOT EXISTS idx_memories_category ON memories(category);
CREATE INDEX IF NOT EXISTS idx_memories_agent ON memories(agent);
`
