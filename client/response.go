package client

/*
	**Payload** - Response as returned by the upstream server.
	Nothing to be changed here, unless very sure to do so
*/

type Team struct {
	ID          int    `json:"id"`
	OptaID      int    `json:"optaId"`
	Country     string `json:"country"`
	CountryName string `json:"countryName"`
	Name        string `json:"name"`
	LogoUrls    []struct {
		Size string `json:"size"`
		URL  string `json:"url"`
	} `json:"logoUrls"`
	IsNational      bool `json:"isNational"`
	HasOfficialPage bool `json:"hasOfficialPage"`
	Competitions    []struct {
		CompetitionID   int    `json:"competitionId"`
		CompetitionName string `json:"competitionName"`
	} `json:"competitions"`
	Players   []Player `json:"players"`
	Officials []struct {
		CountryName  string `json:"countryName"`
		ID           string `json:"id"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Country      string `json:"country"`
		Position     string `json:"position"`
		ThumbnailSrc string `json:"thumbnailSrc"`
		Affiliation  struct {
			Name         string `json:"name"`
			ThumbnailSrc string `json:"thumbnailSrc"`
		} `json:"affiliation"`
	} `json:"officials"`
	Colors struct {
		ShirtColorHome string `json:"shirtColorHome"`
		ShirtColorAway string `json:"shirtColorAway"`
		CrestMainColor string `json:"crestMainColor"`
		MainColor      string `json:"mainColor"`
	} `json:"colors"`
}

type Player struct {
	ID           string `json:"id"`
	Country      string `json:"country"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Name         string `json:"name"`
	Position     string `json:"position"`
	Number       int    `json:"number"`
	BirthDate    string `json:"birthDate"`
	Age          string `json:"age"`
	Height       int    `json:"height"`
	Weight       int    `json:"weight"`
	ThumbnailSrc string `json:"thumbnailSrc"`
	Affiliation  struct {
		Name         string `json:"name"`
		ThumbnailSrc string `json:"thumbnailSrc"`
	} `json:"affiliation"`
}

type Data struct {
	Team Team `json:"team"`
}

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}
