/* The Computer Language Benchmarks Game
 * http://benchmarksgame.alioth.debian.org/
 *
 * based on Go program by The Go Authors.
 * based on C program by Kevin Carson
 * flag.Arg hack by Isaac Gouy
 * modified by victor<ldb8002@163.com> rewrite the bottomUpTree function
 * modified by Jamil Djadala to use goroutines
 * modified by Chai Shushan
 * *reset*
 *
 */

package main

import (
	"flag"
	"fmt"
	//"runtime/debug"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

var minDepth = 4
var n = 0

func main() {
	defer func(t time.Time) {
		fmt.Println(time.Since(t))
	}(time.Now())

	//debug.SetGCPercent(300)

	flag.Parse()
	if flag.NArg() > 0 {
		n, _ = strconv.Atoi(flag.Arg(0))
	}

	maxDepth := n
	if minDepth+2 > n {
		maxDepth = minDepth + 2
	}
	stretchDepth := maxDepth + 1

	check_l := bottomUpTree(stretchDepth).ItemCheck()
	fmt.Printf("stretch tree of depth %d\t check: %d\n", stretchDepth, check_l)

	longLivedTree := bottomUpTree(maxDepth)

	result_trees := make([]int, maxDepth+1)
	result_check := make([]int, maxDepth+1)

	var wg sync.WaitGroup
	for depth_l := minDepth; depth_l <= maxDepth; depth_l += 2 {
		wg.Add(1)
		go func(depth int) {
			iterations := 1 << uint(maxDepth-depth+minDepth)
			check := 0

			for i := 1; i <= iterations; i++ {
				check += bottomUpTree(depth).ItemCheck()
			}
			result_trees[depth] = iterations
			result_check[depth] = check

			wg.Done()
		}(depth_l)
	}
	wg.Wait()

	for depth := minDepth; depth <= maxDepth; depth += 2 {
		fmt.Printf("%d\t trees of depth %d\t check: %d\n",
			result_trees[depth], depth, result_check[depth],
		)
	}
	fmt.Printf("long lived tree of depth %d\t check: %d\n",
		maxDepth, longLivedTree.ItemCheck(),
	)
}

func bottomUpTree(depth int) *Node {
	depth++
	num := 1<<uint(depth) - 1
	nodes := make([]*Node, num)
	nodeSize := int(unsafe.Sizeof(Node{}))
	mem := make([]byte, num*nodeSize)
	ptr := (*[3]uintptr)(unsafe.Pointer(&mem))[0]
	for i := 0; i < depth; i++ {
		currentLayerNum := 1 << uint(i)
		currentLayerFirst := currentLayerNum - 1
		nodes[currentLayerFirst] = (*Node)(unsafe.Pointer(ptr))
		ptr += nodeSize
		if currentLayerNum <= 1 {
			continue
		}
		upperLayerFirst := currentLayerNum>>1 - 1
		nodes[upperLayerFirst].left = nodes[currentLayerFirst]
		for m := currentLayerFirst + 1; m <= 2*currentLayerNum-2; m++ {
			nodes[m] = (*Node)(unsafe.Pointer(ptr))
			ptr += nodeSize
			diff := m - currentLayerFirst
			mod := diff % 2
			offset := diff / 2
			parentNodeIndex := upperLayerFirst + offset
			if mod == 0 {
				nodes[parentNodeIndex].left = nodes[m]
			} else {
				nodes[parentNodeIndex].right = nodes[m]
			}
		}
	}
	return nodes[0]
}

type Node struct {
	left, right *Node
}

func (self *Node) ItemCheck() int {
	if self.left == nil {
		return 1
	}
	return 1 + self.left.ItemCheck() + self.right.ItemCheck()
}
