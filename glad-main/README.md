1. **HTTP-based health check endpoint**:
        ```go
        func main() {
            // ... other setup code ...
            
            http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
                // Add your health check logic here
                if healthy {
                    w.WriteHeader(http.StatusOK)
                    w.Write([]byte("healthy"))
                    return
                }
                w.WriteHeader(http.StatusServiceUnavailable)
            })
            
            // ... start server ...
        }
        ```
        
# PostgreSQL

## JSON types

PostgreSQL provides two JSON data types: JSON and JSONB. JSONB is recommended as it's more efficient for querying and indexing

```
CREATE TABLE example (
    id BIGSERIAL PRIMARY KEY,
    -- Regular JSON (stores exact copy, including whitespace)
    data_json JSON,
    
    -- Binary JSON (more efficient, supports indexing)
    data_jsonb JSONB,
    
    -- With NOT NULL constraint
    config JSONB NOT NULL DEFAULT '{}'::JSONB
);

-- Create index on JSONB field (optional)
CREATE INDEX idx_example_data ON example USING GIN (data_jsonb);
```

## JSON operations

- **Insert JSON data**:
    ```sql
    -- Insert JSON data
    INSERT INTO example (data_jsonb) VALUES (
        '{ 
            "name": "John",
            "address": {
                "street": "123 Main St",
                "city": "Boston"
            },
            "tags": ["developer", "golang"]
        }'
    );
    ```

- **Query JSON fields**:
    ```sql
    -- Query JSON fields
    SELECT 
        data_jsonb->>'name' as name,                    -- Get text
        data_jsonb->'address'->>'city' as city,         -- Nested object
        data_jsonb->'tags'->>0 as first_tag            -- Array element
    FROM example;
    ```

- **Query with conditions**:
    ```sql
    -- Query with conditions
    SELECT * FROM example 
    WHERE data_jsonb->>'name' = 'John';                -- Simple equality

    SELECT * FROM example 
    WHERE data_jsonb @> '{"tags": ["developer"]}';     -- Contains

    SELECT * FROM example 
    WHERE data_jsonb ? 'email';                        -- Has key
    ```

Key differences between `JSON` and `JSONB`:
1. `JSONB` is binary format (faster to process)
2. `JSONB` removes whitespace
3. `JSONB` can be indexed
4. `JSONB` removes duplicate keys

In Go, you can work with JSON fields like this:

```go
type Address struct {
    Street string `json:"street"`
    City   string `json:"city"`
}

type Record struct {
    ID      int64   `db:"id"`
    Name    string  `json:"name"`
    Address Address `json:"address"`
    Tags    []string `json:"tags"`
}

-- Using database/sql
var data []byte
err := db.QueryRow("SELECT data_jsonb FROM example WHERE id = $1", id).Scan(&data)
var record Record
json.Unmarshal(data, &record)

-- Using sqlx
var record Record
err := db.Get(&record, "SELECT data_jsonb FROM example WHERE id = $1", id)
```

