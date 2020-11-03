package menandmice

type NextFreeAddressRespons struct {
	Result struct {
		Address string `json:"address"`
	} `json:"result"`
}

func (c Mmclient) NextFreeAddress(addressRange, startaddress string, ping, excludeDHCP bool, temporaryClaimTime int) (string, error) {
	var re NextFreeAddressRespons
	query := map[string]interface{}{
		"ping":               ping,
		"excludeDHCP":        excludeDHCP,
		"temporaryClaimTime": temporaryClaimTime,
	}
	if startaddress != "" {
		query["startaddress"] = startaddress
	}
	err := c.Get(&re, "Ranges/"+addressRange+"/NextFreeAddress", query, nil)
	return re.Result.Address, err
}
