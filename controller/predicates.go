package controller

import (
	"log"
	"math/rand"
	"strings"

	"k8s.io/api/core/v1"
	extender "k8s.io/kube-scheduler/extender/v1"
)

const (
	LuckyPred        = "Lucky"
	LuckyPredFailMsg = "Sorry, you're not lucky"
)

var predicatesFuncs = map[string]FitPredicate{
	LuckyPred: LuckyPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{LuckyPred}

// it's webhooked to pkg/scheduler/core/generic_scheduler.go#findNodesThatFit()
func filter(args extender.ExtenderArgs) *extender.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(extender.FailedNodesMap)
	pod := args.Pod
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := extender.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	var failReasons []string
	for _, predicateKey := range predicatesSorted {
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

// 简单地检查随机数是否为偶数来判断是否批准节点，如果是的话我们就认为这是一个幸运的节点，否则拒绝批准该节点
func LuckyPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	lucky := rand.Intn(2) == 0
	if lucky {
		log.Printf("pod %v/%v is lucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v is unlucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name)
	return false, []string{LuckyPredFailMsg}, nil
}
