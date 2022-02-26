package ptypes

import (
	pb "github.com/si3nloong/sqlike/v2/protobuf"
	"github.com/si3nloong/sqlike/v2/types"
)

func Key(k *pb.Key) (*types.Key, error) {
	return protoToKey(k), nil
}

func KeyProto(k *types.Key) (*pb.Key, error) {
	return keyToProto(k), nil
}

func protoToKey(pk *pb.Key) *types.Key {
	if pk == nil {
		return nil
	}

	return &types.Key{
		Namespace: pk.Namespace,
		Kind:      pk.Kind,
		NameID:    pk.NameID,
		IntID:     pk.IntID,
		Parent:    protoToKey(pk.Parent),
	}
}

func keyToProto(k *types.Key) *pb.Key {
	if k == nil {
		return nil
	}

	return &pb.Key{
		Namespace: k.Namespace,
		Kind:      k.Kind,
		NameID:    k.NameID,
		IntID:     k.IntID,
		Parent:    keyToProto(k.Parent),
	}
}
