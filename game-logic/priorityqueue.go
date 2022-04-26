package game

type Node struct {
	Pos
	priority int
}

type priorityQueue []Node

func (p priorityQueue) push(pos Pos, priority int) priorityQueue {
	newNode := Node{pos, priority}
	p = append(p, newNode)
	newNodeIdx := len(p)-1
	pIndex, pNode := p.getParent(newNodeIdx)

	for newNode.priority < pNode.priority && newNodeIdx > 0 {
		p.Swap(newNodeIdx, pIndex)
		newNodeIdx = pIndex
		pIndex, pNode = p.getParent(newNodeIdx)
	}
	return p
}

func (p priorityQueue) pop() (priorityQueue, Pos) {
	index := 0
	result := p[index].Pos
	p.Swap(index, len(p)-1)
	p = p[:len(p)-1]

	if len(p) == 0 {
		return p, result
	}

	rootNow := p[index]

	leftExist, leftIndex, left := p.getLeftChild(index)
	rightExist, rightIndex, right := p.getRightChild(index)

	for (leftExist && rootNow.priority > left.priority) || 
	(rightExist && rootNow.priority > right.priority) {
		if !rightExist || left.priority <= right.priority {
			p.Swap(index, leftIndex)
			index = leftIndex
		} else {
			p.Swap(index, rightIndex)
			index = rightIndex
		}

		leftExist, leftIndex, left = p.getLeftChild(index)
		rightExist, rightIndex, right = p.getRightChild(index)
	}

	return p, result
}

func (p priorityQueue) getParent(index int) (int, Node) {
	return (index-1)/2, p[(index-1)/2]
}

func (p priorityQueue) getLeftChild(index int) (bool, int, Node) {
	i := (index*2)+1
	if i < len(p) {
		return true, i, p[i]
	}
	return false, -1, Node{}
}

func (p priorityQueue) getRightChild(index int) (bool, int, Node) {
	i := (index*2)+2
	if i < len(p) {
		return true, i, p[i]
	}
	return false, -1, Node{}
}

func (p priorityQueue) Len() int           { return len(p) }
func (p priorityQueue) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p priorityQueue) Less(i, j int) bool { return p[i].priority < p[j].priority }
