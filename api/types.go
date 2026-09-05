package api

import "encoding/json"

// APIResponse represents the top-level response from the MalwareBazaar API
type APIResponse struct {
	QueryStatus string          `json:"query_status"`
	Data        []MalwareSample `json:"data,omitempty"`
}

// CSCBResponse represents the response for get_cscb
type CSCBResponse struct {
	QueryStatus string      `json:"query_status"`
	Data        []CSCBEntry `json:"data,omitempty"`
}

// CSCBEntry represents an entry in the Code Signing Certificate Blocklist.
// valid_from and valid_to are kept as strings: the API returns ISO-8601
// (2026-03-01T00:04:43Z) while the documentation shows "Jul 2 00:00:00 2020 GMT".
type CSCBEntry struct {
	TimeStamp           string `json:"time_stamp"`
	SerialNumber        string `json:"serial_number"`
	Thumbprint          string `json:"thumbprint"`
	ThumbprintAlgorithm string `json:"thumbprint_algorithm"`
	SubjectCN           string `json:"subject_cn"`
	IssuerCN            string `json:"issuer_cn"`
	ValidFrom           string `json:"valid_from"`
	ValidTo             string `json:"valid_to"`
	CSCBListed          bool   `json:"cscb_listed"`
	CSCBReason          string `json:"cscb_reason"`
}

// MalwareSample represents a single malware sample entry
type MalwareSample struct {
	SHA256Hash      string            `json:"sha256_hash"`
	SHA3_384Hash    string            `json:"sha3_384_hash"`
	SHA1Hash        string            `json:"sha1_hash"`
	MD5Hash         string            `json:"md5_hash"`
	FirstSeen       string            `json:"first_seen"`
	LastSeen        string            `json:"last_seen"`
	FileName        string            `json:"file_name"`
	FileSize        int64             `json:"file_size"`
	FileTypeMIME    string            `json:"file_type_mime"`
	FileType        string            `json:"file_type"`
	FileFormat      string            `json:"file_format,omitempty"`
	FileArch        string            `json:"file_arch,omitempty"`
	Reporter        string            `json:"reporter"`
	OriginCountry   string            `json:"origin_country"`
	Anonymous       int               `json:"anonymous"`
	Signature       string            `json:"signature"`
	Imphash         string            `json:"imphash"`
	TLSH            string            `json:"tlsh"`
	Telfhash        string            `json:"telfhash"`
	Gimphash        string            `json:"gimphash"`
	SSDeep          string            `json:"ssdeep"`
	Magika          string            `json:"magika,omitempty"`
	TrID            []string          `json:"trid,omitempty"`
	DhashIcon       string            `json:"dhash_icon"`
	ArchivePW       string            `json:"archive_pw,omitempty"`
	Tags            []string          `json:"tags"`
	CodeSign        []CodeSign        `json:"code_sign"`
	DeliveryMethod  string            `json:"delivery_method"`
	Intelligence    *Intelligence     `json:"intelligence,omitempty"`
	FileInformation []FileInformation `json:"file_information,omitempty"`
	OLEInformation  json.RawMessage   `json:"ole_information,omitempty"`
	YaraRules       []YaraRule        `json:"yara_rules"`
	VendorIntel     VendorIntel       `json:"vendor_intel"`
	Comment         string            `json:"comment,omitempty"`
	Comments        []Comment         `json:"comments"`
}

// Intelligence represents aggregate intelligence about a sample.
// Downloads and Uploads are strings: the API returns them quoted (e.g. "166").
type Intelligence struct {
	ClamAV    []string        `json:"clamav,omitempty"`
	Downloads string          `json:"downloads,omitempty"`
	Uploads   string          `json:"uploads,omitempty"`
	Mail      json.RawMessage `json:"mail,omitempty"`
}

// FileInformation represents a contextual link or note about a sample
type FileInformation struct {
	Context string `json:"context"`
	Value   string `json:"value"`
}

// CodeSign represents code signing certificate information
type CodeSign struct {
	SubjectCN           string `json:"subject_cn"`
	IssuerCN            string `json:"issuer_cn"`
	Algorithm           string `json:"algorithm"`
	ValidFrom           string `json:"valid_from"`
	ValidTo             string `json:"valid_to"`
	SerialNumber        string `json:"serial_number"`
	Thumbprint          string `json:"thumbprint,omitempty"`
	ThumbprintAlgorithm string `json:"thumbprint_algorithm,omitempty"`
	CSCBListed          bool   `json:"cscb_listed"`
	CSCBReason          string `json:"cscb_reason,omitempty"`
}

// YaraRule represents a YARA rule match
type YaraRule struct {
	RuleName    string `json:"rule_name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Reference   string `json:"reference"`
}

// VendorIntel represents intelligence from various vendors
type VendorIntel struct {
	AnyRun        interface{} `json:"ANY.RUN"`
	Cape          interface{} `json:"CAPE"`
	CertPLMWDB    interface{} `json:"CERT-PL_MWDB"`
	VxCube        interface{} `json:"vxCube"`
	DocGuard      interface{} `json:"DocGuard"`
	FileScanIO    interface{} `json:"FileScan-IO"`
	InQuest       interface{} `json:"InQuest Labs"`
	Intezer       interface{} `json:"Intezer"`
	ReversingLabs interface{} `json:"ReversingLabs"`
	SpamhausHBL   interface{} `json:"Spamhaus_HBL"`
	Triage        interface{} `json:"Triage"`
	UnpacMe       interface{} `json:"UnpacMe"`
	VMRay         interface{} `json:"VMRay"`
	YoroiYomi     interface{} `json:"YOROI_YOMI"`
}

// Comment represents a user comment on a sample
type Comment struct {
	ID            string `json:"id"`
	DateAdded     string `json:"date_added"`
	TwitterHandle string `json:"twitter_handle"`
	DisplayName   string `json:"display_name"`
	Comment       string `json:"comment"`
}

// ParseAPIResponse parses the raw JSON response into an APIResponse struct
func ParseAPIResponse(data []byte) (*APIResponse, error) {
	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
