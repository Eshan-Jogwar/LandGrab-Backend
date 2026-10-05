package pipeline

// BuildMoveChain links a slice of MoveNode instances sequentially and returns the head node.
func BuildMoveChain(nodes ...MoveNode) MoveNode {
	if len(nodes) == 0 {
		return nil
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].SetNext(nodes[i+1])
	}
	return nodes[0]
}

// BuildDecisionChain links a slice of DecisionNode instances sequentially and returns the head node.
func BuildDecisionChain(nodes ...DecisionNode) DecisionNode {
	if len(nodes) == 0 {
		return nil
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].SetNext(nodes[i+1])
	}
	return nodes[0]
}
