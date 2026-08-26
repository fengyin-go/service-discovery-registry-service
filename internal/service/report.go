package service

// RegionReport 区域维度报告。
type RegionReport struct {
	ByRegion map[string]int `json:"by_region"`
	ByZone   map[string]int `json:"by_zone"`
}

// RegionReport 统计实例在区域与可用区上的分布。
func (s *Service) RegionReport() RegionReport {
	r := RegionReport{ByRegion: make(map[string]int), ByZone: make(map[string]int)}
	for _, i := range s.store.ListInstances() {
		if i.Region != "" {
			r.ByRegion[i.Region]++
		}
		if i.Zone != "" {
			r.ByZone[i.Zone]++
		}
	}
	return r
}

// AlgorithmReport 负载均衡算法分布报告。
type AlgorithmReport struct {
	ByAlgorithm map[string]int `json:"by_algorithm"`
}

// AlgorithmReport 统计负载均衡算法使用情况。
func (s *Service) AlgorithmReport() AlgorithmReport {
	r := AlgorithmReport{ByAlgorithm: make(map[string]int)}
	for _, l := range s.store.ListLoadBalanceStrategies() {
		r.ByAlgorithm[l.Algorithm]++
	}
	return r
}

// EventTypeReport 事件类型分布。
func (s *Service) EventTypeReport() map[string]int {
	by := make(map[string]int)
	for _, e := range s.store.ListRegistryEvents() {
		by[e.EventType]++
	}
	return by
}
