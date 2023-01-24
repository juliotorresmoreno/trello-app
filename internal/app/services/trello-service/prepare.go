package trello_service

func (e TrelloService) Prepare() error {
	labels, err := e.GetLabels()
	if err != nil {
		return err
	}
	var predefined = []CreateCardLabel{
		{"task", "blue"},
		{"bug", "red"},
		{"issue", "green"},
		{"maintenance", "purple"},
		{"test", "black"},
	}
	for _, v := range predefined {
		if _, ok := labels[v.Name]; !ok {
			err = e.CreateLabel(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
