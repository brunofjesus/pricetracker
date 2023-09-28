package store

type WortenCategoriesResponse struct {
	Data struct {
		SolveURL struct {
			Typename  string `json:"__typename"`
			ID        string `json:"_id"`
			Code      string `json:"code"`
			Success   bool   `json:"success"`
			Message   any    `json:"message"`
			ContextID string `json:"contextId"`
			Layout    struct {
				Typename string `json:"__typename"`
				ID       string `json:"_id"`
				Modules  []struct {
					Typename   string `json:"__typename"`
					Order      int    `json:"order"`
					Priority   int    `json:"priority"`
					TargetedBy string `json:"targetedBy"`
					Targets    string `json:"targets"`
					Refs       []struct {
						Typename string `json:"__typename"`
						Key      string `json:"_key"`
						Type     string `json:"_type"`
						Valid    bool   `json:"valid"`
						ID       string `json:"id"`
						URL      string `json:"url"`
					} `json:"refs"`
				} `json:"modules"`
			} `json:"layout"`
		} `json:"solveURL"`
	} `json:"data"`
}
