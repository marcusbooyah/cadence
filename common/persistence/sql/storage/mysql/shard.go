// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package mysql

import (
	"database/sql"

	"github.com/uber/cadence/common/persistence/sql/storage/sqldb"
)

const (
	createShardQry = `INSERT INTO shards 
(shard_id, 
owner, 
range_id,
stolen_since_renew,
updated_at,
replication_ack_level,
transfer_ack_level,
timer_ack_level,
cluster_transfer_ack_level,
cluster_timer_ack_level,
domain_notification_version)
VALUES
(:shard_id, 
:owner, 
:range_id,
:stolen_since_renew,
:updated_at,
:replication_ack_level,
:transfer_ack_level,
:timer_ack_level,
:cluster_transfer_ack_level,
:cluster_timer_ack_level,
:domain_notification_version)`

	getShardQry = `SELECT
shard_id,
owner,
range_id,
stolen_since_renew,
updated_at,
replication_ack_level,
transfer_ack_level,
timer_ack_level,
cluster_transfer_ack_level,
cluster_timer_ack_level,
domain_notification_version
FROM shards WHERE
shard_id = ?
`

	updateShardQry = `UPDATE
shards 
SET
shard_id = :shard_id,
owner = :owner,
range_id = :range_id,
stolen_since_renew = :stolen_since_renew,
updated_at = :updated_at,
replication_ack_level = :replication_ack_level,
transfer_ack_level = :transfer_ack_level,
timer_ack_level = :timer_ack_level,
cluster_transfer_ack_level = :cluster_transfer_ack_level,
cluster_timer_ack_level = :cluster_timer_ack_level,
domain_notification_version = :domain_notification_version
WHERE
shard_id = :shard_id
`

	lockShardQry     = `SELECT range_id FROM shards WHERE shard_id = ? FOR UPDATE`
	readLockShardQry = `SELECT range_id FROM shards WHERE shard_id = ? LOCK IN SHARE MODE`
)

// InsertIntoShards inserts one or more rows into shards table
func (mdb *DB) InsertIntoShards(row *sqldb.ShardsRow) (sql.Result, error) {
	row.UpdatedAt = mdb.converter.ToMySQLDateTime(row.UpdatedAt)
	row.TimerAckLevel = mdb.converter.ToMySQLDateTime(row.TimerAckLevel)
	return mdb.conn.NamedExec(createShardQry, row)
}

// UpdateShards updates one or more rows into shards table
func (mdb *DB) UpdateShards(row *sqldb.ShardsRow) (sql.Result, error) {
	row.UpdatedAt = mdb.converter.ToMySQLDateTime(row.UpdatedAt)
	row.TimerAckLevel = mdb.converter.ToMySQLDateTime(row.TimerAckLevel)
	return mdb.conn.NamedExec(updateShardQry, row)
}

// SelectFromShards reads one or more rows from shards table
func (mdb *DB) SelectFromShards(filter *sqldb.ShardsFilter) (*sqldb.ShardsRow, error) {
	var row sqldb.ShardsRow
	err := mdb.conn.Get(&row, getShardQry, filter.ShardID)
	if err != nil {
		return nil, err
	}
	row.UpdatedAt = mdb.converter.FromMySQLDateTime(row.UpdatedAt)
	row.TimerAckLevel = mdb.converter.FromMySQLDateTime(row.TimerAckLevel)
	return &row, err
}

// ReadLockShards acquires a read lock on a single row in shards table
func (mdb *DB) ReadLockShards(filter *sqldb.ShardsFilter) (int, error) {
	var rangeID int
	err := mdb.conn.Get(&rangeID, readLockShardQry, filter.ShardID)
	return rangeID, err
}

// WriteLockShards acquires a write lock on a single row in shards table
func (mdb *DB) WriteLockShards(filter *sqldb.ShardsFilter) (int, error) {
	var rangeID int
	err := mdb.conn.Get(&rangeID, lockShardQry, filter.ShardID)
	return rangeID, err
}
