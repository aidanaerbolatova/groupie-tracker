package store

import (
	"encoding/json"
	"errors"
	"groupie-tracker/internal/models"
	"net/http"
	"strconv"
	"strings"
)

const (
	artistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	relationsURL = "https://groupietrackers.herokuapp.com/api/relation"
)

type Store struct {
	AllArtists []models.Artist
	Result     []models.Artist
}

// fucntion to decode json data
func GetJson(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(result)
}

// GetAllArtists is using to put the artists' data to the fields
func (s *Store) GetAllArtists() error {
	if len(s.AllArtists) != 0 {
		return nil
	}

	if err := GetJson(artistsURL, &s.AllArtists); err != nil {
		return err
	}

	var relations models.Relations
	if err := GetJson(relationsURL, &relations); err != nil {
		return err
	}
	for i, v := range relations.Concerts {
		s.AllArtists[i].DatesLocations = v.DatesLocations
	}

	return nil
}

func (s *Store) GetArtistByID(id int) (models.Artist, error) {
	if id < 1 || id > len(s.AllArtists) {
		return models.Artist{}, errors.New("not found")
	}
	return s.AllArtists[id-1], nil
}

func (s *Store) GetSearchResult(toSearch string) []models.Artist {
	var result []models.Artist
	for _, v := range s.AllArtists {
		if contains(v.Name, toSearch) ||
			contains(strconv.Itoa(v.CreationDate), toSearch) ||
			contains(v.FirstAlbum, toSearch) {
			result = append(result, v)
			continue
		}

		var isContains bool
		for _, member := range v.Members {
			if contains(member, toSearch) {
				result = append(result, v)
				isContains = true
				break
			}
		}

		if isContains {
			continue
		}

		for place, dates := range v.DatesLocations {
			if contains(place, toSearch) {
				result = append(result, v)
				break
			}
			for _, date := range dates {
				if contains(date, toSearch) {
					result = append(result, v)
					isContains = true
					break
				}
			}

			if isContains {
				break
			}
		}
	}
	return result
}

func contains(s, substr string) bool {
	return strings.Contains(
		strings.ToLower(s),
		strings.ToLower(substr),
	)
}

func (s *Store) GetFilterResult(f models.FilterParams) ([]models.Artist, error) {
	var result []models.Artist

	crDateFrom, err1 := strconv.Atoi(f.CreationDate[0])
	crDateTo, err2 := strconv.Atoi(f.CreationDate[1])
	albumDateFrom, err3 := strconv.Atoi(f.FirstAlbumDate[0])
	albumDateTo, err4 := strconv.Atoi(f.FirstAlbumDate[1])

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return nil, errors.New("bad request to filter: not a year")
	}

	for _, artist := range s.AllArtists {
		var isInRange bool
		parts := strings.Split(artist.FirstAlbum, "-")
		n, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}
		if checkRange(n, albumDateFrom, albumDateTo) &&
			checkRange(artist.CreationDate, crDateFrom, crDateTo) {
			isInRange = true
		}

		var isMemberRange bool
		for _, v := range f.MembersCheckbox {
			num, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("bad request to filter: not a year")
			}
			if num == len(artist.Members) {
				isMemberRange = true
				break
			}
			if num == 9 && len(artist.Members) > 8 {
				isMemberRange = true
				break
			}
		}

		if len(f.MembersCheckbox) == 0 {
			isMemberRange = true
		}

		var isLocation bool
		for place := range artist.DatesLocations {
			place = strings.ReplaceAll(place, "-", ", ")
			correctLocation := strings.ReplaceAll(place, "_", "-")
			if contains(correctLocation, f.Location) {
				isLocation = true
				break
			}
		}

		if len(f.Location) == 0 {
			isLocation = true
		}

		if isInRange && isMemberRange && isLocation {
			result = append(result, artist)
		}
	}
	return result, nil
}

func checkRange(target, from, to int) bool {
	if target >= from && target <= to {
		return true
	}
	return false
}
