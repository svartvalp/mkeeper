package mkeeper

type Stats struct {
	OnPut    func()
	OnDelete func()
	OnGet    func()
	OnMiss   func()
}

func (s *Stats) ProcessEvents(events []Event) {
	for _, e := range events {
		switch e.Type {
		case Put:
			s.OnPut()
		case Get:
			s.OnGet()
		case Miss:
			s.OnMiss()
		case Delete:
			s.OnDelete()
		}
	}
}
