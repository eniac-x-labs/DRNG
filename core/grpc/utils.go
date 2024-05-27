package grpc

import (
	"errors"
	"fmt"

	"DRNG/common/types"
	_proto "DRNG/core/grpc/proto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"

	//"google.golang.org/protobuf/proto"
	//"github.com/gogo/protobuf/proto"
	"google.golang.org/protobuf/proto"
)

func DecodeTransportMessage(data *_proto.TransportMessage) (interface{}, string, error) { // return (unmarshaledMsg, type, error)
	var err error

	buf := data.Data
	switch data.MsgType {
	case types.StartType:
		log.Debug("decoding", "type", types.StartType, "detail", hexutil.Encode(buf))
		msg := &_proto.StartNewRoundMessage{}
		if err = proto.Unmarshal(buf, msg); err != nil {
			return nil, types.StartType, err
		}
		return msg, types.StartType, nil
	case types.SingleRandType:
		log.Debug("decoding", "type", types.SingleRandType)

		msg := &_proto.SingleRandMessage{}
		if err = proto.Unmarshal(buf, msg); err != nil {
			return nil, types.SingleRandType, err
		}
		return msg, types.SingleRandType, nil
	case types.CollectionRandType:
		log.Debug("decoding", "type", types.CollectionRandType)

		msg := &_proto.CollectionRandMessage{}
		if err = proto.Unmarshal(buf, msg); err != nil {
			return nil, types.CollectionRandType, err
		}
		return msg, types.CollectionRandType, nil
	case types.MergedRandType:
		log.Debug("decoding", "type", types.MergedRandType)

		msg := &_proto.MergedRandMessage{}
		if err = proto.Unmarshal(buf, msg); err != nil {
			return nil, types.MergedRandType, err
		}
		return msg, types.MergedRandType, nil
	default:
		return nil, types.UnknownType, errors.New(fmt.Sprintf("unknown msg type: %s", data.MsgType))
	}
}
