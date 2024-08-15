package constants

import c "github.com/Newt6611/apollo/constants"

type stableConfig struct {
	OrderAddress   string
	PoolAddress    string
	NftAsset       string
	LpAsset        string
	Assets         []string
	Multiples      []uint64
	Fee            uint64
	AdminFee       uint64
	FeeDenominator uint64
}

var StableConfig = map[c.Network][]stableConfig{
	c.MAINNET: {
		{
			OrderAddress: "addr1w9xy6edqv9hkptwzewns75ehq53nk8t73je7np5vmj3emps698n9g",
			PoolAddress:  "addr1wy7kkcpuf39tusnnyga5t2zcul65dwx9yqzg7sep3cjscesx2q5m5",
			NftAsset:     "5d4b6afd3344adcf37ccef5558bb87f522874578c32f17160512e398444a45442d695553442d534c50",
			LpAsset:      "2c07095028169d7ab4376611abef750623c8f955597a38cd15248640444a45442d695553442d534c50",
			Assets: []string{
				"8db269c3ec630e06ae29f74bc39edd1f87c819f1056206e879a1cd61446a65644d6963726f555344",
				"f66d78b4a3cb3d37afa0ec36461e51ecbde00f26c8f0a68f94b6988069555344",
			},
			Multiples:      []uint64{1, 1},
			Fee:            1000000,
			AdminFee:       5000000000,
			FeeDenominator: 10000000000,
		},
		{
			OrderAddress: "addr1w93d8cuht3hvqt2qqfjqgyek3gk5d6ss2j93e5sh505m0ng8cmze2",
			PoolAddress:  "addr1wx8d45xlfrlxd7tctve8xgdtk59j849n00zz2pgyvv47t8sxa6t53",
			NftAsset:     "d97fa91daaf63559a253970365fb219dc4364c028e5fe0606cdbfff9555344432d444a45442d534c50",
			LpAsset:      "ac49e0969d76ed5aa9e9861a77be65f4fc29e9a979dc4c37a99eb8f4555344432d444a45442d534c50",
			Assets: []string{
				"25c5de5f5b286073c593edfd77b48abc7a48e5a4f3d4cd9d428ff93555534443",
				"8db269c3ec630e06ae29f74bc39edd1f87c819f1056206e879a1cd61446a65644d6963726f555344",
			},
			Multiples:      []uint64{1, 100},
			Fee:            1000000,
			AdminFee:       5000000000,
			FeeDenominator: 10000000000,
		},
		{
			OrderAddress: "addr1wxtv9k2lcum5pmcc4wu44a5tufulszahz84knff87wcawycez9lug",
			PoolAddress:  "addr1w9520fyp6g3pjwd0ymfy4v2xka54ek6ulv4h8vce54zfyfcm2m0sm",
			NftAsset:     "96402c6f5e7a04f16b4d6f500ab039ff5eac5d0226d4f88bf5523ce85553444d2d695553442d534c50",
			LpAsset:      "31f92531ac9f1af3079701fab7c66ce997eb07988277ee5b9d6403015553444d2d695553442d534c50",
			Assets: []string{
				"c48cbb3d5e57ed56e276bc45f99ab39abe94e6cd7ac39fb402da47ad0014df105553444d",
				"f66d78b4a3cb3d37afa0ec36461e51ecbde00f26c8f0a68f94b6988069555344",
			},
			Multiples:      []uint64{1, 1},
			Fee:            1000000,
			AdminFee:       5000000000,
			FeeDenominator: 10000000000,
		},
		{
			OrderAddress: "addr1wxr9ppdymqgw6g0hvaaa7wc6j0smwh730ujx6lczgdynehsguav8d",
			PoolAddress:  "addr1wxxdvtj6y4fut4tmu796qpvy2xujtd836yg69ahat3e6jjcelrf94",
			NftAsset:     "07b0869ed7488657e24ac9b27b3f0fb4f76757f444197b2a38a15c3c444a45442d5553444d2d534c50",
			LpAsset:      "5b042cf53c0b2ce4f30a9e743b4871ad8c6dcdf1d845133395f55a8e444a45442d5553444d2d534c50",
			Assets: []string{
				"8db269c3ec630e06ae29f74bc39edd1f87c819f1056206e879a1cd61446a65644d6963726f555344",
				"c48cbb3d5e57ed56e276bc45f99ab39abe94e6cd7ac39fb402da47ad0014df105553444d",
			},
			Multiples:      []uint64{1, 1},
			Fee:            1000000,
			AdminFee:       5000000000,
			FeeDenominator: 10000000000,
		},
	},

	c.TESTNET: {
		{
			OrderAddress: "addr_test1zq8spknltt6yyz2505rhc5lqw89afc4anhu4u0347n5dz8urajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upqa63kst",
			PoolAddress:  "addr_test1zr3hs60rn9x49ahuduuzmnlhnema0jsl4d3ujrf3cmurhmvrajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upqcgz9yc",
			NftAsset:     "06fe1ba957728130154154d5e5b25a7b533ebe6c4516356c0aa69355646a65642d697573642d76312e342d6c70",
			LpAsset:      "d16339238c9e1fb4d034b6a48facb2f97794a9cdb7bc049dd7c49f54646a65642d697573642d76312e342d6c70",
			Assets: []string{
				"e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed7274444a4544",
				"e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed727469555344",
			},
			Multiples:      []uint64{1, 1},
			Fee:            1000000,
			AdminFee:       5000000000,
			FeeDenominator: 10000000000,
		},
		{
			OrderAddress: "addr_test1zp3mf7r63u8km2d69kh6v2axlvl04yunmmj67vprljuht4urajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upqhelj6n",
			PoolAddress:  "addr_test1zzc8ar93kgntz3lv95uauhe29kj4yj84mxhg5v9dqj4k7p5rajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upqujv25l",
			NftAsset:     "06fe1ba957728130154154d5e5b25a7b533ebe6c4516356c0aa69355757364632d757364742d76312e342d6c70",
			LpAsset:      "8db03e0cc042a5f82434123a0509f590210996f1c7410c94f913ac48757364632d757364742d76312e342d6c70",
			Assets: []string{
				"e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed727455534443",
				"e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed727455534454",
			},
			Multiples:      []uint64{1, 1},
			Fee:            1000000,
			AdminFee:       5000000000,
			FeeDenominator: 10000000000,
		},
		{
			OrderAddress: "addr_test1zqpmw0kkgm6fp9x0asq5vwuaccweeqdv3edhwckqr2gnvzurajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upq9z8vxj",
			PoolAddress:  "addr_test1zqh2uv0wvrtt579e92q35ktkzcj3lj3nzdm3xjpsdack3q5rajt8r8wqtygrfduwgukk73m5gcnplmztc5tl5ngy0upqud27a8",
			NftAsset:     "06fe1ba957728130154154d5e5b25a7b533ebe6c4516356c0aa69355646a65642d697573642d6461692d76312e342d6c70",
			LpAsset:      "492fd7252d5914c9f5acb7eeb6b905b3a65b9a952c2300de34eb86c5646a65642d697573642d6461692d76312e342d6c70",
			Assets: []string{
				"e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed7274444a4544",
				"e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed727469555344",
				"e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed7274444149",
			},
			Multiples:      []uint64{1, 1, 1},
			Fee:            1000000,
			AdminFee:       5000000000,
			FeeDenominator: 10000000000,
		},
	},
}
