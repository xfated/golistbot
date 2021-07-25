package constants

// Finite state machine for handling adding items
type State int

const (
	Idle State = iota

	/* #### Adding place #### */
	ReadyForNextAction
	SetName
	SetAddress
	SetNotes
	SetURL
	SetImages
	SetTags
	/* ######## */

	/* #### Query #### */
	QuerySelectType

	QueryOneTagOrName
	QueryOneSetName

	QuerySetTags
	QueryFewSetNum
	QueryRetrieve
	/* ######## */

	Finished
)

type PlaceDetails struct {
	Name    string          `json:"name"`
	Address string          `json:"address"`
	Notes   string          `json:"notes"`
	URL     string          `json:"url"`
	Images  map[string]bool `json:"images"`
	Tags    map[string]bool `json:"tags"`
}

func (placeData *PlaceDetails) GetImageIDs() []string {
	if placeData.Images == nil {
		return []string{}
	}
	imageIDs := make([]string, 0)
	for id := range placeData.Images {
		imageIDs = append(imageIDs, id)
	}
	return imageIDs
}

func IsAddingNewPlace(state State) bool {
	switch state {
	case ReadyForNextAction,
		SetName,
		SetAddress,
		SetNotes,
		SetURL,
		SetImages,
		SetTags:
		return true
	default:
		return false
	}
}

func IsQuery(state State) bool {
	switch state {
	case QuerySelectType,

		QueryOneTagOrName,
		QueryOneSetName,

		QueryFewSetNum,

		QuerySetTags,
		QueryRetrieve:
		return true
	default:
		return false
	}
}
