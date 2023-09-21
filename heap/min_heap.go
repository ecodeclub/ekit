package heap

import (
	"errors"
	"github.com/ecodeclub/ekit"
)

var (
	ErrMinHeapComparatorIsNull = errors.New("ekit: MinHeap 的 Comparator 不能为 nil")
	ErrMinHeapIsEmpty          = errors.New("ekit: MinHeap 中没有元素")
	ErrValueNotInMinHeap       = errors.New("ekit: 要查找的元素不在 MinHeap 内")
)

// MinHeap 小根堆
type MinHeap[T any] struct {
	compare ekit.Comparator[T] //堆元素怎么比大小
	data    []T                //堆内数据
	size    int                //堆内元素数量
}

func NewMinHeap[T any](compare ekit.Comparator[T], initSize int) (*MinHeap[T], error) {
	if compare == nil {
		return nil, ErrMinHeapComparatorIsNull
	}

	return &MinHeap[T]{
		compare: compare,
		data:    make([]T, 0, initSize),
		size:    0,
	}, nil
}

// getParentIndex 获取父结点下标
// 假设有一个小根堆 0,1,2,3,4,5,6
// (1-1)/2=>0;(2-1)/2=>0.5(0);
// (3-1)/2=>1;(4-1)/2=>1.5(1);(5-1)/2=>2;(6-1)/2=>2.5(2);
func (mh *MinHeap[T]) getParentIndex(index int) int {
	return (index - 1) / 2
}

// getLeftChildIndex 获取左子结点下标
// 假设有一个小根堆 0,1,2,3,4,5,6
// 2*0+1=>1;
// 2*1+1=>3;2*2+1=>5
func (mh *MinHeap[T]) getLeftChildIndex(index int) int {
	return 2*index + 1
}

// getRightChildIndex 获取右子结点下标
// 假设有一个小根堆 0,1,2,3,4,5,6
// 2*0+2=>2;
// 2*1+2=>4;2*2+2=>6
func (mh *MinHeap[T]) getRightChildIndex(index int) int {
	return 2*index + 2
}

// Size 堆内元素数量
func (mh *MinHeap[T]) Size() int {
	return mh.size
}

// Add 插入元素
func (mh *MinHeap[T]) Add(value T) {
	mh.data = append(mh.data, value)
	mh.size++
	mh.heapifyUp(mh.size - 1)
}

// GetTop 获取堆顶元素，不移除元素
func (mh *MinHeap[T]) GetTop() (T, error) {
	if mh.size < 1 {
		var t T
		return t, ErrMinHeapIsEmpty
	}
	return mh.data[0], nil
}

// ExtractTop 提取堆顶元素，移除元素
func (mh *MinHeap[T]) ExtractTop() (T, error) {
	if mh.size < 1 {
		var t T
		return t, ErrMinHeapIsEmpty
	}
	return mh.extractSubHeapTop(0), nil
}

// extractSubHeapTop 提取子堆堆顶元素
func (mh *MinHeap[T]) extractSubHeapTop(index int) T {
	min := mh.data[index] //提取堆顶元素

	lastIndex := mh.size - 1
	mh.data[index] = mh.data[lastIndex] //把最后一个元素放到堆顶
	mh.data = mh.data[:lastIndex]       //移除最后一个元素
	mh.size--                           //修改堆内元素数量

	mh.heapifyDown(index) //从堆顶开始向下调整堆结构

	return min
}

// Delete 移除堆内任意一个元素
// 思路和移除堆顶元素是一样的，多了一个查询的步骤
// 要删除的元素就相当于小根堆里面一个子堆的堆顶
func (mh *MinHeap[T]) Delete(value T) error {
	// 查询要删除的元素的位置
	index := -1
	for i, v := range mh.data {
		if 0 == mh.compare(v, value) {
			index = i
			break
		}
	}
	if index == -1 {
		return ErrValueNotInMinHeap
	}

	mh.extractSubHeapTop(index)

	return nil
}

// heapifyUp 向上调整堆结构
func (mh *MinHeap[T]) heapifyUp(index int) {
	for index > 0 && mh.compare(mh.data[index], mh.data[mh.getParentIndex(index)]) == -1 {
		//比较index和其父结点，如果index比父结点小，就交换
		mh.data[index], mh.data[mh.getParentIndex(index)] = mh.data[mh.getParentIndex(index)], mh.data[index]
		index = mh.getParentIndex(index)
	}
}

// heapifyDown 向下调整堆结构
func (mh *MinHeap[T]) heapifyDown(index int) {
	//比较index和其两个子结点，把index和子结点中最小的那个交换
	minIndex := index
	leftChildIndex := mh.getLeftChildIndex(index)
	if leftChildIndex < mh.size && mh.compare(mh.data[leftChildIndex], mh.data[minIndex]) == -1 {
		minIndex = leftChildIndex
	}
	rightChildIndex := mh.getRightChildIndex(index)
	if rightChildIndex < mh.size && mh.compare(mh.data[rightChildIndex], mh.data[minIndex]) == -1 {
		minIndex = rightChildIndex
	}
	if index != minIndex {
		//如果发生了交换就要继续向下调整堆结构
		mh.data[index], mh.data[minIndex] = mh.data[minIndex], mh.data[index]
		mh.heapifyDown(minIndex)
	}
}
