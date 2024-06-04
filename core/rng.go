package core

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"DRNG/common"
	"DRNG/common/blsSignatures"
	_errors "DRNG/common/errors"
	"DRNG/common/eth"
	"DRNG/common/types"
	_conf "DRNG/config"
	_grpc "DRNG/core/grpc"
	_proto "DRNG/core/grpc/proto"
	"DRNG/core/random"
	_db "DRNG/db"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RNG struct {
	_proto.UnimplementedRandomNumberGeneratorServer
	ctx context.Context
	wg  *sync.WaitGroup

	//config *_conf.Config

	nodeID string // local
	leader string
	//Generators         []string // include self
	role common.RoleType

	nodeInfoFromConfig   map[string]_conf.NodeInfo
	unConnectedNodeTable map[string]struct{}

	// nodes information record
	nodeURLTable       map[string]string                             // key=nodeID, value=url
	nodeKeyTable       map[string]blsSignatures.PublicKey            //*ecdsa.PublicKey
	nodeRpcClientTable map[string]_proto.RandomNumberGeneratorClient // remote:client
	connectionMap      map[string]*grpc.ClientConn

	// BLS signature
	privKey      blsSignatures.PrivateKey // for bls sign
	blsThreshold int

	// random number messages from remote nodes
	generatorRecords GeneratorRandomNumberRecord
	leaderRecords    LeaderRandomNumberRecord
	rpcCh            chan interface{} // grpc msg push into rpcCh, main loop listen rpcCh

	errCh          chan types.ErrMessage
	grpcServerPort string
	pendingSeq     atomic.Int64
	//leaderReqs     types.LeaderReqs

	lastDiscardSeq  int64
	discardDuration time.Duration

	// every random number request has a time limit，nodes reject corresponding msg if timeout
	timeoutDuration time.Duration

	retryUnconnectedNodeDuration time.Duration

	mustOnChain           bool
	randHashStorageClient *eth.RandHashStorage // client for saving {reqID, randHashes} on chain

	randomNumbersPool *common.RandomNumbersPool

	EvilChecker EvilCheckerInterface

	metaDB _db.MetaDB
}

func NewRNG(ctx context.Context, conf *_conf.Config) (*RNG, error) {
	db, err := _db.NewDB(conf.MetaDBPath)
	if err != nil {
		return nil, err
	}
	var pendingSeq atomic.Int64
	pendingSeqInt64, err := db.GetPendingSeq()
	if err != nil {
		log.Error("get pendingSeq from level db failed, use default pendingSeq 0", "err", err)
	} else {
		log.Debug("get pendingSeq from level db", "pendingSeq", pendingSeqInt64)
		pendingSeq.Store(pendingSeqInt64)
	}
	return &RNG{
		UnimplementedRandomNumberGeneratorServer: _proto.UnimplementedRandomNumberGeneratorServer{},
		ctx:                                      ctx,
		wg:                                       &sync.WaitGroup{},

		//config: conf,

		nodeID: conf.Local,
		leader: conf.Leader,
		//Generators:                               conf.Nodes,
		role: conf.Role,

		//nodeInfoFromConfig:   conf.NodeInfoFromConfig,
		//unConnectedNodeTable: conf.UnConnectedNodeTable,

		nodeURLTable:       conf.NodeURLTable,
		nodeKeyTable:       conf.NodeKeyTable,
		connectionMap:      conf.ConnectionMap,
		nodeRpcClientTable: conf.NodeRpcClientTable,

		privKey:      conf.PrivKey,
		blsThreshold: conf.BlsThreshold,

		generatorRecords: NewGeneratorRandomNumberRecord(),
		leaderRecords:    NewLeaderRandomNumberRecord(),

		rpcCh:          make(chan interface{}, conf.RpcChanLen),
		errCh:          make(chan types.ErrMessage, 10),
		grpcServerPort: conf.GrpcServerPort,
		pendingSeq:     pendingSeq,
		//leaderReqs:     types.LeaderReqs{PendingSeqNum: pendingSeq},

		lastDiscardSeq:  conf.LastDiscardSeq,
		discardDuration: conf.DiscardDuration,
		timeoutDuration: conf.TimeoutDuration,

		//retryUnconnectedNodeDuration: conf.RetryUnconnectedNodeDuration,

		mustOnChain:           conf.MustOnChain,
		randHashStorageClient: conf.RandStorageCli,

		randomNumbersPool: common.NewRandomNumbersPool(conf.RandomNumbersPoolUpperLimit, pendingSeqInt64-1),

		EvilChecker: NewEvilChecker(),
		metaDB:      db,
	}, nil
}

/**
grpc service imply
*/

// RequestRandomNumber get new request from client
// return reqID
func (n *RNG) RequestRandomNumber(ctx context.Context, empty *emptypb.Empty) (*_proto.UnaryResponse, error) {
	log.Debug("receive  RequestRandomNumber")
	// 1. check self role is leader
	if n.nodeID != n.leader {
		log.Warn("local is generator, not leader. only leader handle RequestRandomNumber")
		// should send RequestMessage to leader, but local is not leader
		return &_proto.UnaryResponse{
			From:    n.nodeID,
			Message: types.Fail, // fail: 400
		}, errors.New("should send RequestMessage to leader, but local is not leader")
	}

	// 2. new reqID and get pending seq
	reqMsg := &_proto.RequestMessage{}
	reqMsg.ReqID = uuid.New().String()
	reqMsg.SeqNum = n.pendingSeq.Load()
	newPendingSeq := reqMsg.SeqNum + 1
	n.pendingSeq.Store(newPendingSeq)
	if err := n.metaDB.SetPendingSeq(newPendingSeq); err != nil {
		log.Error("save pendingSeq into leveldb failed", "err", err, "pendingSeq", newPendingSeq)
	}
	log.Debug("leader get reqID and seNum for new request", "seqNum", reqMsg.GetSeqNum(), "reqID", reqMsg.GetReqID(), "pendingSeq", newPendingSeq)

	// 3. set timeout
	timeout := time.Now().Add(n.timeoutDuration)
	timeoutByte := common.TimeToBytes(timeout)

	reqMsg.Timeout = timeoutByte

	// 4. push msg into chan
	n.rpcCh <- reqMsg
	log.Debug("pushed requestMessage into rpcCh successfully", "seq", reqMsg.GetSeqNum())
	return &_proto.UnaryResponse{
		From:    n.nodeID,
		ReqID:   reqMsg.ReqID,
		SeqNum:  reqMsg.SeqNum,
		Message: types.Success,
	}, nil
}

// SendInternalMessage process proto message from internal remote nodes
func (n *RNG) SendInternalMessage(ctx context.Context, message *_proto.TransportMessage) (*_proto.UnaryResponse, error) {
	log.Debug("receive  SendInternalMessage")
	msg, msgType, err := _grpc.DecodeTransportMessage(message)
	if err != nil {
		log.Error("decode transport message", "err", err)
		return &_proto.UnaryResponse{
			From:    n.nodeID,
			Message: types.Fail,
		}, err
	}

	log.Debug("finish decode transport message, pushing into rpcCh", "msgType", msgType)
	// todo: push buffer list if chan is full
	// todo: start a new loop to get msg from buffer and push into chan
	n.rpcCh <- msg
	log.Debug("pushed Internal message into rpcCh successfully")
	return &_proto.UnaryResponse{
		From:    n.nodeID,
		Message: types.Success,
	}, nil
}

// Start prepare the grpc server, and listen grpc request.
// Run main loop as leader or generator
func (n *RNG) Start() error {
	// grpc server
	log.Debug("try to listen", "addr", n.grpcServerPort, "network", "tcp")
	lis, err := net.Listen("tcp", n.grpcServerPort)
	if err != nil {
		log.Error("failed to listen", "err", err, "address", n.grpcServerPort)
		return err
	}

	s := grpc.NewServer()
	_proto.RegisterRandomNumberGeneratorServer(s, n)
	reflection.Register(s) // support the reflection API for grpcurl
	log.Debug("finished registering random number generator server, finished reflection.Register")
	go s.Serve(lis)

	n.wg.Add(1)
	n.run()
	n.wg.Wait()

	return nil
}

// run the main thread that handles leadership and RPC requests.
func (n *RNG) run() {
	log.Debug("running role")
	defer n.wg.Done()
	switch n.role.GetRole() {
	case common.LEADER:
		log.Info("run node as leader")
		n.wg.Add(2)
		go n.runLeader()
		go n.leaderDiscard()
	case common.GENERATOR:
		log.Info("run node as generator")
		n.wg.Add(2)
		go n.runGenerator()
		go n.generatorDiscard()
	default:
		log.Error("error node role", "role", n.role.GetRole())
		return
	}

	go n.listenErrMessage()
}

// runGenerator is run as goroutine
func (n *RNG) runGenerator() {
	log.Info("start run generator")
	defer n.wg.Done()
	if n.leader == n.nodeID {
		log.Error("config local role is generator, but config's leader is also local node")
		return
	}

	var isLocked bool
	for n.role.GetRole() == common.GENERATOR {
		if isLocked {
			n.generatorRecords.Unlock()
			isLocked = false
		}

		select {
		case <-n.ctx.Done():
			for _, conn := range n.connectionMap {
				conn.Close()
			}
			n.metaDB.CloseDB()
			return
		case e := <-n.rpcCh:
			if e == nil {
				log.Error("get nil from rpcCh")
				break
			}
			switch e.(type) {
			case *_proto.StartNewRoundMessage:
				msg := e.(*_proto.StartNewRoundMessage)
				metaData := msg.GetMetaData()
				seqNum, reqID, from, leader := metaData.GetSeqNum(), metaData.GetReqID(), msg.GetFrom(), metaData.GetLeader()
				log.Debug("generator received StartNewRoundMessage", "seqNum", seqNum, "reqID", reqID)

				// 0. check remote evil
				if n.EvilChecker.IsEvil(from) {
					log.Warn("get message from evil remote", "remote", from)
					break
				}

				// 1. check timeout
				timeByte := metaData.GetTimeout()
				timeout, err := common.BytesToTime(timeByte)
				if err != nil {
					log.Error("generator check StartNewRoundMessage time limit", "err", err, "seqNum", seqNum, "reqID", reqID)
					break
				}

				if time.Now().After(timeout) {
					log.Error("checked StartNewRoundMessage, msg timeout", "timeout", timeout, "seqNum", seqNum, "reqID", reqID)
					break
				}

				// 2. check from
				if from != leader {
					log.Error("get StartNewRoundMessage from no-leader", "seqNum", seqNum, "reqID", reqID)
					break
				}

				if n.leader != leader {
					log.Error("get msg from no-leader", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}

				n.generatorRecords.Lock()
				isLocked = true

				// 3. check status
				status := n.generatorRecords.seqStatusMap[seqNum]
				if status&types.Generated != 0 {
					log.Warn("have recieved same seqNum msg from other node", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}
				if status&types.Done != 0 {
					log.Warn("this seq has already processed", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}

				leaderSigData := &_proto.MetaData{
					Leader:  n.leader,
					SeqNum:  seqNum,
					ReqID:   reqID,
					Timeout: metaData.GetTimeout(),
				}

				if !n.verifyLeaderSignature(leaderSigData, metaData.GetSig()) {
					n.EvilChecker.AddEvilRecord(from)
					log.Error("verifyLeaderSignature failed, add evil record", "result", false, "from", from, "leader", leader, "seqNum", seqNum, "reqID", reqID)
					break
				}

				// 5. check status is Idle, init msg buffer
				if status&types.Idle == 0 {
					n.generatorRecords.msgBuffer[seqNum] = &GeneratorMsgBuffer{
						reqID:             reqID,
						collectRandomNums: make(map[string]*collectionElem),
					}
					n.generatorRecords.seqStatusMap[seqNum] = types.Idle
					n.generatorRecords.seqTimeoutQueue.Enqueue(common.TimeoutInfo{
						Seq:     seqNum,
						Timeout: timeout,
					})
					//n.generatorRecords.seqTimeoutMap[seqNum] = timeout
				}

				// 6. generate self random num, build internal message and route it
				routeMsg, collElem, err := n.generateRandNumberAndBuildSingleRandMsg(metaData)
				if err != nil {
					log.Error("generateRandNumberAndBuildSingleRandMsg", "err", err, "from", from, "leader", leader, "seqNum", seqNum, "reqID", reqID)
					break
				}
				go n.routeInternalMessage(routeMsg, types.SingleRandType)

				n.generatorRecords.seqStatusMap[seqNum] |= types.Generated
				log.Debug("generate random num and routing to remotes", "from", from, "seqNum", seqNum, "reqID", reqID, "myRand", collElem.randNum)

				// 7. collect myRand
				n.generatorRecords.msgBuffer[seqNum].singleRandomNum = collElem.randNum
				n.generatorRecords.msgBuffer[seqNum].collectRandomNums[n.nodeID] = collElem
				n.generatorRecords.seqStatusMap[seqNum] |= types.Collecting
				log.Debug("collected self generated random number", "from", from, "seqNum", seqNum, "reqID", reqID)

				n.generatorRecords.Unlock()
				isLocked = false
				log.Debug("finish switch case", "type", "StartNewRoundMessage", "seqNum", seqNum)

			case *_proto.SingleRandMessage:
				msg := e.(*_proto.SingleRandMessage)
				metaData := msg.GetMetaData()
				from, leader, seqNum, reqID := msg.GetFrom(), metaData.GetLeader(), metaData.GetSeqNum(), metaData.GetReqID()

				// 0. check remote evil
				if n.EvilChecker.IsEvil(from) {
					log.Warn("get message from evil remote", "remote", from)
					break
				}

				// 1. check timeout
				timeByte := metaData.GetTimeout()
				timeout, err := common.BytesToTime(timeByte)
				if err != nil {
					log.Error("generator check SingleRandMessage time limit", "err", err, "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}
				if time.Now().After(timeout) {
					log.Error("generator checked SingleRandMessage, msg timeout", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}

				// 2. check from and leader, from should not in the msg buffer
				if _, ok := n.nodeURLTable[from]; !ok {
					log.Error("rpc msg from unknown", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}
				if leader != n.leader {
					log.Error("leader not match", "expected", n.leader, "got", leader)
					break
				}

				n.generatorRecords.Lock()
				isLocked = true

				status := n.generatorRecords.seqStatusMap[seqNum]
				if status&types.Idle != 0 { // msgBuffer[seqNum] is not nil
					if _, ok := n.generatorRecords.msgBuffer[seqNum].collectRandomNums[from]; ok {
						log.Warn("got singleRandMessage from same remote duplicated", "seq", seqNum, "from", from)
						break
					}
				}

				// 3. check inside signatures
				leaderSigData := &_proto.MetaData{
					Leader:  n.leader,
					SeqNum:  seqNum,
					ReqID:   reqID,
					Timeout: metaData.GetTimeout(),
				}
				if !n.verifyLeaderSignature(leaderSigData, metaData.GetSig()) {
					n.EvilChecker.AddEvilRecord(from)
					log.Error("verifyLeaderSignature", "result", false, "from", from, "leader", leader, "seqNum", seqNum, "reqID", reqID)
					break
				}

				// 4. check status
				// init msgBuffer
				// generate self random number and route it
				if status&types.Idle == 0 { // init
					n.generatorRecords.msgBuffer[seqNum] = &GeneratorMsgBuffer{
						reqID:             reqID,
						singleRandomNum:   0,
						collectRandomNums: make(map[string]*collectionElem),
					}
					n.generatorRecords.seqStatusMap[seqNum] = types.Idle
					n.generatorRecords.seqTimeoutQueue.Enqueue(common.TimeoutInfo{
						Seq:     seqNum,
						Timeout: timeout,
					})
					//n.generatorRecords.seqTimeoutMap[seqNum] = timeout
				}
				if status&types.Generated == 0 {
					routeMsg, collElem, err := n.generateRandNumberAndBuildSingleRandMsg(metaData)
					if err != nil {
						log.Error("generateRandNumberAndBuildSingleRandMsg", "err", err, "from", from, "seqNum", seqNum, "reqID", reqID)
						break
					}
					go n.routeInternalMessage(routeMsg, types.SingleRandType)

					n.generatorRecords.seqStatusMap[seqNum] |= types.Generated
					log.Debug("generate random num and routing to remotes", "seqNum", seqNum, "reqID", reqID, "myRand", collElem.randNum)

					n.generatorRecords.msgBuffer[seqNum].singleRandomNum = collElem.randNum
					n.generatorRecords.msgBuffer[seqNum].collectRandomNums[n.nodeID] = collElem
					n.generatorRecords.seqStatusMap[seqNum] |= types.Collecting
					log.Debug("collected self generated random number", "seqNum", seqNum)
				}

				// 5. collect remote random number
				remoteRand := msg.GetRandomNum()
				remoteSig, err := blsSignatures.SignatureFromBytes(msg.GetSig()) // not check sig now
				if err != nil {
					log.Error("get signature from singleRandMessage.Sig", "err", err, "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}
				// get signed message
				plaintextMsg := &_proto.SingleRandMessage{
					MetaData:  metaData,
					From:      msg.GetFrom(),
					RandomNum: msg.GetRandomNum(),
				}
				routeMsgByte, err := proto.Marshal(plaintextMsg)
				if err != nil {
					log.Error("marshal signedMsg plaintext failed", "err", err)
					break
				}
				n.generatorRecords.msgBuffer[seqNum].collectRandomNums[from] = &collectionElem{
					randNum:   remoteRand,
					blsSig:    remoteSig,
					signedMsg: routeMsgByte,
				}
				n.generatorRecords.seqStatusMap[seqNum] |= types.Collecting
				log.Debug("collect random number from remote", "remote", from, "remoteRandNum", remoteRand, "seqNum", seqNum, "reqID", reqID)

				// 6. compare number of random numbers to the threshold
				// aggregate signatures and verify
				// build collectionMessage and route it
				collNumbsLen := len(n.generatorRecords.msgBuffer[seqNum].collectRandomNums)
				n.generatorRecords.Unlock()
				isLocked = false
				if collNumbsLen >= n.blsThreshold {
					log.Info("collect enough random numbers, trying to aggregate and route", "seqNum", seqNum, "reqID", reqID)
					go n.handleCollectionRandNums(seqNum, metaData)
				}
				log.Debug("finish switch case", "type", "SingleRandMessage", "seqNum", seqNum)

			default:
				log.Warn("ignore msg, GENERATOR only process StartNewRoundMessage and SingleRandMessage")
			}
		}
	}
}

// runLeader
func (n *RNG) runLeader() {
	log.Info("start runLeader")
	defer n.wg.Done()
	if n.leader != n.nodeID {
		log.Error("config leader is not local node, but local role is leader")
		return
	}

	var isLocked bool
	// handle remote msg
	for { //  n.role.GetRole() == common.LEADER {
		if isLocked {
			n.leaderRecords.Unlock()
			isLocked = false
		}

		select {
		case <-n.ctx.Done():
			for _, conn := range n.connectionMap {
				conn.Close()
			}
			n.metaDB.CloseDB()
			return
		case e := <-n.rpcCh:
			if e == nil {
				log.Error("get nil from rpcCh")
				break
			}
			switch e.(type) {
			case *_proto.RequestMessage: // client -> leader
				msg := e.(*_proto.RequestMessage)
				seqNum, reqID := msg.GetSeqNum(), msg.GetReqID()

				log.Debug("leader received RequestMessage", "seqNum", seqNum, "reqID", reqID)
				// 1. check timeout
				timeout, err := common.BytesToTime(msg.GetTimeout())
				if err != nil {
					log.Error("leader check RequestMessage time limit", "err", err, "seqNum", seqNum, "reqID", reqID)
					break
				}
				if time.Now().After(timeout) {
					log.Error("leader checked RequestMessage, msg timeout", "seqNum", seqNum, "reqID", reqID)
					break
				}
				log.Debug("leader checked RequestMessage time limit ok", "seqNum", seqNum, "reqID", reqID)

				// 2. sign leader signature as inside signature int the metaData
				metaData := &_proto.MetaData{
					Leader:  n.nodeID,
					SeqNum:  seqNum,
					ReqID:   reqID,
					Timeout: msg.GetTimeout(),
				}

				metaDataByte, err := proto.Marshal(metaData)
				if err != nil {
					log.Error("leader handle client request message, marshal metaData", "err", err, "seqNum", seqNum)
					break
				}
				leaderSig, err := blsSignatures.SignMessage(n.privKey, metaDataByte)
				if err != nil {
					log.Error("leader sign metadata", "err", err, "seqNum", seqNum)
					break
				}
				metaData.Sig = blsSignatures.SignatureToBytes(leaderSig)
				log.Debug("leader finished signing metaData", "seqNum", seqNum)

				// 3. rout
				startNewRoundMsg := &_proto.StartNewRoundMessage{
					MetaData: metaData,
					From:     n.nodeID,
				}
				log.Debug("leader routing internal message", "seqNum", msg.GetSeqNum(), "type", types.StartType)
				go n.routeInternalMessage(startNewRoundMsg, types.StartType)

				// update record
				seq := metaData.GetSeqNum()
				n.leaderRecords.Lock()
				n.leaderRecords.seqStatusMap[seq] |= types.Idle // leader的Idle表示已经开始了，“processing”
				//n.leaderRecords.seqTimeoutMap[seq] = timeout
				n.leaderRecords.seqTimeoutQueue.Enqueue(common.TimeoutInfo{
					Seq:     seq,
					Timeout: timeout,
				})
				n.leaderRecords.msgBuffer[seq] = NewCollectionBuffer()
				n.leaderRecords.Unlock()
				log.Debug("finish switch case", "type", "RequestMessage", "seqNum", seqNum)
			case *_proto.CollectionRandMessage: // leader receive msg from generator
				msg := e.(*_proto.CollectionRandMessage)
				metaData := msg.GetMetaData()
				from, leader, timeoutByte, seqNum, reqID := msg.GetFrom(), metaData.GetLeader(), metaData.Timeout, metaData.GetSeqNum(), metaData.GetReqID()
				log.Debug("leader received CollectionRandMessage", "from", from, "seqNum", seqNum, "reqID", reqID)

				// 0. check remote evil
				if n.EvilChecker.IsEvil(from) {
					log.Warn("get message from evil remote", "remote", from)
					break
				}

				n.leaderRecords.Lock()
				isLocked = true
				// 1. check status and timeout
				if n.leaderRecords.seqStatusMap[seqNum]&types.Idle == 0 {
					log.Error("got CollectionRandMessage but has no status", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}
				timeout, err := common.BytesToTime(timeoutByte)
				if err != nil {
					log.Error("leader check CollectionRandMessage time limit", "err", err)
					break
				}
				if time.Now().After(timeout) {
					log.Error("leader checked CollectionRandMessage, msg timeout", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}

				// 2. check leader and from, prevent duplicate msg
				if leader != n.leader {
					log.Error("metaData.leader not match", "got", leader, "want", n.leader, "seqNum", seqNum, "reqID", reqID)
					break
				}
				if _, ok := n.nodeURLTable[from]; !ok {
					log.Error("msg from unknown node", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}
				if _, ok := n.leaderRecords.msgBuffer[seqNum].remote2Collection[from]; ok {
					log.Warn("received duplicate collectionRandMessage from same node", "from", from, "seqNum", seqNum, "reqID", reqID)
					break
				}

				// 3. verify inside signature
				leaderSigData := &_proto.MetaData{
					Leader:  n.leader,
					SeqNum:  seqNum,
					ReqID:   reqID,
					Timeout: metaData.GetTimeout(),
				}
				if !n.verifyLeaderSignature(leaderSigData, metaData.GetSig()) {
					n.EvilChecker.AddEvilRecord(from)
					log.Error("verifyLeaderSignature", "result", false, "from", from, "leader", leader, "seqNum", seqNum, "reqID", reqID)
					break
				}

				// 4. keep outside signature for aggregating
				sig, err := blsSignatures.SignatureFromBytes(msg.GetSig())
				if err != nil {
					log.Error("decode signature", "err", err, "seqNum", seqNum, "from", from, "sigByteHex", hexutil.Encode(msg.GetSig()))
					break
				}
				if n.leaderRecords.msgBuffer[seqNum].plainAndSig == nil {
					n.leaderRecords.msgBuffer[seqNum].plainAndSig = make(map[string]struct {
						plain []byte
						sig   blsSignatures.Signature
					})
				}
				n.leaderRecords.msgBuffer[seqNum].plainAndSig[from] = struct {
					plain []byte
					sig   blsSignatures.Signature
				}{plain: msg.GetSigPlaintext(), sig: sig}

				n.leaderRecords.Unlock()
				isLocked = false

				// 5. append collection records
				// 一个from只能发一组collection
				// 不同组collection中，来自同一个node的随机数应该相同，如果不同，取第一个？还是取数量最多的？
				// 取第一不行，恶意节点可以以最快速度虚构来自其他节点的随机数，第一时间发出，抢占有效权
				// 应该取最多
				// 或者以randNum为map的key？randNum:struct{}
				// 要防止一个node分别给不同其他node发不同的randNum
				go func() {
					log.Debug("go handle collection numbers", "seqNum", seqNum, "reqID", reqID)
					// 1. keep record
					collection := msg.CollectionRandomNumbers // node:rand

					n.leaderRecords.Lock()
					if n.leaderRecords.msgBuffer[seqNum] == nil {
						n.leaderRecords.msgBuffer[seqNum] = &CollectionBuffer{
							remote2Collection: make(map[string][]uint64),
							nodeRandTimes:     make(map[string]RandTimes),
						}
					}

					// append remote2Collection, nodeRandTimes
					remoteCollection := make([]uint64, 0)
					for node, rand := range collection {
						remoteCollection = append(remoteCollection, rand)
						if n.leaderRecords.msgBuffer[seqNum].nodeRandTimes[node] == nil {
							n.leaderRecords.msgBuffer[seqNum].nodeRandTimes[node] = make(RandTimes) // randNum:times
						}
						// increase leader receive node1's random number's times according to different collection
						// if evil exist, maybe receive different random number from different collection, but the generator is the same
						n.leaderRecords.msgBuffer[seqNum].nodeRandTimes[node][rand]++
					}
					n.leaderRecords.msgBuffer[seqNum].remote2Collection[from] = remoteCollection
					collectionsNum := len(n.leaderRecords.msgBuffer[seqNum].remote2Collection)
					n.leaderRecords.Unlock()

					// if upto threshold, go to get mergedRandNums from map
					if collectionsNum >= n.blsThreshold {
						n.leaderRecords.Lock()
						n.leaderRecords.seqStatusMap[seqNum] |= types.Mergeing
						n.leaderRecords.Unlock()
						go func() {
							// 1. aggregate sig and verify
							// todo； with pool
							sigs := make([]blsSignatures.Signature, 0)
							plains := make([][]byte, 0)
							pubs := make([]blsSignatures.PublicKey, 0)
							nodes := make([]string, 0)

							n.leaderRecords.RLock()
							for node, plainAndSig := range n.leaderRecords.msgBuffer[seqNum].plainAndSig {
								sigs = append(sigs, plainAndSig.sig)
								plains = append(plains, plainAndSig.plain)
								pubs = append(pubs, n.nodeKeyTable[node])
								nodes = append(nodes, node)
							}
							n.leaderRecords.RUnlock()

							aggrSignature := blsSignatures.AggregateSignatures(sigs)
							valid, err := blsSignatures.VerifyAggregatedSignatureDifferentMessages(aggrSignature, plains, pubs)

							if !valid {
								log.Error("collectionRandMessages aggregatedSignature invalid", "seqNum", seqNum, "reqID", reqID)
								if err != nil {
									log.Error("leader verify collectionRandMessages aggregatedSignature failed", "err", err, "seqNum", seqNum, "reqID", reqID)

								}

								// check evil situation one by one
								evils := n.EvilChecker.CheckOneByOne(sigs, plains, pubs, nodes)
								if len(evils)+n.blsThreshold > len(n.nodeRpcClientTable) {
									log.Warn("impossible to receive enough signature", "seqNum", seqNum, "reqID", reqID, "evils number", len(evils))
								} else {
									log.Debug("still has chance to get enough validate signature, del the evil record")
									for _, evilNode := range evils {
										n.leaderRecords.Lock()
										delete(n.leaderRecords.msgBuffer[seqNum].nodeRandTimes, evilNode)
										delete(n.leaderRecords.msgBuffer[seqNum].plainAndSig, evilNode)
										delete(n.leaderRecords.msgBuffer[seqNum].remote2Collection, evilNode)
										n.leaderRecords.Unlock()
									}
								}

								return
							}
							log.Debug("leader verified collections aggregated signature", "seqNum", seqNum, "reqID", reqID)

							// 2. merge
							// leader store {seq, map[random]from} on chain by contract
							// leader store merged([]random) in memory for client
							// todo 要不要增加各节点在不同阶段结果上链，txhash作为参数参加到流程的数据传递中，每一步的txhash都作为下一步的data，可追溯
							merged := make([]uint64, 0)

							n.leaderRecords.RLock()
							for node, randTimes := range n.leaderRecords.msgBuffer[seqNum].nodeRandTimes {
								if randTimes == nil || len(randTimes) == 0 {
									log.Warn("for leaderRecords.msgBuffer[seqNum].nodeRandTimes, get empty randTimes", "seqNum", seqNum, "node", node)
									continue
								}

								// len(randTimes) > 1: one node sent different randNums to different collectors,
								// choose the random number with largest times
								// len(randTimes) == 1, expected status
								randomNum, largeTimes := uint64(0), uint64(0)
								for rand, times := range randTimes { // 这个节点产生过的所有随机数遍历，选择被记录次数最多的那个
									if times > largeTimes {
										randomNum, largeTimes = rand, times
									}
								}
								if largeTimes != 0 {
									merged = append(merged, randomNum)
								}
							}
							n.leaderRecords.RUnlock()

							n.leaderRecords.Lock()
							// 3. keep in memory
							n.leaderRecords.mergedRandNums[seqNum] = merged
							n.leaderRecords.seqStatusMap[seqNum] |= types.Done
							log.Debug("leader merged random numbers", "randNumList", merged, "seqNum", seqNum, "reqID", reqID)

							// 4. save mergedRandNums[seqNum] on chain or da
							// todo impl save with Da
							// save random numbers hash on chain
							if n.randHashStorageClient != nil {
								txHash, err := n.randHashStorageClient.LeaderSaveResultsOnChain(reqID, merged)
								if err != nil {
									log.Error("leader save random numbers hashes on chain failed", "err", err, "seqNum", seqNum, "reqID", reqID)
									if n.mustOnChain {
										n.leaderRecords.Unlock()
										return
									}
								}
								log.Debug("leader save random numbers hash on chain successfully", "txHash", txHash, "seqNum", seqNum, "reqID", reqID)

							}
							n.leaderRecords.seqStatusMap[seqNum] |= types.Safe
							n.leaderRecords.Unlock()

							// 5. results append pool
							log.Debug("putting random numbers into pool", "seqNum", seqNum, "reqID", reqID)
							go func() {
								n.randomNumbersPool.PutRNGResults(reqID, merged)
								n.randomNumbersPool.SetLatestSeqInPool(seqNum)
							}()
						}()

					}
				}()
				log.Debug("finish switch case", "type", "CollectionRandMessage", "seqNum", seqNum)

			default:
				log.Warn("ignore msg, LEADER only process RequestMessage and MergedRandMessage")
			}
		}
	}
}

//func (n *RNG) manageConnection(ctx context.Context) {
//	ticker := time.NewTicker(n.retryUnconnectedNodeDuration)
//	var (
//		conn *grpc.ClientConn
//		err  error
//	)
//
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case <-ticker.C:
//			//for node, conn := range n.connectionMap {
//			//	newState := conn.GetState()
//			//	log.Debug("Connection state changed", "node", node, "newState", newState)
//			//
//			//	if newState >= connectivity.TransientFailure {
//			//		log.Debug("Connection is transiently unavailable, reconnecting...")
//			//		// 连接断开，尝试重新建立连接
//			//		for {
//			//			err := conn.ResetChannelz()
//			//			if err == nil {
//			//				log.Println("Connection reestablished successfully.")
//			//				break
//			//			}
//			//			log.Printf("Failed to reestablish connection: %v", err)
//			//		}
//			//	}
//			//}
//
//			// retry unconnected grpc connection
//			if len(n.unConnectedNodeTable) == 0 {
//				// 第一轮连上，后续宕机断连，再次run起来，外界进来的连接会在被调用方法时主动连接上，向外的连接会主动创建，
//				return
//			}
//			log.Debug("重连 unConnectedNodeTable", "unConnectedNodeTable", n.unConnectedNodeTable)
//			for node, _ := range n.unConnectedNodeTable {
//				log.Debug("Attempting to connect to remote node", "nodeID", node, "attempt")
//				conn, err = grpc.Dial(
//					n.nodeInfoFromConfig[node].Url,
//					grpc.WithTransportCredentials(insecure.NewCredentials()),
//					//grpc.WithBlock(),
//				)
//				if err == nil {
//
//					delete(n.unConnectedNodeTable, node)
//					n.connectionMap[node] = conn
//					n.nodeRpcClientTable[node] = _proto.NewRandomNumberGeneratorClient(conn)
//					n.nodeURLTable[node] = n.nodeInfoFromConfig[node].Url
//				}
//			}
//		}
//
//	}
//}

func (n *RNG) listenErrMessage() {
	for {
		select {
		case <-n.ctx.Done():
		// handle the left err in errCh

		case e := <-n.errCh:
			// store into da and chain
			log.Debug("handle errMsg", "type", e.ErrType, "seqNum", e.SeqNum, "reqID", e.ReqID)
		}
	}
}

func (n *RNG) routeInternalMessage(msg interface{}, msgType string) {
	var (
		msgByte []byte
		err     error
	)
	// wrap msg
	switch msgType {
	case types.StartType:
		log.Debug("routing StartType msg")
		msgByte, err = proto.Marshal(msg.(*_proto.StartNewRoundMessage))
		if err != nil {
			log.Error("routeInternalMessage marshal StartNewRoundMessage", "err", err)
			return
		}
	case types.SingleRandType:
		log.Debug("routing SingleRandType msg")
		msgByte, err = proto.Marshal(msg.(*_proto.SingleRandMessage))
		if err != nil {
			log.Error("routeInternalMessage marshal SingleRandMessage", "err", err)
			return
		}
	case types.CollectionRandType:
		log.Debug("routing CollectionRandType msg")
		msgByte, err = proto.Marshal(msg.(*_proto.CollectionRandMessage))
		if err != nil {
			log.Error("routeInternalMessage marshal CollectionRandMessage", "err", err)
			return
		}
	default:
		log.Error("routeInternalMessage unknown type", "type", msgType)
		return
	}

	transMsg := &_proto.TransportMessage{
		Version: "",
		MsgType: msgType,
		Data:    msgByte,
	}

	for remote, rpcCli := range n.nodeRpcClientTable {
		if remote == n.nodeID {
			continue
		}
		_, err = rpcCli.SendInternalMessage(n.ctx, transMsg)
		if err != nil {
			log.Error("routeInternalMessage.SendInternalMessage failed", "err", err, "from", n.nodeID, "to", remote)
		}
	}
}

func (n *RNG) generateRandNumberAndBuildSingleRandMsg(metaData *_proto.MetaData) (*_proto.SingleRandMessage, *collectionElem, error) {
	// 1. generate self random number
	myRand, err := random.GenerateRandomNum()
	if err != nil {
		log.Error("generate self random number", "err", err)
		return nil, nil, err
	}

	// 2. build route msg and sign
	routeMsg := &_proto.SingleRandMessage{
		MetaData:  metaData,
		From:      n.nodeID,
		RandomNum: myRand,
	}

	routeMsgByte, err := proto.Marshal(routeMsg)
	if err != nil {
		log.Error("marshal routeMsg", "err", err)
		return nil, nil, err
	}
	sig, err := blsSignatures.SignMessage(n.privKey, routeMsgByte)
	if err != nil {
		log.Error("bls sign routeMsgByte", "err", err)
		return nil, nil, err
	}

	routeMsg.Sig = blsSignatures.SignatureToBytes(sig)

	return routeMsg, &collectionElem{
		myRand,
		sig,
		routeMsgByte,
	}, nil
}

// handleCollectionRandNums call by generator
func (n *RNG) handleCollectionRandNums(seqNum int64, metaData *_proto.MetaData) {
	reqID := metaData.GetReqID()

	// 1. aggregate pubkey and signedReq_1
	sigList := make([]blsSignatures.Signature, 0)
	signedMsgList := make([][]byte, 0)
	pubList := make([]blsSignatures.PublicKey, 0)
	randNums := make([]uint64, 0)
	randNumsWithFrom := make(map[string]uint64) // remote:rand
	nodes := make([]string, 0)

	n.generatorRecords.RLock()
	for remote, elem := range n.generatorRecords.msgBuffer[seqNum].collectRandomNums {
		sigList = append(sigList, elem.blsSig)
		signedMsgList = append(signedMsgList, elem.signedMsg)
		pubList = append(pubList, n.nodeKeyTable[remote])
		randNums = append(randNums, elem.randNum)
		randNumsWithFrom[remote] = elem.randNum
		nodes = append(nodes, remote)
	}
	n.generatorRecords.RUnlock()

	aggrSignature := blsSignatures.AggregateSignatures(sigList)
	log.Debug("handleCollectionRandNums and aggregate signatures", "aggregated signature", hexutil.Encode(blsSignatures.SignatureToBytes(aggrSignature)))

	// verify aggregator signature
	log.Debug("VerifyAggregatedSignatureDifferentMessages", "len(signedMsgList)", len(signedMsgList), "len(pubList)", len(pubList))
	ok, err := blsSignatures.VerifyAggregatedSignatureDifferentMessages(aggrSignature, signedMsgList, pubList)

	if !ok {
		log.Error("verify aggregator signature not passed")
		if err != nil {
			log.Error("VerifyAggregatedSignatureDifferentMessages", "err", err, "seqNum", seqNum, "reqID", reqID)
		}

		evils := n.EvilChecker.CheckOneByOne(sigList, signedMsgList, pubList, nodes)
		if len(evils)+n.blsThreshold > len(n.nodeRpcClientTable) {
			log.Warn("impossible to receive enough signature", "seqNum", seqNum, "reqID", reqID, "evils number", len(evils))
		} else {
			log.Debug("still has chance to get enough validate signature, del the evil record")
			for _, evilNode := range evils {
				n.generatorRecords.Lock()
				delete(n.generatorRecords.msgBuffer[seqNum].collectRandomNums, evilNode)
				n.generatorRecords.Unlock()
			}
		}
		return
	}
	log.Debug("handleCollectionRandNums verify aggregator signature successfully")

	// build collection msg and route it
	routeMsg := &_proto.CollectionRandMessage{
		MetaData:                metaData,
		From:                    n.nodeID,
		CollectionRandomNumbers: randNumsWithFrom,
	}

	routeMsgByte, err := proto.Marshal(routeMsg)
	if err != nil {
		log.Error("marshal CollectionRandMessage failed", "err", err, "seqNum", seqNum, "reqID", reqID)
		return
	}

	sig, err := blsSignatures.SignMessage(n.privKey, routeMsgByte)
	if err != nil {
		log.Error("Sign CollectionRandMessage failed", "err", err, "seqNum", seqNum, "reqID", reqID)
		return
	}

	routeMsg.Sig = blsSignatures.SignatureToBytes(sig)
	routeMsg.SigPlaintext = routeMsgByte

	log.Debug("generator routing collectionRandType message", "seqNum", seqNum, "reqID", reqID)
	go n.routeInternalMessage(routeMsg, types.CollectionRandType)

	if n.randHashStorageClient != nil {
		log.Debug("generator saving collectionRandTypeMessage on chain", "seqNum", seqNum, "reqID", reqID)
		go n.randHashStorageClient.GeneratorSaveCollectionOnChain(reqID, randNums)
	}
}

func (n *RNG) verifyLeaderSignature(leaderSigData *_proto.MetaData, metaSig []byte) bool {
	//leaderSigDataByte, err := json.Marshal(leaderSigData)
	leaderSigDataByte, err := proto.Marshal(leaderSigData)
	if err != nil {
		log.Error("verifyLeaderSignature marshal leaderSign data", "err", err)
		return false
	}
	leaderSigByte, err := blsSignatures.SignatureFromBytes(metaSig)
	if err != nil {
		log.Error("verifyLeaderSignature SignatureFromBytes", "err", err)
		return false
	}
	ok, err := blsSignatures.VerifySignature(leaderSigByte, leaderSigDataByte, n.nodeKeyTable[n.leader])
	if err != nil {
		log.Error("verifyLeaderSignature check leader signature failed", "err", err)
		return false
	}
	if !ok {
		log.Error("verifyLeaderSignature check leader signature", "ok", ok)
		return false
	}
	return true
}

// JudgeAndStartNewGenerateRound check the number of random numbers which will be generated later
// if random numbers in processing is not enough, start a new generate round
func (n *RNG) JudgeAndStartNewGenerateRound(poolSize int) error {
	if poolSize >= n.GetPoolUpperLimit() {
		return nil
	}

	latestSeqInPool, processingLatestSeq, lastDiscardSeq, generatorNum := n.GetPoolLatestSeq(), n.GetLatestProcessingSeq(), n.GetLastDiscardSeq(), n.GetGeneratorsNum()

	// maybe some seq will not get results caused by connection error or generator nodes error, choose the bigger seq
	var biggerSeq int64
	if latestSeqInPool > lastDiscardSeq {
		biggerSeq = latestSeqInPool
	} else {
		biggerSeq = lastDiscardSeq
	}
	if float64((processingLatestSeq-biggerSeq)*int64(generatorNum)) < float64(n.GetPoolUpperLimit())*1.2 {
		log.Info("The number of random numbers expected to be generated does not exceed 120% of upperLimit, call the leader to start generating random numbers", "latestSeqInPool", latestSeqInPool, "processingLatestSeq", processingLatestSeq)
		resp, err := n.RequestRandomNumber(context.Background(), nil)
		if err != nil {
			log.Error("api RequestRandomNumber failed", "err", err)
			//c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return err
		}
		log.Info("api call RequestRandomNumber succ,", "reqID", resp.GetReqID())
		return _errors.ErrLeftRandoNumbersNotEnough
	}
	log.Debug("not start a new generating round", "latestSeqInPool", latestSeqInPool, "lastDiscardSeq", lastDiscardSeq, "processingLatestSeq", processingLatestSeq)
	// return nil: not start a new generate round
	return nil
}

func (n *RNG) GetRandNumFromPool(num int) ([]common.PoolElem, int) {
	return n.randomNumbersPool.DequeueMultiRandNums(num)
}

func (n *RNG) GetPoolUpperLimit() int {
	return n.randomNumbersPool.GetUpperLimit()
}

func (n *RNG) GetPoolLatestSeq() int64 {
	return n.randomNumbersPool.GetLatestSeqInPool()
}

func (n *RNG) GetLatestProcessingSeq() int64 {
	return n.pendingSeq.Load()
}

func (n *RNG) GetGeneratorsNum() int {
	return len(n.nodeURLTable)
}

func (n *RNG) GetLastDiscardSeq() int64 {
	return atomic.LoadInt64(&n.lastDiscardSeq)
}

func (n *RNG) GetSizeofPool() int {
	return n.randomNumbersPool.Size()
}

func (n *RNG) GetSingleRandomNumber() (common.PoolElem, error) {
	res, poolSize := n.GetRandNumFromPool(1)
	n.JudgeAndStartNewGenerateRound(poolSize)
	if res == nil || len(res) == 0 {
		return common.PoolElem{}, _errors.ErrLeftRandoNumbersNotEnough
	}
	return res[0], nil
}

func (n *RNG) GetMultiRandomNumbers(num int) ([]common.PoolElem, error) {
	res, pooSize := n.GetRandNumFromPool(num)
	n.JudgeAndStartNewGenerateRound(pooSize)
	if res == nil || len(res) == 0 {
		return nil, _errors.ErrLeftRandoNumbersNotEnough
	}

	return res, nil
}
