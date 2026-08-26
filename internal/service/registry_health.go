package service

// HealthScore 注册中心整体健康评分。
type HealthScore struct {
	Score        float64 `json:"score"`
	Grade        string  `json:"grade"`
	HealthyRatio float64 `json:"healthy_ratio"`
	StaleCount   int     `json:"stale_count"`
	CycleCount   int     `json:"cycle_count"`
}

// HealthScore 综合健康率、心跳超时与依赖环计算健康分（0-100）。
func (s *Service) HealthScore() HealthScore {
	stats := s.StatsRegistry()
	stale := s.StaleInstances(30)
	cycles := s.DetectCycles()
	score := stats.HealthyRatio * 100
	score -= float64(len(stale)) * 2
	score -= float64(len(cycles)) * 10
	if score < 0 {
		score = 0
	}
	grade := "D"
	switch {
	case score >= 90:
		grade = "A"
	case score >= 75:
		grade = "B"
	case score >= 60:
		grade = "C"
	}
	return HealthScore{
		Score:        score,
		Grade:        grade,
		HealthyRatio: stats.HealthyRatio,
		StaleCount:   len(stale),
		CycleCount:   len(cycles),
	}
}
