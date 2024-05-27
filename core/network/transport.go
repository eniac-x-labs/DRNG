package network

//
//import (
//	"context"
//	"log"
//	"sync"
//	"time"
//
//	_proto "DRNG/core/grpc/proto"
//	"google.golang.org/grpc"
//)
//
//type ServerAddress string
//
//type NetworkTransport struct {
//	connPool     map[ServerAddress][]*grpc.ClientConn
//	connPoolLock sync.Mutex
//
//	consumeCh chan RPC // 这个 chan
//
//	heartbeatFn     func(RPC)
//	heartbeatFnLock sync.Mutex
//
//	maxPool     int
//	maxInFlight int
//
//	serverAddressProvider ServerAddressProvider
//
//	shutdown     bool
//	shutdownCh   chan struct{}
//	shutdownLock sync.Mutex
//
//	stream StreamLayer
//
//	// streamCtx is used to cancel existing connection handlers.
//	streamCtx     context.Context
//	streamCancel  context.CancelFunc
//	streamCtxLock sync.RWMutex
//
//	timeout      time.Duration
//	TimeoutScale int
//
//	msgpackUseNewTimeFormat bool
//}
//
//// getExistingConn is used to grab a pooled connection.
//func (n *NetworkTransport) getPooledConn(target ServerAddress) *grpc.ClientConn {
//	n.connPoolLock.Lock()
//	defer n.connPoolLock.Unlock()
//
//	conns, ok := n.connPool[target]
//	if !ok || len(conns) == 0 {
//		return nil
//	}
//
//	var conn *grpc.ClientConn
//	num := len(conns)
//	conn, conns[num-1] = conns[num-1], nil
//	n.connPool[target] = conns[:num-1]
//	return conn
//}
//
//// getConn is used to get a connection from the pool.
//func (n *NetworkTransport) getConn(target ServerAddress) (*grpc.ClientConn, error) {
//	// Check for a pooled conn
//	if conn := n.getPooledConn(target); conn != nil { // 连接池
//		return conn, nil
//	}
//
//	conn, err := grpc.Dial(string(target), grpc.WithInsecure())
//	if err != nil {
//		log.Fatal("did not connect: %v", err)
//	}
//	//// Dial a new connection
//	//conn, err := n.stream.Dial(target, n.timeout)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//// Wrap the conn
//	//netConn := &netConn{
//	//	target: target,
//	//	conn:   conn,
//	//	dec:    codec.NewDecoder(bufio.NewReader(conn), &codec.MsgpackHandle{}),
//	//	w:      bufio.NewWriterSize(conn, connSendBufferSize),
//	//}
//	//
//	//netConn.enc = codec.NewEncoder(netConn.w, &codec.MsgpackHandle{
//	//	BasicHandle: codec.BasicHandle{
//	//		TimeNotBuiltin: !n.msgpackUseNewTimeFormat,
//	//	},
//	//})
//
//	// Done
//	return conn, nil
//}
//
//// returnConn returns a connection back to the pool.
//func (n *NetworkTransport) returnConn(conn *grpc.ClientConn) { // 连接还回来
//	n.connPoolLock.Lock()
//	defer n.connPoolLock.Unlock()
//
//	key := conn.Target()
//	conns := n.connPool[ServerAddress(key)]
//
//	if !n.IsShutdown() && len(conns) < n.maxPool {
//		n.connPool[ServerAddress(key)] = append(conns, conn)
//	} else {
//		conn.Close()
//	}
//}
//
//// CloseStreams closes the current streams.
//func (n *NetworkTransport) CloseStreams() {
//	n.connPoolLock.Lock()
//	defer n.connPoolLock.Unlock()
//
//	// Close all the connections in the connection pool and then remove their
//	// entry.
//	for k, e := range n.connPool {
//		for _, conn := range e {
//			conn.Close()
//		}
//
//		delete(n.connPool, k)
//	}
//
//	// Cancel the existing connections and create a new context. Both these
//	// operations must always be done with the lock held otherwise we can create
//	// connection handlers that are holding a context that will never be
//	// cancelable.
//	n.streamCtxLock.Lock()
//	n.streamCancel()
//	n.setupStreamContext()
//	n.streamCtxLock.Unlock()
//}
//
//func (n *NetworkTransport) GetGrpcClient(target ServerAddress) (_proto.RandomNumberGeneratorClient, error) {
//	conn, err := n.getConn(target)
//	if err != nil {
//		return nil, err
//	}
//	//conn, err := grpc.Dial(address, grpc.WithInsecure())
//	//if err != nil {
//	//	log.Fatal("did not connect: %v", err)
//	//}
//	return _proto.NewRandomNumberGeneratorClient(conn), nil
//}
