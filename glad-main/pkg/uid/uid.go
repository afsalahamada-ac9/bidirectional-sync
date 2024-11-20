/*
 * Copyright 2024 AboveCloud9.AI Products and Services Private Limited
 * All rights reserved.
 * This code may not be used, copied, modified, or distributed without explicit permission.
 */

package uid

import (
	"time"
)

var gSeq = uint64(0)

// TODO: Add epoch time & subtract that from the UnixMilli()

// Ref: https://instagram-engineering.com/sharding-ids-at-instagram-1cf5a71e5a5c
func Get(shard int) uint64 {

	msec := uint64(time.Now().UnixMilli())
	id := (msec << (64 - 41)) | (uint64(shard) << (64 - 41 - 13)) | (gSeq % 1024)
	gSeq++

	// Postgres doesn't support value with high bit set; so setting the high bit to zero
	mask := ^(uint64(1) << 63)
	return (id & mask)
}
