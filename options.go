// Author: yann
// Date: 2022/5/28
// Desc: consistent

package consistent

type Option func(consistent *Consistent)

func ReplicasOption(replicas int) func(consistent *Consistent) {
	return func(consistent *Consistent) {
		if replicas == 0 {
			replicas = defaultReplicas
		}
		consistent.replicas = replicas
	}
}

func HashOption(hash Hash) func(consistent *Consistent) {
	return func(consistent *Consistent) {
		if hash == nil {
			hash = defaultHash
		}
		consistent.hash = hash
	}
}
