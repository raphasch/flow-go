package topology

import (
	"fmt"

	"github.com/onflow/flow-go/crypto/random"
	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/flow/filter"
	"github.com/onflow/flow-go/network/gossip/libp2p/channel"
	"github.com/onflow/flow-go/state/protocol"
)

// RandomizedTopology is a generates a random topology per channel.
// By random topology we mean a node is connected to any other co-channel nodes with some
// edge probability.
type RandomizedTopology struct {
	me      flow.Identifier             // used to keep identifier of the node
	state   protocol.ReadOnlyState      // used to keep a read only protocol state
	subMngr channel.SubscriptionManager // used to keep track topics the node subscribed to
	chance  uint64                      // used to translate connectedness probability into a number in [0, 100]
	rng     random.Rand                 // used as a stateful random number generator to sample edges
}

// NewRandomizedTopology returns an instance of the RandomizedTopology.
func NewRandomizedTopology(nodeID flow.Identifier, edgeProb float64, state protocol.ReadOnlyState, subMngr channel.SubscriptionManager) (*RandomizedTopology, error) {
	// edge probability should be a positive value between 0 and 1. However,
	// we like it to be strictly greater than zero. Also, at the current scale of
	// Flow, we need it to be greater than 0.01 to support probabilistic connectedness.
	if edgeProb < 0.01 || edgeProb > 1 {
		return nil, fmt.Errorf("randomized topology probability should in in range of [0.01, 1], wrong value: %f", edgeProb)
	}

	// generates seed and random number generator
	seed, err := byteSeedFromID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("could not generate seed from id:%w", err)
	}
	rng, err := random.NewRand(seed)
	if err != nil {
		return nil, fmt.Errorf("could not generate random number generator: %w", err)
	}

	return &RandomizedTopology{
		me:      nodeID,
		state:   state,
		subMngr: subMngr,
		chance:  uint64(100 * edgeProb),
		rng:     rng,
	}, nil

}

// GenerateFanout receives IdentityList of entire network and constructs the fanout IdentityList
// of this instance. A node directly communicates with its fanout IdentityList on epidemic dissemination
// of the messages (i.e., publish and multicast).
// Independent invocations of GenerateFanout on different nodes collaboratively must construct a cohesive
// connected graph of nodes that enables them talking to each other.
func (r RandomizedTopology) GenerateFanout(ids flow.IdentityList) (flow.IdentityList, error) {
	myChannelIDs := r.subMngr.GetChannelIDs()
	if len(myChannelIDs) == 0 {
		// no subscribed channel id, hence skip topology creation
		// we do not return an error at this state as invocation of MakeTopology may happen before
		// node subscribing to all its channels.
		return flow.IdentityList{}, nil
	}

	var myFanout flow.IdentityList

	// generates a randomized subgraph per channel
	for _, myChannel := range myChannelIDs {
		topicFanout, err := r.subsetChannel(ids, myChannel)
		if err != nil {
			return nil, fmt.Errorf("failed to derive list of peer nodes to connect for topic %s: %w", myChannel, err)
		}
		myFanout = myFanout.Union(topicFanout)
	}

	if len(myFanout) == 0 {
		return nil, fmt.Errorf("topology size reached zero")
	}
	return myFanout, nil
}

// subsetChannel returns a random subset of the identity list that is passed. `shouldHave` represents set of
// identities that should be included in the returned subset.
// Returned identities should all subscribed to the specified `channel`.
// Note: this method should not include identity of its executor.
func (r RandomizedTopology) subsetChannel(ids flow.IdentityList, channel string) (flow.IdentityList, error) {
	// excludes node itself
	sampleSpace := ids.Filter(filter.Not(filter.HasNodeID(r.me)))

	if _, ok := engine.IsClusterChannelID(channel); ok {
		return r.clusterChannelHandler(sampleSpace)
	} else {
		return r.nonClusterChannelHandler(sampleSpace, channel)
	}
}

// sampleFanout receives two lists: all and shouldHave. It then samples a connected fanout
// for the caller that includes the shouldHave set. Independent invocations of this method over
// different nodes, should create a connected graph.
// Fanout is the set of nodes that this instance should get connected to in order to create a
// connected graph.
func (r RandomizedTopology) sampleFanout(all flow.IdentityList) (flow.IdentityList, error) {
	if len(all) == 0 {
		return nil, fmt.Errorf("empty identity list")
	}

	fanout := flow.IdentityList{}
	for _, id := range all {
		if r.tossBiasedBit() {
			fanout = append(fanout, id)
		}
	}

	return fanout, nil
}

// clusterPeers returns the list of other nodes within the same cluster as this node.
func (r RandomizedTopology) clusterPeers() (flow.IdentityList, error) {
	currentEpoch := r.state.Final().Epochs().Current()
	clusterList, err := currentEpoch.Clustering()
	if err != nil {
		return nil, fmt.Errorf("failed to extract cluster list %w", err)
	}

	myCluster, _, found := clusterList.ByNodeID(r.me)
	if !found {
		return nil, fmt.Errorf("failed to find the cluster for node ID %s", r.me.String())
	}

	return myCluster, nil
}

// clusterChannelHandler returns a connected graph fanout of peers in the same cluster as executor of this instance.
func (r RandomizedTopology) clusterChannelHandler(ids flow.IdentityList) (flow.IdentityList, error) {
	// extracts cluster peer ids to which the node belongs to.
	clusterPeers, err := r.clusterPeers()
	if err != nil {
		return nil, fmt.Errorf("failed to find cluster peers for node %s: %w", r.me.String(), err)
	}

	// excludes node itself from cluster
	clusterPeers = clusterPeers.Filter(filter.Not(filter.HasNodeID(r.me)))

	// checks all cluster peers belong to the passed ids list
	nonMembers := clusterPeers.Filter(filter.Not(filter.In(ids)))
	if len(nonMembers) > 0 {
		return nil, fmt.Errorf("cluster peers not belonged to sample space: %v", nonMembers)
	}

	// samples fanout from cluster peers
	return r.sampleFanout(clusterPeers)
}

// clusterChannelHandler returns a connected graph fanout of peers from `ids` that subscribed to `channel`.
// The returned sample contains `shouldHave` ones that also subscribed to `channel`.
func (r RandomizedTopology) nonClusterChannelHandler(ids flow.IdentityList, channel string) (flow.IdentityList, error) {
	if _, ok := engine.IsClusterChannelID(channel); ok {
		return nil, fmt.Errorf("could not handle cluster channel: %s", channel)
	}

	// extracts flow roles subscribed to topic.
	roles, ok := engine.RolesByChannelID(channel)
	if !ok {
		return nil, fmt.Errorf("unknown topic with no subscribed roles: %s", channel)
	}

	// samples fanout among interacting roles
	return r.sampleFanout(ids.Filter(filter.HasRole(roles...)))
}

// tossBiasedBit returns true with probability equal `r.chance/100`, and returns false otherwise.
func (r *RandomizedTopology) tossBiasedBit() bool {
	if r.rng.UintN(100) <= r.chance {
		return true
	} else {
		return false
	}
}