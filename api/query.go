package api

import "fmt"

// queryAPI helper to make API requests and print results
func queryAPI(queryType, queryKey, queryValue string, limit int, apiKey string) error {
	data := map[string]string{
		"query": queryType,
	}

	if queryKey != "" && queryValue != "" {
		data[queryKey] = queryValue
	}

	if limit > 0 {
		data["limit"] = fmt.Sprintf("%d", limit)
	}

	response, error := MakeRequest(data, nil, apiKey)
	if error != nil {
		return fmt.Errorf("Error retrieving samples for %s %s: %v", queryType, queryValue, error)
	}

	fmt.Printf("Samples found searching for %s: '%s':\n", queryKey, queryValue)
	dataJSON, err := PrintData(response)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(dataJSON)
	}
	return nil
}

// QueryByHash retrieves malware samples by hash
func QueryByHash(hash string, limit int, apiKey string) error {
	return queryAPI("get_info", "hash", hash, limit, apiKey)
}

// QueryByTag retrieves malware samples with a specific tag
func QueryByTag(tag string, limit int, apiKey string) error {
	return queryAPI("get_taginfo", "tag", tag, limit, apiKey)
}

// QueryBySignature retrieves malware samples with a specific signature
func QueryBySignature(signature string, limit int, apiKey string) error {
	return queryAPI("get_siginfo", "signature", signature, limit, apiKey)
}

// QueryByFileType retrieves malware samples of a specific file type
func QueryByFileType(filetype string, limit int, apiKey string) error {
	return queryAPI("get_file_type", "file_type", filetype, limit, apiKey)
}

// QueryByClamAV retrieves samples by ClamAV signature
func QueryByClamAV(clamav string, limit int, apiKey string) error {
	return queryAPI("get_clamavinfo", "clamav", clamav, limit, apiKey)
}

// QueryByImpHash retrieves samples by Imphash signature
func QueryByImpHash(imphash string, limit int, apiKey string) error {
	return queryAPI("get_imphash", "imphash", imphash, limit, apiKey)
}

// QueryByTLSH retrieves samples by TLSH signature
func QueryByTLSH(tlshash string, limit int, apiKey string) error {
	return queryAPI("get_tlsh", "tlsh", tlshash, limit, apiKey)
}

// QueryByTelfHash retrieves samples by Telfhash signature
func QueryByTelfHash(telfhash string, limit int, apiKey string) error {
	return queryAPI("get_telfhash", "telfhash", telfhash, limit, apiKey)
}

// QueryByDHash retrieves samples by DHash signature
func QueryByDHash(dhash string, limit int, apiKey string) error {
	return queryAPI("get_dhash_icon", "dhash_icon", dhash, limit, apiKey)
}

// QueryByGimphash retrieves samples by Gimphash signature
func QueryByGimphash(dhash string, limit int, apiKey string) error {
	return queryAPI("get_gimphash", "gimphash", dhash, limit, apiKey)
}

// QueryByYara retrieves samples by Yara Rule
func QueryByYara(yara string, limit int, apiKey string) error {
	return queryAPI("get_yarainfo", "yara_rule", yara, limit, apiKey)
}

// QueryByIssuerCN retrieves samples by Issuer Common Name
func QueryByIssuerCN(issuerCN string, limit int, apiKey string) error {
	return queryAPI("get_issuerinfo", "issuer_cn", issuerCN, limit, apiKey)
}

// QueryBySubjectCN retrieves samples by Subject Common Name
func QueryBySubjectCN(subjectCN string, limit int, apiKey string) error {
	return queryAPI("get_subjectinfo", "subject_cn", subjectCN, limit, apiKey)
}

// QueryBySerialNumber retrieves samples by Serial Number
func QueryBySerialNumber(serialNumber string, limit int, apiKey string) error {
	return queryAPI("get_certificate", "serial_number", serialNumber, limit, apiKey)
}
