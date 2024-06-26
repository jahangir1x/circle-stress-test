package serializer

type LoginReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	LongLived bool   `json:"long_lived"`
}

type LoginResp struct {
	UserId       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CsvData struct {
	LocationName string  `csv:"Location name"`
	Coordinates  string  `csv:"Coordinates"`
	Radius       float64 `csv:"Radius (Meters)"`
}

type CreatePlaceReq struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CliArgs struct {
	Operation      string
	Start          int
	End            int
	EmailPrefix    string
	EmailSuffix    string
	Password       string
	LocationSpread float64
	WaitTimeMin    float64
	WaitTimeMax    float64
}
