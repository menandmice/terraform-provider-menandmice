package menandmice

type NextFreeAddressRespons struct {
	Result struct {
		Address string `json:"address"`
	} `json:"result"`
}

func (c Mmclient) NextFreeAddress(addressRange, startAddress string, ping, excludeDHCP bool, temporaryClaimTime int) (string, error) {
	var re NextFreeAddressRespons
	query := map[string]interface{}{
		"ping":               ping,
		"excludeDHCP":        excludeDHCP,
		"temporaryClaimTime": temporaryClaimTime,
	}
	if startAddress != "" {
		query["startAddress"] = startAddress
	}
	err := c.Get(&re, "Ranges/"+addressRange+"/NextFreeAddress", query, nil)
	return re.Result.Address, err
}
