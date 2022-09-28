package models

type Relations struct {
	Concerts []struct {
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

// Artist is using to save all the info about an artist
type Artist struct {
	ID             int      `json:"id"`
	Image          string   `json:"image"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	DatesLocations map[string][]string
}

type FilterParams struct {
	CreationDate    [2]string
	FirstAlbumDate  [2]string
	MembersCheckbox []string
	Location        string
}
