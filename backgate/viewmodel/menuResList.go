package viewmodel

import pb "backgate/training"

type MenuResList struct {
	Menu          []*pb.MenuRes `json:"menu"`
	Permissions   []string      `json:"permissions"`
	DashboardGrid []string      `json:"dashboardGrid"`
}
