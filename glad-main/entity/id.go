/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package entity

import (
	"math/rand"
	"strconv"
	"sudhagar/glad/pkg/uid"
)

// ID entity ID
type ID uint64

// NewID create a new entity ID
func NewID() ID {
	randomShardID := rand.Intn(1024)
	return ID(uid.Get(randomShardID))
}

// NewIDWithShard create a new entity ID with given Shard ID
func NewIDWithShard(shardID int) ID {
	return ID(uid.Get(shardID))
}

// StringToID convert a string to an entity ID
func StringToID(s string) (ID, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	return ID(id), err
}

func (id ID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}
