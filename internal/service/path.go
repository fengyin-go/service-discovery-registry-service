package service

// DependencyPath 依赖路径。
type DependencyPath struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Path   []string `json:"path"`
	Exists bool     `json:"exists"`
}

// FindDependencyPath 使用 BFS 查找 from 到 to 的最短依赖路径。
func (s *Service) FindDependencyPath(fromID, toID string) DependencyPath {
	adj := make(map[string][]string)
	for _, d := range s.store.ListDependencies() {
		adj[d.ServiceID] = append(adj[d.ServiceID], d.DependsOnID)
	}
	if fromID == toID {
		return DependencyPath{From: fromID, To: toID, Path: s.serviceNames([]string{fromID}), Exists: true}
	}
	prev := make(map[string]string)
	visited := map[string]bool{fromID: true}
	queue := []string{fromID}
	found := false
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == toID {
			found = true
			break
		}
		for _, next := range adj[cur] {
			if visited[next] {
				continue
			}
			visited[next] = true
			prev[next] = cur
			queue = append(queue, next)
		}
	}
	if !found {
		return DependencyPath{From: fromID, To: toID, Path: []string{}, Exists: false}
	}
	path := []string{toID}
	for cur := toID; cur != fromID; {
		p, ok := prev[cur]
		if !ok {
			break
		}
		path = append([]string{p}, path...)
		cur = p
	}
	return DependencyPath{From: fromID, To: toID, Path: s.serviceNames(path), Exists: true}
}
