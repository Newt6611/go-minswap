package constants

type v2config struct {
	FactoryAsset                  string
	PoolAuthenAsset               string
	GlobalSettingAsset            string
	LpPolicyId                    string
	GlobalSettingScriptHash       string
	GlobalSettingScriptHashBech32 string
	OrderScriptHash               string
	PoolScriptHash                string
	PoolScriptHashBech32          string
	PoolCreationAddress           string
	FactoryScriptHashBech32       string
	FactoryScriptHash             string
	FactoryAddress                string
	ExpiredOrderCancelAddress     string
	PoolBatchingAddress           string
	OrderEnterpriseAddress        string
}

var V2Config = map[NetworkId]v2config{
	NetworkIdMainnet: {
		FactoryAsset:                  "f5808c2c990d86da54bfc97d89cee6efa20cd8461616359478d96b4c4d5346",
		PoolAuthenAsset:               "f5808c2c990d86da54bfc97d89cee6efa20cd8461616359478d96b4c4d5350",
		GlobalSettingAsset:            "f5808c2c990d86da54bfc97d89cee6efa20cd8461616359478d96b4c4d534753",
		LpPolicyId:                    "f5808c2c990d86da54bfc97d89cee6efa20cd8461616359478d96b4c",
		GlobalSettingScriptHash:       "f5808c2c990d86da54bfc97d89cee6efa20cd8461616359478d96b4c",
		GlobalSettingScriptHashBech32: "script17kqgctyepkrd549le97cnnhxa73qekzxzctrt9rcm945c880puk",
		OrderScriptHash:               "c3e28c36c3447315ba5a56f33da6a6ddc1770a876a8d9f0cb3a97c4c",
		PoolScriptHash:                "ea07b733d932129c378af627436e7cbc2ef0bf96e0036bb51b3bde6b",
		PoolScriptHashBech32:          "script1agrmwv7exgffcdu27cn5xmnuhsh0p0ukuqpkhdgm800xksw7e2w",
		PoolCreationAddress:           "addr1z84q0denmyep98ph3tmzwsmw0j7zau9ljmsqx6a4rvaau66j2c79gy9l76sdg0xwhd7r0c0kna0tycz4y5s6mlenh8pq777e2a",
		FactoryScriptHash:             "7bc5fbd41a95f561be84369631e0e35895efb0b73e0a7480bb9ed730",
		FactoryScriptHashBech32:       "script100zlh4q6jh6kr05yx6trrc8rtz27lv9h8c98fq9mnmtnqfa47eg",
		FactoryAddress:                "addr1z9aut775r22l2cd7ssmfvv0qudvftmaskulq5ayqhw0dwvzj2c79gy9l76sdg0xwhd7r0c0kna0tycz4y5s6mlenh8pqgjw6pl",
		ExpiredOrderCancelAddress:     "stake178ytpnrpxax5p8leepgjx9cq8ecedgly6jz4xwvvv4kvzfq9s6295",
		PoolBatchingAddress:           "stake17y02a946720zw6pw50upt2arvxsvvpvaghjtl054h0f0gjsfyjz59",
		OrderEnterpriseAddress:        "addr1w8p79rpkcdz8x9d6tft0x0dx5mwuzac2sa4gm8cvkw5hcnqst2ctf",
	},
	NetworkIdTestnet: {
		FactoryAsset:                  "d6aae2059baee188f74917493cf7637e679cd219bdfbbf4dcbeb1d0b4d5346",
		PoolAuthenAsset:               "d6aae2059baee188f74917493cf7637e679cd219bdfbbf4dcbeb1d0b4d5350",
		GlobalSettingAsset:            "d6aae2059baee188f74917493cf7637e679cd219bdfbbf4dcbeb1d0b4d534753",
		LpPolicyId:                    "d6aae2059baee188f74917493cf7637e679cd219bdfbbf4dcbeb1d0b",
		GlobalSettingScriptHash:       "d6aae2059baee188f74917493cf7637e679cd219bdfbbf4dcbeb1d0b",
		GlobalSettingScriptHashBech32: "script1664wypvm4msc3a6fzayneamr0enee5sehham7nwtavwsk2s2vg9",
		OrderScriptHash:               "da9525463841173ad1230b1d5a1b5d0a3116bbdeb4412327148a1b7a",
		PoolScriptHash:                "d6ba9b7509eac866288ff5072d2a18205ac56f744bc82dcd808cb8fe",
		PoolScriptHashBech32:          "script166afkagfatyxv2y075rj62scypdv2mm5f0yzmnvq3ju0uqqmszv",
		PoolCreationAddress:           "addr_test1zrtt4xm4p84vse3g3l6swtf2rqs943t0w39ustwdszxt3l5rajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upqhns793",
		FactoryScriptHashBech32:       "6e23fe172b5b50e2ad59aded9ee8d488f74c7f4686f91b032220adad",
		FactoryScriptHash:             "script1dc3lu9ettdgw9t2e4hkea6x53rm5cl6xsmu3kqezyzk66vpljxc",
		FactoryAddress:                "addr_test1zphz8lsh9dd4pc4dtxk7m8hg6jy0wnrlg6r0jxcrygs2mtvrajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upqjgg24z",
		ExpiredOrderCancelAddress:     "stake_test17rytpnrpxax5p8leepgjx9cq8ecedgly6jz4xwvvv4kvzfqz6sgpf",
		PoolBatchingAddress:           "stake_test17rann6nth9675m0y5tz32u3rfhzcfjymanxqnfyexsufu5glcajhf",
		OrderEnterpriseAddress:        "addr_test1wrdf2f2x8pq3wwk3yv936ksmt59rz94mm66yzge8zj9pk7s0kjph3",
	},
}