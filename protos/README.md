# Protocol Buffers Schema Docs

*Last Updated: **August 8th 2021***  
*If you find any issues such as typos or inconsistencies in this documentation, please raise an Issue or submit a Pull Request.*

## Code Generation
Run the following command to generate Go files based on the Protocol Buffer specifications in the ``.proto`` files. This command requires the ``protoc`` and ``protoc-gen-go`` binaries to be installed.
```
make protos
```

## Schema Docs
The Protocol Buffers defined in this package are used for the p2p communication between peers on the Weave network.   
This file documents the various concepts associated with those various message schemas. 

### Message
A ``Message`` is the only type of message buffer that is propogated across the network. It acts as a wrapper for any underlying message type. The type of the underlying message is defined within the buffer and is constrained by the ``msgtype`` enum along with the message body which is one of ``Entity``, ``Query`` and ``Response``.

### Entity
An ``Entity`` is a message buffer that contains the data for a chain entity such as a transaction or a clock. The type of the entity is defined within the buffer and is constrained by the ``entitytype`` enum. The buffer also contains one of either the ``Block`` or ``Txn`` message types within it.

These messages are published on the network when either a new block is proposed or a new transaction is performed by the publishing peer.

#### Block
A ``Block`` is a message buffer that contains the data of a block for the blockchain. The buffer contains the hash of the block as bytes and the gob encoded bytes of the block data. This gob encoded data must decode to a ``weave.core.Block`` struct.

#### Txn
A ``Txn`` is a message buffer that contains the data of a transaction for the blockchain.
The buffer contains the hash of the transaction as bytes and the gob encoded bytes of the transaction data. This gob encoded data must decode to a ``weave.core.Transaction`` struct.

### Query
A ``Query`` is a message buffer that contains an underlying query message buffer. The type of query is defined within the buffer and is constrained by the ``querytype`` enum along with the query body which is one of ``TxnQuery``, ``BlockQuery``, ``StateQuery``, ``StatusQuery`` or ``InventoryQuery``.

These messages are published on the network as part of the networks state synchronization, such as when a new node joins the network and needs to update its local data to match the network state.

#### StateQuery
A ``StateQuery`` is a message buffer that contains the parameters for a chain state query. The buffer contains the peer ID of the peer that sent the query and a boolean indicating if the miner configuration state is also requested.

These messages are published when a peer is trying to determine the state of the network and chain in order to determine the necessary steps to synchronize its local state.

#### TxnQuery
A ``TxnQuery`` is a message buffer that contains the parameters for a transaction entity query. The buffer contains the peer ID of the peer that sent the query and the hash of requested transaction as bytes.

These messages are published when a peer is trying to build its mempool with in-flight transactions that occured prior to the node joining the network.

#### BlockQuery
A ``BlockQuery`` is a message buffer that contains the parameters for a block entity query. The buffer contains the peer ID of the peer that sent the query and the hash of requested block as bytes.

These messages are published when a peer is trying to update its outdated chain with blocks that were created and added to the chain while it was disconnected from the network.

#### InventoryQuery
An``InventoryQuery`` is a message buffer that contains the parameters for a node inventory query. The buffer contains the peer ID of the peer that sent the query. The query parameter include the height upto which block inventory must be included. This block height is the height of the chain locally on the requesting peer. The buffer also defines a boolean indicating if the transactions from the mempool of the node must be included in the inventory.

These messages are published when a peer is trying to determine the entities it needs to query in order to match the network state. The requesting peer directly queries the node with the best chain height based on its state query. 

#### StatusQuery
A ``StatusQuery`` is a message buffer that contains the parameters for an entity status query. The buffer contains the peer ID of the peer that sent the query, the type of the entity constrained by the ``entitytype`` enum and the hash of the entity.

These messages are published when a litenode is trying to determine the status of a transaction or a block.

### Response
A ``Response`` is a message buffer that contains an underlying query response buffer. The type of the response is defined within the buffer and is constrained by the ``responsetype`` enum along with the reponse body which is one of ``TxnResponse``, ``BlockResponse``, ``StateResponse``, ``StatusResponse`` or ``InventoryResponse``.

These messages are published on the network as part of the networks state synchronization, such as when a new node joins the network and needs to update its local data to match the network state.

#### StateResponse
A ``StateResponse`` is a message buffer that contains the response for a chain state query. The buffer contains the peer ID of the peer that sent the response along with the chain state on the responding peer. The chain state includes the current height of the chain and the miner configuration (included if included in the query).

The miner configuration is defined by the ``MinerConfig`` buffer that contains the size of the mempool (determines the number of transaction in a block) along with the current mining difficulty and reward.

These messages are published when a peer responds to another peer that is trying to determine the state of the network and chain and publishes a ``StateQuery``.

#### TxnQuery
A ``TxnResponse`` is a message buffer that contains the response for a transaction entity query. The buffer contains the peer ID of the peer that sent the response and the transaction data as a ``Txn`` buffer.

These messages are published when a peer responds to another peer is trying to build its mempool with in-flight transactions that occured prior to the node joining the network and publishes a ``TxnQuery``.

#### BlockResponse
A ``BlockResponse`` is a message buffer that contains the response for a block entity query. The buffer contains the peer ID of the peer that sent the response and the block data as a ``Block`` buffer.

These messages are published when a peer responds to another peer that is trying to update its outdated chain with blocks that were created and added to the chain while it was disconnected from the network and publishes a ``BlockQuery``.

#### InventoryResponse
An``InventoryResponse`` is a message buffer that contains the response for a node inventory query. The buffer contains the peer ID of the peer that sent the reponse along with the node inventory itself. The inventory includes a collection of blocks on the chain of the peer (included blocks are determined based on query) and a collection of transactions currently in the node mempool (included based on query).

These messages are published when a peer responds to another peer that is trying to determine the entities it needs to query in order to match the network state and publishes an ``InventoryQuery``. 

#### StatusResponse
A ``StatusResponse`` is a message buffer that contains the response for an entity status query. The buffer contains the peer ID of the peer that sent the response, the type of the entity constrained by the ``entitytype`` enum, the hash of the entity and the status of the entity.

These messages are published when a peer responds to another peer with a litenode is trying to determine the status of a transaction or a block.
