package memory

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// UpsertFact inserts or updates a durable fact.
func (m *MemoryManager) UpsertFact(ctx context.Context, fact Fact) error {
	fact.Scope = strings.TrimSpace(fact.Scope)
	fact.EntityID = strings.TrimSpace(fact.EntityID)
	fact.Key = strings.TrimSpace(fact.Key)
	fact.Value = strings.TrimSpace(strings.ReplaceAll(fact.Value, "\x00", ""))
	fact.Source = strings.TrimSpace(fact.Source)

	if fact.Scope == "" || fact.EntityID == "" || fact.Key == "" || fact.Value == "" {
		return fmt.Errorf("scope, entity_id, key and value are required")
	}
	if fact.Source == "" {
		fact.Source = "unknown"
	}

	query := `
		INSERT INTO memory_facts (scope, entity_id, key, value, source, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(scope, entity_id, key) DO UPDATE SET
			value = excluded.value,
			source = excluded.source,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := m.db.ExecContext(ctx, query, fact.Scope, fact.EntityID, fact.Key, fact.Value, fact.Source)
	if err != nil {
		return fmt.Errorf("failed to upsert fact: %w", err)
	}
	return nil
}

// GetFact retrieves a single durable fact.
func (m *MemoryManager) GetFact(ctx context.Context, scope, entityID, key string) (Fact, bool, error) {
	query := `
		SELECT id, scope, entity_id, key, value, source, updated_at
		FROM memory_facts
		WHERE scope = ? AND entity_id = ? AND key = ?
		LIMIT 1
	`

	var fact Fact
	err := m.db.QueryRowContext(ctx, query, strings.TrimSpace(scope), strings.TrimSpace(entityID), strings.TrimSpace(key)).
		Scan(&fact.ID, &fact.Scope, &fact.EntityID, &fact.Key, &fact.Value, &fact.Source, &fact.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Fact{}, false, nil
		}
		return Fact{}, false, fmt.Errorf("failed to query fact: %w", err)
	}

	return fact, true, nil
}

// ListFacts returns all facts for a scope/entity pair ordered by most recent update first.
func (m *MemoryManager) ListFacts(ctx context.Context, scope, entityID string) ([]Fact, error) {
	query := `
		SELECT id, scope, entity_id, key, value, source, updated_at
		FROM memory_facts
		WHERE scope = ? AND entity_id = ?
		ORDER BY updated_at DESC, id DESC
	`
	rows, err := m.db.QueryContext(ctx, query, strings.TrimSpace(scope), strings.TrimSpace(entityID))
	if err != nil {
		return nil, fmt.Errorf("failed to query facts: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var facts []Fact
	for rows.Next() {
		var fact Fact
		if err := rows.Scan(&fact.ID, &fact.Scope, &fact.EntityID, &fact.Key, &fact.Value, &fact.Source, &fact.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan fact row: %w", err)
		}
		facts = append(facts, fact)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return facts, nil
}
