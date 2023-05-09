// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

// RPCHeader is a common sub-structure used to pass along protocol version and
// other information about the cluster. For older Raft implementations before
// versioning was added this will default to a zero-valued structure when read
// by newer Raft versions.
type RPCHeader struct {
	// ProtocolVersion is the version of the protocol the sender is
	// speaking.
	ProtocolVersion ProtocolVersion
	// ID is the ServerID of the node sending the RPC Request or Response
	ID []byte
	// Addr is the ServerAddr of the node sending the RPC Request or Response
	Addr []byte
}

// WithRPCHeader is an interface that exposes the RPC header.
type WithRPCHeader interface {
	GetRPCHeader() RPCHeader
}

// AppendEntriesRequest is the command used to append entries to the
// replicated log.
/*
	用途:
		1.复制日志到follower中
		2.heartbeat
*/
type AppendEntriesRequest struct {
	RPCHeader

	// Provide the current term and leader
	Term uint64

	// Deprecated: use RPCHeader.Addr instead
	Leader []byte

	// Provide the previous entries for integrity checking
	/* TODO 为什么要提供PrevLogEntry和PrevLogTerm?
	 * raft中的机制：强制follower复制leader中的日志，就算follower日志和leader有冲突
	 * 会有两种情况：
	 * 		1.复制成功：follower会将PrevLogEntry和PrevLogTerm之后的日志替换成Entries中的日志
	 * 		2.复制失败：PrevLogEntry和PrevLogTerm没有能够和follower中日志匹配上的，leader需要往前继续找直到有日志和follower中对上的或者全量复制
	 */
	PrevLogEntry uint64
	PrevLogTerm  uint64

	// New entries to commit
	Entries []*Log

	// Commit index on the leader
	/* TODO leader发送的commitIndex是多少，follower就设置为多少,这里是为了保证follower的commitIndex和leader一致
	 * 		万一leader的commitIndex比follower的commitIndex小，follower会不会出现问题？
	 *
	 */
	// leader告诉follower LeaderCommitIndex之前的日志都可以应用到自己的状态机中
	LeaderCommitIndex uint64
}

// GetRPCHeader - See WithRPCHeader.
func (r *AppendEntriesRequest) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}

// AppendEntriesResponse is the response returned from an
// AppendEntriesRequest.
type AppendEntriesResponse struct {
	RPCHeader

	// Newer term if leader is out of date
	Term uint64

	// Last Log is a hint to help accelerate rebuilding slow nodes
	LastLog uint64

	// We may not succeed if we have a conflicting entry
	Success bool

	// There are scenarios where this request didn't succeed
	// but there's no need to wait/back-off the next attempt.
	NoRetryBackoff bool
}

// GetRPCHeader - See WithRPCHeader.
func (r *AppendEntriesResponse) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}

// RequestVoteRequest is the command used by a candidate to ask a Raft peer
// for a vote in an election.
type RequestVoteRequest struct {
	RPCHeader

	// Provide the term and our id
	Term uint64

	// Deprecated: use RPCHeader.Addr instead
	Candidate []byte

	// Used to ensure safety
	/* a想得到b的投票，那么以下两种情况b会直接拒绝a的请求
	 *	1.a的lastLogTerm < b的lastLogTerm
	 *	2.a的lastLogTerm = b的lastLogTerm && a的lastLogIndex < b的lastLogIndex
	 */
	// 这两个字段保证了选举出的leader拥有最新的日志
	LastLogIndex uint64
	LastLogTerm  uint64

	// Used to indicate to peers if this vote was triggered by a leadership
	// transfer. It is required for leadership transfer to work, because servers
	// wouldn't vote otherwise if they are aware of an existing leader.
	LeadershipTransfer bool
}

// GetRPCHeader - See WithRPCHeader.
func (r *RequestVoteRequest) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}

// RequestVoteResponse is the response returned from a RequestVoteRequest.
type RequestVoteResponse struct {
	RPCHeader

	// Newer term if leader is out of date.
	Term uint64

	// Peers is deprecated, but required by servers that only understand
	// protocol version 0. This is not populated in protocol version 2
	// and later.
	Peers []byte

	// Is the vote granted.
	Granted bool
}

// GetRPCHeader - See WithRPCHeader.
func (r *RequestVoteResponse) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}

// InstallSnapshotRequest is the command sent to a Raft peer to bootstrap its
// log (and state machine) from a snapshot on another peer.
type InstallSnapshotRequest struct {
	RPCHeader
	SnapshotVersion SnapshotVersion

	Term   uint64
	Leader []byte

	// These are the last index/term included in the snapshot
	LastLogIndex uint64
	LastLogTerm  uint64

	// Peer Set in the snapshot.
	// but remains here in case we receive an InstallSnapshot from a leader
	// that's running old code.
	// Deprecated: This is deprecated in favor of Configuration
	Peers []byte

	// Cluster membership.
	Configuration []byte
	// Log index where 'Configuration' entry was originally written.
	ConfigurationIndex uint64

	// Size of the snapshot
	Size int64
}

// GetRPCHeader - See WithRPCHeader.
func (r *InstallSnapshotRequest) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}

// InstallSnapshotResponse is the response returned from an
// InstallSnapshotRequest.
type InstallSnapshotResponse struct {
	RPCHeader

	Term    uint64
	Success bool
}

// GetRPCHeader - See WithRPCHeader.
func (r *InstallSnapshotResponse) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}

// TimeoutNowRequest is the command used by a leader to signal another server to
// start an election.
type TimeoutNowRequest struct {
	RPCHeader
}

// GetRPCHeader - See WithRPCHeader.
func (r *TimeoutNowRequest) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}

// TimeoutNowResponse is the response to TimeoutNowRequest.
type TimeoutNowResponse struct {
	RPCHeader
}

// GetRPCHeader - See WithRPCHeader.
func (r *TimeoutNowResponse) GetRPCHeader() RPCHeader {
	return r.RPCHeader
}
