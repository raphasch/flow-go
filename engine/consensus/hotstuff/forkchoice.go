package hotstuff

import (
	"bytes"

	"github.com/dapperlabs/flow-go/engine/consensus/hotstuff/types"
)

type ForkChoice struct {
	lockedBlock        *types.BlockProposal
	finalizedBlock     *types.BlockProposal
	highestQC          *types.QuorumCertificate
	incorporatedBlocks *LevelledForrest
}

// Queries

func (fc *ForkChoice) GetQCForNextBlock(view uint64) *types.QuorumCertificate {
	return fc.highestQC
}

func (fc *ForkChoice) FindProposalsByView(view uint64) []*types.BlockProposal {
	vertices := fc.incorporatedBlocks.FindVerticesByLevel(view)

	proposals := make([]*types.BlockProposal, len(vertices))
	for i, v := range vertices {
		proposal, ok := v.(*types.BlockProposal)
		if !ok {
			return nil
		}
		proposals[i] = proposal
	}
	return proposals
}

func (fc *ForkChoice) FindProposalByViewAndBlockMRH(view uint64, blockMRH types.MRH) (*types.BlockProposal, bool) {
	vertex, exists := fc.incorporatedBlocks.FindVerticeByID(blockMRH)
	if !exists {
		return nil, false
	}

	proposal, ok := vertex.(*types.BlockProposal)
	if !ok {
		return nil, false
	}

	if proposal.View() != view {
		return nil, false
	}
	return proposal, true
}

func (fc *ForkChoice) FinalizedView() uint64 {
	return fc.finalizedBlock.Block.View
}

func (fc *ForkChoice) LockedQC() *types.QuorumCertificate {
	return fc.lockedBlock.Block.QC
}

// LockedView is the view number of the locked QC
func (fc *ForkChoice) LockedView() uint64 {
	return fc.LockedQC().View
}

func (fc *ForkChoice) IsSafeNode(proposal *types.BlockProposal) bool {
	qc := proposal.Block.QC
	if qc.View > fc.LockedView() {
		return true
	}
	if qc.View < fc.LockedView() {
		return false
	}
	return bytes.Equal(qc.BlockMRH, fc.LockedQC().BlockMRH)
}

// return true only if the following conditions are all true
// 1. above the finalized block's View
// 2. block's QC is pointing to a leaf block on the tree
func (fc *ForkChoice) CanIncorporate(bp *types.BlockProposal) bool {
	return fc.incorporatedBlocks.CanIncorporate(bp)
}

// Updates

func (fc *ForkChoice) UpdateValidQC(qc *types.QuorumCertificate) (highestQCUpdated bool, finalizedBlocks []*types.BlockProposal) {
	// updateHighestQC
	// updateLockedQC
	// updateFinalizedBlock
	if qc.View > fc.highestQC.View {
		fc.highestQC = qc
		highestQCUpdated = true
	}
	if qc.View > fc.finalizedBlock.View() {
		parent, grantParent, greatGrantParent := fc.incorporatedBlocks.FindThreeParentsOf(qc.BlockMRH)
		finalizedQC, finalized := isFinalized(parent, grantParent, greatGrantParent)
		if finalized {
			finalizedBlocks := fc.findFinalizedBlocks(finalizedQC)
		}
	}
}

func (fc *ForkChoice) AddNewProposal(bp *types.BlockProposal) (incorperatedBlock *types.BlockProposal, added bool) {
	// TODO: UpdateValidQC
	return fc.incorporatedBlocks.AddIncorporatableVertex(bp)
}
