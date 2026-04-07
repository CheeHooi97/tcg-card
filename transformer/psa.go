package transformer

type PSAGetCert struct {
	PSACert PSACert `json:"PSACert"`
	DNACert DNACert `json:"DNACert"`
}

type PSACert struct {
	CertNumber                   string `json:"CertNumber"`
	SpecID                       int64  `json:"SpecID"`
	SpecNumber                   string `json:"SpecNumber"`
	LabelType                    string `json:"LabelType"`
	ReverseBarCode               bool   `json:"ReverseBarCode"`
	Year                         string `json:"Year"`
	Brand                        string `json:"Brand"`
	Category                     string `json:"Category"`
	CardNumber                   string `json:"CardNumber"`
	Subject                      string `json:"Subject"`
	Variety                      string `json:"Variety"`
	IsPSADNA                     bool   `json:"IsPSADNA"`
	IsDualCert                   bool   `json:"IsDualCert"`
	GradeDescription             string `json:"GradeDescription"`
	CardGrade                    string `json:"CardGrade"`
	TotalPopulation              int64  `json:"TotalPopulation"`
	TotalPopulationWithQualifier int64  `json:"TotalPopulationWithQualifier"`
	PopulationHigher             int64  `json:"PopulationHigher"`
}

type DNACert struct {
	CertNumber           string `json:"CertNumber"`
	ItemDescription      string `json:"ItemDescription"`
	PrimarySubjects      []any  `json:"PrimarySubjects"`
	AuthenticationResult string `json:"AuthenticationResult"`
	DNAItemType          string `json:"DNAItemType"`
}

type PSAGetSpecPopulation struct {
	SpecID      int64  `json:"SpecID"`
	Description string `json:"Description"`
	PSAPop      PSAPop `json:"PSAPop"`
	PSADNAPop   PSAPop `json:"PSADNAPop"`
}

type PSAPop struct {
	Total     int64 `json:"Total"`
	Auth      int64 `json:"Auth"`
	Grade1    int64 `json:"Grade1"`
	Grade1Q   int64 `json:"Grade1Q"`
	Grade1_5  int64 `json:"Grade1_5"`
	Grade1_5Q int64 `json:"Grade1_5Q"`
	Grade2    int64 `json:"Grade2"`
	Grade2Q   int64 `json:"Grade2Q"`
	Grade2_5  int64 `json:"Grade2_5"`
	Grade3    int64 `json:"Grade3"`
	Grade3Q   int64 `json:"Grade3Q"`
	Grade3_5  int64 `json:"Grade3_5"`
	Grade4    int64 `json:"Grade4"`
	Grade4Q   int64 `json:"Grade4Q"`
	Grade4_5  int64 `json:"Grade4_5"`
	Grade5    int64 `json:"Grade5"`
	Grade5Q   int64 `json:"Grade5Q"`
	Grade5_5  int64 `json:"Grade5_5"`
	Grade6    int64 `json:"Grade6"`
	Grade6Q   int64 `json:"Grade6Q"`
	Grade6_5  int64 `json:"Grade6_5"`
	Grade7    int64 `json:"Grade7"`
	Grade7Q   int64 `json:"Grade7Q"`
	Grade7_5  int64 `json:"Grade7_5"`
	Grade8    int64 `json:"Grade8"`
	Grade8Q   int64 `json:"Grade8Q"`
	Grade8_5  int64 `json:"Grade8_5"`
	Grade9    int64 `json:"Grade9"`
	Grade9Q   int64 `json:"Grade9Q"`
	Grade10   int64 `json:"Grade10"`
}
