package assets

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"sort"
)

// AssetNames returns a list of all assets
func AssetNames() []string {
	an := make([]string, len(data))
	i := 0
	for k := range data {
		an[i] = k
		i++
	}
	sort.Strings(an)
	return an
}

// Get returns an asset by name
func Get(an string) ([]byte, bool) {
	if d, ok := data[an]; ok {
		return d, true
	}
	return nil, false
}

// MustGet returns an asset by name or explodes
func MustGet(an string) []byte {
	if r, ok := Get(an); ok {
		return r
	}
	panic(errors.New("could not find asset: " + an))
}

func decompress(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	r, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.Bytes()
}

var data = make(map[string][]byte, 60)

func init() {
	data["core/01_branching.lisp"] = decompress("H4sIAAAAAAAA/8ySzU7rMBCF936Koy56nUV1220iBO+BkHCdCbVkPGA7IITosyOnaX6cdgMbuqiiM/GZfF9SVVWF8NLGSB6aPZXYe+X0wbgnIWRNzbPSnvF+ICeAz7JmvVEhUMQqZasvAdxHChFr7Ln+eBDAozQNjl0oa8bxLg0KOGOLIi/dOI4/KXbGTsvnxZpdnZem7FS6hraqDRRSpUzLIAO93p7jQgBAWiVvsIO05IbRaQbIPsB2SKwJEf96l6Vp+nx273i9KyZzdmE8mR4U0nf2uv9hd/ebcqolprpM2RtBQtpmSIi+pfMw512QlmQDXWa2FEeoN9KR/TD8rxLV2DXBvyoOs9OTODOmBmH5a8o7coHsc3/s/6Y+9r+wNy+6LjItWX5y4jsAAP//2G3rUyUEAAA=")
	data["core/02_functions.lisp"] = decompress("H4sIAAAAAAAA/3zRsWrDMBAG4F1P8ZOhlYZAbG82tH2PUqgqnSDUOhWd3KU0zx7iOInA4O047vuH+4dhGCA/UymU4VKmHmFiV46JRSntKUTrckJgBfz1Prm9FaGCXeDdvwLenxBSjvKhAD0epeB5ieuj/aa9G5NMmRQAaJfYzxOgLXvoF+iR+JpgcDDQv+RKyq/Q8w4HY8xj3rDN2jaVbbZsu7ZtZdst261tV9nuZnsahW6nxtz/IY9/jTZ+ebsEG1O/39O6gMvuWgHbSKh6+LxAnOa1DrxMp7d78jkAAP//oLEuavYBAAA=")
	data["core/03_concurrency.lisp"] = decompress("H4sIAAAAAAAA/4SRwW7DIAyG7zyF1UkrHKrek8veo6o05DhrNII7IFOraXv2KSEkkE5aTsT++f3xu67rGvx1CIEcIDuqANni4BxZvAshG2p7jY5B+7tFAfBVNYwH7T0F2E3F3bcAOD1Dy673ZwEgTecD7GfbqtfvdEDDfnAkAADkJ2Fgp+IPsvWruOHDZBrdlFI5wxtZcjrQFiPV/yMxFHKApXHEy1aKF20tmciYvlVvZn01votyo80N6ruQtNM5l84JlJNT0H+0DIWV/fYYXQptQSjvHzF7z/5WjK88fZRsRfLtEOL2itxj9SH1VzmCnq6O+87TE8j5pM5xYFywXPsNw8/LsvBRk3pKid8AAAD//+F/cyCjAgAA")
	data["core/04_sequences.lisp"] = decompress("H4sIAAAAAAAA/7SSzUrEMBSF93mKwywkBecF2o3vMQwY7tw4hebH5FYR0WeXJp1RWosbZxfOOXx8cNN1XYccRxFOoJC4RebnkT1xVkqf2HpQ8GREAe/tKdDe5MyCXU13Hwo43MGG5PJRAfqJPScjrABAz+jWhrRnQ2ccpmWd34MHru9jWQOaXS8lbpqmuQjYfhBOS4GaVgE7esL/Wbye2UMXaHGZ86XfVdCZuLRzJt5ErQj8UPvWcIZSgITJINDS55KvDvaoTYzDG0oLXc+Kz4cyWMOHPq++whxvoqfyb/ILk4TVla/FJr3Wv/G/AgAA//8hnKlQ3QIAAA==")
	data["docstring/and.md"] = decompress("H4sIAAAAAAAA/0SOQa7CMAwF9z3Fk/4inwp6BxasOYPVOk2kNK5sl8LtUSgSO+tpPJo//FOdEEWXvj9hZW2ngWBJ1C9j1nHLnuuMQHUK3e1BZSNngyfGAUeVBYWjwwWa5+QDrgYTqSCDVAb/3gSEEKkYvwLazGfsuRQo+6YVnsiPfcDdE+uejZH9gFaVkXlqmq+ztbWWyk//BA3dOwAA///2Hivp1gAAAA==")
	data["docstring/apply.md"] = decompress("H4sIAAAAAAAA/1TPTUrAMBAF4P2c4kEXJggFFdy76A3cl7GZ2kCaxPzUentJgwWXM3zvDTNAcYzuB2v1C7J8abTZSganz7qLLxklgC9QbPA0HewqF8komyCmcFgjpmWr+EXA3twdtmQ0LlfJP//XBxVS9wt/ONHgHrvPj0TDgDeP6eQ9OiEClJEVJx7UE57xonVb9T8ecWqi981mSPf4ts4hSanJY36dR/oNAAD//8CDF8f1AAAA")
	data["docstring/assoc.md"] = decompress("H4sIAAAAAAAA/2zQz0oDQQwG8Ps8xcf2YosWwZuI4EF8AcFzOpvpBmczS5Kt7duL3eIf8Jjw4/tIVrgi95bx8M4nHKjO/LjZrJGNKdhBUP7AmQiFHBgeNueYjdOb1HqB/7gdOfdoihgYk7WD9NzjuwUTifk1msE4ZlOQgscpTn9SpEAbSrPRQfYTtMXrII4yaw5pCvEv10spbKyBGGhp/h1WJdiowk8adAQfM09nGpBAJsWOEeeDeoiCUMQ8bnIldxTyQZpuU1qt8KR4PtI4VU4Jlx8mALhXGhmdT3MEW7fsaM/A7fZumSrpHt1L69bpMwAA///9TMgsfwEAAA==")
	data["docstring/async.md"] = decompress("H4sIAAAAAAAA/zyPS07EMBBE9z5FkWzsWcwdZsENuEBPXMEWjjvyZ3BujxIBy269V6WaYaUeecGqZbvdHK4rFM3aazrAl6QujRWCZ9LlC+YjEHvRV/T0l1bxHVPCk/+0R8wQVO5SpBEtFIqHruDg0lvUfMcjHyisPbWYP3GaPIkW+Nv0l+pjXaR4+rsx84xHxvuQbU80BrCeKwbsXnSLlc4Z4HxfOwwA2IEpMCWdnLuM2grscJhwaH+bnPkJAAD//43Lk9AFAQAA")
	data["docstring/channel.md"] = decompress("H4sIAAAAAAAA/1STwY7jLBCE7zxFaeZiRzN5gNxGv+b031Z7W63kDm7HaDFkoEkm+/Srxk4c+2QVRXfXB7yisSOFwL6FTUzCGYQSXO8SW3ExkMfiMB/3Pzh19SSELKlYKYkhI4kulMw9JOLEgRMJg+Dp7w2ZvwoHy4gDLuQL5z1+jozEuXiZS46Ux/eJzrAxZJfFhZPaKaDjyUmHoYQ61RsInfUx80YLPejRaI8fLMnxRatQAHueOAiGFCfIyOtEE92wO/po/+zecCVX+w4xVVfgb5kH1lRHhk4i3COmh+fOZTbUufo9PtVYm6/76eG9Ou9BPkfUzihBnK/VjmUYOCmRwZc8cg9SOAuoOMCJkq1wLpw2WfbGvL7iv6XH/3zLxuwOmb920G8Te572fk692R00WjUqrhozPfBqY90+xDSha3R5jtV2ZneomXUrrRskzig2iLqmam1XB/0I+Pym6ezZGKDxLPhlx/VSGiwfT2jqeLDjqlqPZm6t8u+HXn8ayrdgF63hCS9DjC/ts3CktArWt62Zd0p8v7CVmNAoOy2uawvWjMTnxFkvk0bz8fru+cIeJfSczi4EF075Dqy7E36+qtfR2VGP2AVMMQssZc7Vf048cErc40q3zUvavKPl+fCMD3SMF4aNxfd6qomvSc8vHCrYp0D3amah4OQZzCJUMG1r/gUAAP//p+iBAyAEAAA=")
	data["docstring/concat.md"] = decompress("H4sIAAAAAAAA/1SOzWrDMBCE7/sUAz7UplDo772UvkFuIRghbyKBIjlayY7z9EEymOQ6336z06DVwWuVIHx57eDUzboFa8ZeJZZCMnvNQn+Ra6LK3bIBzCYIFymxT7CCZBiRJbuEcHxos/5U2aRcZimMlTYYY5jswMPW+EbUNPj1+L+q8+iYCNvQ/Ts+8HnAS/uFb/x0HdHOWMFsnUPklKOvP54n9m3VsDo93QMAAP//uAS+0/kAAAA=")
	data["docstring/cond.md"] = decompress("H4sIAAAAAAAA/7SSMW/bMBCFd/6KB3uIFQSGMmRIEKTo0AKd2z0+U2eTCH0UyFMl/fvipMQtCnSsJlF47+m+e9xi57N0eO4Ld9DA8nJ7C06VPzXouZxyuVSYJGrMQgnHQuJDlLP7mguYfMDhD/cBPtFQ+c6OsO/RkzLGmBKODP5JaSDl7g4kHeIJURErtAwaZuwkK06ULMBeJaZmSbqx9BvYPCa/xiwphXUoYplZA5cxVl5MwpO+z2OmvmTPtXK3d267xWfBl4kufWLngF3HJ0x4fGyWkyE7ANg9Y8JD2wCbkSoS1woNJHhoN6vgBRPu27ZZBefCpFxWzX37Lvr9LKIoOLKOzLJpnPsm0BArPBn269+K1+vyPjj3+J7iOWiaccmFQVw1sEZPKc3oE1ONcsaYh9SZTzNsBSR4smrxxvOYS/f0P0jXP/yD9EcwCKrZiElX7DGXt2r9HNkvVVl1hyXn8DHqtfBqMOtl2btfAQAA//+f2055vgIAAA==")
	data["docstring/conj.md"] = decompress("H4sIAAAAAAAA/2SQvW7jMBCEez7FAG7OuIOAy3/rIm8QIEVgWCtyZTGmdhVyZcV5+kCR7SYt55vhzqzwx6u8o/AHWs393zUohAJO3LNYgSloVkcWz27zW5vtGoWaxFeuwksXCxru6Bg1Y4opIcS25YzAA0uIsocKrOM5wGc2hp0GrrBBisUWy5B/4H8gHNnbJYmG5XXqYmKQgEpRH8nikdHTgQtEsR8pkxhzATU6GjQHzlH25+vaUbxFlSVT1DBpPmCK1iHR1+napqCMvgMVqHDBkDWMngOaE+qehhqaUbcxGee6cm61wkbw/En9kNg5nBd++48b3OJui3s84BFPa+de558z25iXKc4ldxf2gm53lfsOAAD//0Od9nWqAQAA")
	data["docstring/cons.md"] = decompress("H4sIAAAAAAAA/3ySP2/cMAzFd3+KB2SoLQQHJP2zd+jWsUDHmJF4NgEd5ZLy/fn2hexchqLtyEc+8kdKD+hjUceZ8spw/jUgltOrKDvoTb1InUEtubJG7n5usaJYYuP0nniEr3EGNWcWryiGM8da7BF1Zhj7miuk5ZUvf9TsY1rdPnUxXlgTJ9SyycVkEqV8wB1g1f8gKMi9RKEqZwaZ0W3DMG4EWjCtZKSVG9hElkQnLMWlStFD1/2YGUonxtgONO7YZ/YqE6MccZlZ8V18gZyWzCfWSs3qCM0Qqq2xcnrf1NuqkXN2vN6wkFgbSBgj2Yg+xFC0sra647ZvoJCSsTsWstrUYIEn8co23B9ljOkf7hQ42ob1N/+h6x4e8FXx7UoNv+uAPvERV3zoP+ITPuPLMNzF29svecZ1F7foCbcBaJcSB+99cJGcYVxXU7z0T3jGvdvLofsdAAD//yLnvbNuAgAA")
	data["docstring/def.md"] = decompress("H4sIAAAAAAAA/2SSzYrbQBCE73qKAh/WXjaGbJ5gDzkEcgzkEALbkkqegZmR0t1jW3n6MLLxHnJpRP98Vd2aHfYjJxTJxDRrPkDM4qkYZEvaIgNxFo3SJ3Y/Y0r3DsgjjVjggRiqKot/DL7gEuIQEA1P1ahP6FeMnKQmP+ItpQfCIErEnKtvRCkjlFaTN7gUUHVW9IzlBJVoHBGnVhB35sWbRpaR8BnKT3ePHpiP+BGioWeQc5y1NY5xmrhZ9SAFeTbH92iLvUCsDa24tE1PLFRJacUkMcFiYvG0NktWh4BBjHbsut0ObwVfr5KXxK7D7ajXDgD2WZbtA9hPBb/W31jx/IzXwz1r/PO5hdcWvhwOXbf55Y12MzIoxQlBkr8rsiy4BCpBGQLOkioxT9s/8KAkFp3PceTYkJVloG1rz7VPHFGXuUBbxfyIb37T8MCCx93EP7R83tD/vwe8X9+P3b8AAAD//70yXNVEAgAA")
	data["docstring/defmacro.md"] = decompress("H4sIAAAAAAAA/1SQT2/CMAzF7/0UTyCtgU1oXCdNiMM+ww6Ig0ncNVKadLE72LefEv4Meqhe/LPfszyHcdwNZHNCpIExsNIGu6KXyz26lIfnBUjEf0UBITM5zqgTzacPAY5toMygcxE+QnuGnXLmqNVVRrL8gmPvbQ8vaCfh3OLwC8cdTUFX2F6nBXwaKTp2hRenS6SPSLkITaCgRfQMSVO2DJsctwJHWlYcMwtHJfUp4sBdygyv1fuHwkTKbtU08zm2ER8nGsbATYO7U9gUXQMAM1/gwFGlxrWFtPUss9qwe4INNAnLvr7NsecII/y9uYJFBYDxHcw71jCB4w1eKWAuJbze1YIXRSvjpMr5zXc38tD/r9eLh44U5bwzTGbRy/+WXb/mLwAA//982CfHBQIAAA==")
	data["docstring/defn.md"] = decompress("H4sIAAAAAAAA/0yOzU7DMBCE736KUXNo0gJqOFZCiAPPwKHqYeOsyUrOpvhHgBDvjmyhNr54vrFndxq0IzuF0syYOdEzTkXvdme4Jcz7DhSjvGsE1U/xQpbhstoki5o38R4jW0+BQVcfokgTw+YQWNMteYfPSewEidjmyGGL4RsjO8o+PRjTNHhRvH7RfPFsDP7bORkMAGwCpxxKFwV7nsvoxdVNToZFyVpB5I/ManlTIz/HROKRQubfapzkXO/WLjpWBbRPEBw6HNbcd+jX/HjjI/vIxd+jdTKgva/v3Yr6rhzzFwAA///wRPyAYAEAAA==")
	data["docstring/do.md"] = decompress("H4sIAAAAAAAA/zzOParDMBDE8V6n+IMbSTweSZkmkCJnSC2Sdbwgy0YfQccPdsDV7vCbYgbsa2Fc8uy9Qz4htlClMLdYdY2yUzEPjfFQJDynHdBEbTn9kWW7mt7USRg1haOvSyIUtBaylBbrvzHDwC1x72FeoxgDVkfslc7FbYsMgF1zorvf7z2d82lLzplvAAAA//+liQqetwAAAA==")
	data["docstring/drop.md"] = decompress("H4sIAAAAAAAA/1SQT0vEMBTE7/kUA3vYBqWg+0evHjx4FzyILDGZ0kKa1ORVWz+9ZLvq7vXN8Jt5s0LlUhxg4xgEmR9XGuWQIS3RdCkL6NkzSEZsimNksMzqpfMeiTKmAANvvuc/EdIaASfrR8dz0voYs74gFnVI8bNzdP94VDYGa4TBCJ2u8bQ4cxyT5ZnPJCK3MQlTiQ1HVypyFrrlr2uYAPaDzPBdFnyV6u88taerlVqt8BDwOJl+8FQKqBwbTFhXN7jFBlutf48zXnfY4w73b8ut7LfBhFkr9dx2GVwwS85polLrcqRDtcUOe32o1U8AAAD//6EiwBaGAQAA")
	data["docstring/eq.md"] = decompress("H4sIAAAAAAAA/3SQsW4jMQxEe33FnF3Yxhn7D1dcESBAypQGrR1ZQmTKFrlx8vfBeh0gTRoSGIAzj7PGllek1s/38XcHp7mhJAiMjpbwLnWiQTpRRqqXKBXe4JlIpZuH11IrOn3qikOSajxADNaazrs4qLFN6uwGWeI8i6MYtPkvtgNePLPfihG3nwneJx6GENZr/FP8/5DzpTIEANiOTMhYZdbaVrtF4/VbQN6F8FzeiHMzx6VzLFGctofnYkiTRi9NEUUh1RqOhPIkzhHHz/ngQh2Lnu6Ymz8bxCxdorMPeNLFJYpx/8BcyI98wHOcq51/llrvRSzF8jpJHcJXAAAA//+7FebRjwEAAA==")
	data["docstring/eval.md"] = decompress("H4sIAAAAAAAA/1JW0EgtS8xRSMsvytVUADFLE0tSixUSwSJcgAAAAP//DEc3fR8AAAA=")
	data["docstring/filter.md"] = decompress("H4sIAAAAAAAA/1yQzU7rMBCF936KI2UT615VasuSDUIs2LOrqsqkYzLSMA7+aROeHjkRgbKz5xx932gatJ4lU4Qv2iHRxz8LcZ8sE5Yg1WEh7SiZx0guU4KrlWkNcO1DInRBM2kGJ+SeECkVyQgebhhkYn2bx0MMFz7TeRZmDooc5uDipFCq/Zvaat/g2f8B11+Fc+dmUjXHkvsJrYYM7yTRf9Snstha1x8VriyCVwJrJ6WqWH/x677f7o0xTYMHxdPo3gchY7DerfWKw3hEe48Re2tx2GKHPe6O1piXntPiiZRLXPi3tzu1W+zsyXwFAAD//1JglZSKAQAA")
	data["docstring/first.md"] = decompress("H4sIAAAAAAAA/3yQMU+EQBSE+/0Vk1whJIZCL0YqY2Fhb3/AMisvPnZhd/Hw35sDo7nGdibzvTdzQOEkpozEuURkXqJPyAOxy1SO9BnBbWLivNBbmrdBEtzibZbgcRbVn/A/2YlWnLD/pdwiRHhRyDUdksBxyl8VXjPOYdEeHdHR04mVVpED7ED7ARciWny2Kn9YLEn8O5oicX7aijXo6EIkbKu6eduHzeV+E5lyUxlzOODZ42Vtx0lpDFD0dFhxU9Q1Ho64v8PxsSwvxt5vLY3BNgT30NUOp7o+VeY7AAD//0ew0thgAQAA")
	data["docstring/fn.md"] = decompress("H4sIAAAAAAAA/2yRwW7bQAxE7/qKQXKo7SZG7aZBb4EP/YMCPRhGTK+oaNFdUlhyXevvC0lBD22PXA4ed2buseoEQplfkNnpBcdp2GxO6LTkj2uEwuRsIAGJypi1GroqwaNK8yOm9C75vwLekyPTiAtjIDNuQUWrtIgCQheL+WNIZIZMIly2TXN/j4Pg243ykLhpgFXLHVqtl8QNAKwSO465JuxP8wNmH3eLxHClVNnucLydsNpscEOuab1er2eY6+OVg2tZWJmGdzY+rHbY4zOe8AXPs/x7Hw28/AS/JreFvZbJF2Oh4PW4xxOe8RW7T9jtT6/bpjn8k4RhKNrWwCCEpFYLL+EEHSLbAqQSaXYwb3q6Mi7MgsIdF5bALbqieRZbLXOQUd5gQQfe4pCci5DHK6fxASqMQPKnoL9aKZwpiiHo5M45jSB3Cj23cEV0Q9JAaaHjMqLadOycKF9aOiOKOVML7XDu5PwA77W+9fAptOlwmrauyJy1jNP407bN7wAAAP//M1mS8XUCAAA=")
	data["docstring/future.md"] = decompress("H4sIAAAAAAAA/1yRUZLbIBBE/zlFl/0jqxLfIR+5QGoPIAytFRU8uGDQrm6fAinOOnzOMP1mus8Y5qo1E3PK93G8gKuN1SoLLAoVae6tAls2cUtOkmqJm/lFrVnar0MgCHTZddpQq4vTkOSKt4V45LQGT3+ofYQYceMT59t8Iz5stkrokml9E+InXW0632DFw8oGZ2Ms0LQDD8wuOd5icr9HVNEQnwsVLHYlbqTApfsjUhm3f/CrMeczfgh+ftrWNQYYPOd221+HDAAMmr6vdJoyhncK2669AQy8B8Up058ur6WUrbzz/+rGGNPH6dJe581VL8a8LaGA+xrNNF8dv7jcTqfHNFedoItVuCQrs5Z+bGapUUsPQF4iQzcGQTTBYr9hT2bq5KlLH7n0r18s3Nn5GXnzjVfzJwAA//872/RsPgIAAA==")
	data["docstring/generate.md"] = decompress("H4sIAAAAAAAA/1SSzW7bMBCE73qKgX2x0sTv0ENeoAh6KYp4Ta4sotSuyh8n6tMXJCXbOYrkzjczqz0OFxYOlBiDhulbj+07ghD5b2YxDIqLmDGoaI5+6V6v5HN9k0ZGnNm4wbGtEhFO6uhMVTaNgclCB/Anm5ycyhE/OOUgXxBppIQP5z1caob0ygEkS5ktnAKtSEqgwODJpcT2iJ/t4uEMOTq5gODVkPcLotG5GMxiioNNshjG6VCmmnx/OuJtvVgVq2u2oGLX6EvQnJzwMyYmKZRmCAWEpNVCC/J09mr+PCFLcr7yjIbAcVaxZZA9TywJLiJwVH9li/NSKRLzxGGzuZV07Lr9Ht8Fr580zZ67DjhYHmDUa4j3ZXYA0GLtAttd/3iggeTCX88W9l4/dn3fV82kL1c2ScOq3Hfd2+giuHFbPIrRXaQa9PRvua8y1O2WKOXu9oNtBbVdBkdnzzg1wFr7A7c8rpha4q1DSvdZFzHkst21L/sMElv0ZfNQWU3x/Vet4pb/Fvr3+7H7HwAA//8ocJVqCQMAAA==")
	data["docstring/has-meta.md"] = decompress("H4sIAAAAAAAA/5SRsY7bMBBEe33F2ClsOYmB5AdSB0jpLgiElbgyiVBLgbu0nb8PKPnOh4Oba7Uzo3nDT9hPbPQDY8rT5xbGaoqrZ/OcYZ4x53QJjt2iUHi6MKrFkVHzcwTFeD/xhWIhY1gCiSQjY4f6kfVLzRKYD4qxyGAhCa4hRmS2kgWd5cLdESfPGENWgyT5+oi5a7VEQ5Cl2WvOGhHkjG6kqNwdmwbYOx4FLpU+Mn7f/mB/OOCG721bjyv1em2b5lSL8Y2mOfKTXuh5oKKPfyoGEvT8AD02za/wlzElNcyZXRjIVvC3zNVGUVP1Cp8Xtv5fNcwsrjJUtN1mh8FTpsE411GCYmIShXmyZfM0PnueqSzLLZJ0fWkXkqyTbFbsrecY0xbfPkCuloOcFSQOUqae87LBzt6t8D8AAP//HDug1VICAAA=")
	data["docstring/if.md"] = decompress("H4sIAAAAAAAA/3SPvU7DMBRGdz/FUTPEliIkKH8TiIGBx3Ca68aS40S2Q9O3R04HWNhs3e+c+90G7R1LkoEySkRClnfDIsnNacpkPy1B6JONp9HHs/pyNYh827DaIsPO+pMtgs+UtJbxio5zwdmQpaM+ow+m27m2bmmp8pr/1dg4kKSsKcrQMZdR0sVnuUG11Q3q8A4brx0XHwK9/KO4U6pp+Ih8brZeoBToQRwbrb7ngSOPPPHMC6/G7FPv0G/oIJHNcDQK4LDVmr0/H/5882RDOBj1EwAA//8lsX8FPQEAAA==")
	data["docstring/is-assoc.md"] = decompress("H4sIAAAAAAAA/4yRsc7UMBCEez/F/Lni58QRCehoEAUFEuV1CEUbZxNbOHbk3eQ4oXt3ZOc4AaKgtLwzO/vNAS9IJNn3GFOeXx6hLCq4OFbHGeoYS06bH3ioEwLKjCrxpH5jMZ9GUAj3X94orKQMTaC4D56KTYQ6LxjXaNWniIsPAZl1zRGd5pW7FmfHGH0WRUzxVdX+mpM1KHysgR4eu9zHCd1IQbhrjTkc8CHi43eal8DG4HHfj3eRZkbT+xCa2+O58DRdmxu+vMYbvP16NOZccvJu8GfMfQl6trQK1yzqfN7JwAsIG1tNuTXms//GmJMolsyDt6Qsp78Y2EIoSELPiDyR8oD+WgQLx6HcVVY8Pz3DOspklXOB5AUzUxSoI63s0/ivpua1ktRi/3tjbcXydOfSOA4hNf8FYO/J/AwAAP//ysZCyzUCAAA=")
	data["docstring/is-list.md"] = decompress("H4sIAAAAAAAA/3yRT2sbMRTE7/oUY/tgm7oL/V+fSg89FHr0zZTlWfvWEtFKi95bO/n2Qdo4kBByFTNP85tZYRO86C/0KQ8ftlAWFVwdq+MMdYwxp4vvuKsKAWVGcYj524NCeHrmC4WJlKEJVAW74o5Q5wX9FK36FHH1ISCzTjmi1Txx2+DgGL3PoogpfizWm0ymoPCxxng+Mbt9PKPtKQi3jTGrFX5H/LmnYQxsDG5Q680nfMYXfN3i+A3f8QM//2+NOZRMPKtfRpov4sSWJuH6sbBNcaaHFxAubDXlxph//o4xJFGMmTtvSVl2r4AtRVCQhBMj8pmUO5weimHk2BWK8sd6sYZ1lMkq59KIFwxMUaCOtNac+rfWGKZam5bzdZWm0i9m/KXjENISx/3+fex5CfMYAAD//6e3brsNAgAA")
	data["docstring/is-mapped.md"] = decompress("H4sIAAAAAAAA/4yRv87UMBDEez/FfLni48QRCehoEAUFEuV1CEWbZBJbOLZlb+44oXt3lD8cAlFQ2tqdnfnNAS8mSYn9ewwxTy+PUBYtuFqqZYZaIuV4cT37daJAMrHtmE8DxPv9nxfxsyihEbJPQG+Jp0UlQK0rGObQqYsBV+c9MnXOAY3mmU2NsyUGl4sixPBqV9gHy+wVLqyGHiLbvgsjmkF8YVMbczjgQ8DH7zIlT2PwO+CPd0Emomqd99X98Uwcx1t1x5fXeIO3X4/GnBen3BT+NLpdQctO5sLVjFqXNzRwBYILO425Nuaz+0ZMsShSZu86UZbTXxQ6CRBfIloicBRlj/a2LCSGfgm2nHh+ekZnJUunzAsmVzBRQoFa0bWCOPyrqmleWeoiv1GoVyJPv5BUlt7H6r+ybyX9DAAA///YqAyuMQIAAA==")
	data["docstring/is-nil.md"] = decompress("H4sIAAAAAAAA/3yQP4/TQBDF+/0U7y5FLuKwxJ8mFaKgQKJMh5A1scfZEeNZa3ecwLdHuw5IIHTt7syb3+/t8GSiHzClPL86wLl4wS2yR87wyFhyusrIY5sooMww0fB5AqneH/lKupIzPNXP57po8CgF02qDSzLcRBWZfc2G3vPKfYdTZEySi8OSvTbR31NlVYdYA/iTsC2LXdBPpIX7LoTdDh8Nn37QvCiHgLvN/ukN3uId3h8qzyGEU2XhbexvlC0KZx5oLbxdbEhVDVJAUCnehfBFvjPmVBxL5lEGci7P/1gOZCAtCWeG8YWcR5x/1oWFbazs9cD+YY8hUqbBOdcapGBmsgKP5K3YNP2v/XltXXmLF+2a8UNTfoysmh7x9Xj89qLw1n34FQAA//8j90f0+QEAAA==")
	data["docstring/is-seq.md"] = decompress("H4sIAAAAAAAA/4yQzU4CQRCE7/MUBRyAiCT+y8l48GDikZsxm2a3l5k4O7NM94K+vZnlJ9Fw8NzVle+rESbCmyfUMTUXUyiLCnaW1XKCWkab4tZVXPUJASWG8KbjULKY1xrk/eHEW/IdKUMjCFvyrjpFZ7krQK0T1F0o1cWAnfMeibVLAYWmjos5lpZRuySKEMPl8f0Ylc4rXOjBTjX7BhfWKGrywsXcmNEIzwEvX9S0no3BwXI8ucI1bnA7xfsd7vGAx4+pMcuMxfvwOSpj3twno4miaBNXriRlmf3RKSmAvESsGIHXpFxh9Z0fWg5V5svY48EYpaVEpXLKvk7QMAWBWtJ+zlifW77p+lE01/8eV+a94aBXHFr2Pg6xWPxH7CcAAP//82gtTf8BAAA=")
	data["docstring/is-vector.md"] = decompress("H4sIAAAAAAAA/5SRvW7jMBCEez7F2C5s43wC7v9cHa644oCU7oxAWFMrkwhFCtyVnbx9QMnOH9Kk5sxwvtkFVie2mvIftCl3n9ZQFhWcHavjDHWMPqeTb7gZFQLKjMkj5n8LCuHywCcKAylDE+gi2ZSECHVe0A7Rqk8RZx8CMuuQI2rNA9cVdo7R+iyKmOLnyXwVyhAUPo5lnkImv49H1C0F4boyZrHA34h/99T1gY3BM9xy9QVf8Q3f19j/wE/8wu/btTG70osn/etaUyYObGkQnr4e6xVWeAEheNHKmBt/x+iSKPrMjbekLJs3xJYiKEjCgRH5SMoNDg/F0HNsCkT5YDlbwjrKZJVzmcQLOqYoUEc6Lp3a907SDeNuWuIvp6lG/NmVf+44hDTHfrv9MLewTbF5AT6FVuYxAAD////9lvg9AgAA")
	data["docstring/len.md"] = decompress("H4sIAAAAAAAA/1SPwWrzMBCE73qKAR/+GILhb/MCPfTQe6HHIKvjSCDvOtKKNH36YjtQelx2vm93OhwyBZXXHoXWilRYJDLlYhE6bVPltVEC3UfK+ZHbFtLmkWWNMXOmWEUSeARtYn7Mv+iAt78upIrsv1O+I+i8NOPnEb7eJcSioq0eoQVqkeWWKpEk+GVT6oSRSS77lRWzmCqmJsGSCm7bk36FvIClaBmc6zq8CF6//LxkOoe9+L/DfzzhGae+d+591XBPPCx71fPpPLifAAAA///Twvx/KwEAAA==")
	data["docstring/let.md"] = decompress("H4sIAAAAAAAA/1yRP0/DMBDFd3+KJ3VIo6JI/GdASAwMLExIDFWHw7k2lhy75C5NwqdHcdJS8GL55N/ze88LLD0r1o+BasY2NvXTapP2VQ4Scbsg8NGSx4EaR5+exXw472EbJmUQAnfzDbFxzxcz5sIOWjFkz9ZtHZe/AtAIrUgnoMCrohslteIAPpBvR+W/8GhJ0DmtXDiDQaFEw9o2IQENS+sVcZtOnkSPgi6GwpjFAs8BLz3Ve8/GYI7fI1te4grXuMkN5jVgfYs73ONhs0nDpY3BkqLHkOfGvFdOwJPUFODUiXeisOQ9l8j6LLkkHNhqbE7zIfvv3tP3gOkRDsnylCQKQ/ir5WBZCrzFVA9pgrSLZ9WWESEquB8dxFbFlXysY/qYmoOmOgvzEwAA///hYwBj/gEAAA==")
	data["docstring/list.md"] = decompress("H4sIAAAAAAAA/0zMO2rFMBCF4V6rONjN9SVkDymyhtQiOkaCkWRG49fug1+Q9p8zX4+XpGYYq+b3e8Cv0hsbPApXHCf3k0Tu/i9jjbURFGYWa/BKWCS4eJm9MZxkw6R1SYHhA1WhtFnLtcuT7ZeURpR6zw/mefl0ru/xVfC9+TwJnQNegSM2dJEitRuesqOzSOVVTnXDPri/AAAA///Sf+9X3wAAAA==")
	data["docstring/map.md"] = decompress("H4sIAAAAAAAA/1SOzWrDMBCE7/sUA75YKQTqhtJrKX2D3kIIS7KqBWtJ1U9i9+mL1RDIdb7dma9DP3GErf6ELD9PBsq/ThdMHPOaVPEnyfSRhItk8MqXO8B1DFlwYa0rTIIyCpLkqgXBgmPUxfnvFscULu4s57ZWXPAooYHbe7CPZ/f1LVHX4d3jc+YpqhDhX7u3Hvv5gH6zwYzBGOyfMeAFu4Mh+hpdxtWpIkmpybf2R/1jP2CHV7yZ45b+AgAA//+tps20DAEAAA==")
	data["docstring/meta.md"] = decompress("H4sIAAAAAAAA/2yPsU7EMBBE+/2KUdIkp4iCMh0FBT0dQqe9eEMsOetTvOaCEP+OfE4HnTVvxs9u0a1ijDlua49NLG+aYIugxI4rQi3Qy3xH5QyfwKrR2MQNJda/M1vYavvmQ8BFDoM4cJmDU4qTZ/OfgmRbnixv8kDUtnhSPO+8XoMQAZ2TWeFivgTB2/6O7nTCjse+L/D+hQp7otfFpyqstuPNIcSb14//nSPRGfge7esqaOask/mozYBReZXj7gGjixOa5udMvwEAAP//4WKBkTkBAAA=")
	data["docstring/ns.md"] = decompress("H4sIAAAAAAAA/0yOsW6EMBBEe3/FSDRYCvxDihT5ghRRijUeEktmjeyN4P7+dBw6Ue7O6L3p0GuDykKPSvuv2iDH3VaZ6L5Szmdw/WMqatwNQRojiiLZEzPic4b98dKNhQ1aDJIrJd7APTV7QzJsD3wgpkoxRgTOpRKBSX9PLePoXNfhXfGxy7JmOgf0mYZvPda3snBQbsNL6X8cAPRrVaj37h4AAP//3EC5U+YAAAA=")
	data["docstring/nth.md"] = decompress("H4sIAAAAAAAA/3SSMY8TMRCF+/0VT5eCRDoiAaFGV1BElCBR3s2u32YteceLx85d+PXI3iSigNLPb2a+efYGW80TjL/g1fENjqOUkL/skJiT55mGMcUZUk2FOhD9ZTV3P30I1VeSIk/EWUIhvP5tzpNkDKLoiTEWdZDczLZw8KOnW5vtcRybvnJ4QywZcdX6Wmi30633Yz0p6PPE1G6u9DcQu8LRISaIginF1HTxRrfHN3KpwLNXt6K28TRIIn4zxfe9VGfXbTZ4Unx9k3kJ7Dpg6zgi4N32Az7iEw67XRVrnAEHPLymqKeHXdf9mLyBa90VyPB8vX9Gz0GK3fY+YFsX+ZynVdhV2p6X2PgIqrvHcA8weMsr4bHRO3y/xf+0LMEPkn3U7qiWKa2+mNdT6/KieXrBWHSonsfr/u4essHKMEGsTTGIOpw55JisPasEi5jkTEg6lZlaPXUoHZxPHHK4IMc6a97/P7Z/RPZaf9fFM6ybm8w1PqvvK9akJfHsY7mnu+/+BAAA//834ihQ0AIAAA==")
	data["docstring/or.md"] = decompress("H4sIAAAAAAAA/0SOQQrCMBBF9z3FBxfRor2DC9eeIdRpE0gzZfJj9fYSK7gbPm8e74CjGia1pe9PWMXaWeBRghovY7SxRsY8w6m57vb0qXpKAYNgZyfTBUkmggqLc+CAa0FRzfAFmgXyf1N4OFpleDu0Wc7YYkowYbUMBs99H3BnENtiEUTu0Go6ijya5udsaa0ly4vfoKH7BAAA//9She7z1AAAAA==")
	data["docstring/promise.md"] = decompress("H4sIAAAAAAAA/1SQMW7DMAxFd53iI1mSIMgdMnTo2hswEh0TpUlBluz29oXgxkVHiR/vP/KIUy4+ycxn5OKpRZ5BhmaFZ9eFExbSxuGDays2g/CbhxjqyBi8TPABhKFZrOJ2w/sAylmFE1apo7famVSebWKrV9RR5j2OVVRxeajHz8sVK0kVe3YuaCtHdTwYiVUWLl1JCIRIqn3ULV5SO7SOVCEWtaVto1f7Dfc9Hcngpt//6W6RbyEcj7gb3r5oysohAKfEA/Lfwc79L+Mwsqofttc5/AQAAP//etNsWFIBAAA=")
	data["docstring/quote.md"] = decompress("H4sIAAAAAAAA/0zPQU7DMBCF4b1P8aQu2kgoUoELsGDBgh0H6DSZkJHssfFM2ub2KC4gtrbs9387HL6W7Iwp19Shsi9VDT4zrPAgk/DY7iCKkZyQ8sjhnUlFP+EzOaKYG0hH2JrOORquEiM0O84MvlBcyHns8TGLIdFQM8TA08SDy4Xjet+jxCBDqVxYx+17UvCtVDaTrLiKz9sRlWxec5kZh33Xh7Db4UXxeqNUIoeAX9PhiEc84bnrQmjbretubJNRnCvFJkAln7luJIXXtfEyqJSfPl3SmSuOD1ukTBDHlSuDMC06uGTtgTffaLZq1jXlxe7R2/N/kNP+r+zUh+8AAAD//+9CmPqDAQAA")
	data["docstring/read.md"] = decompress("H4sIAAAAAAAA/xTKuw2AMBAD0J4pLNHALixxil2kSaL7CMZH6V7xTlwuIyL9xlbA7cXTY6FNCn3kRA19ywZF0NL2rpbliuMPAAD//w7hDINBAAAA")
	data["docstring/reduce.md"] = decompress("H4sIAAAAAAAA/1ySTY/TMBCG7/4Vr9QDrbZELFA+jiBx4L6cEFqZZNKM1rXDzDht/z2yk7SIU+TRTJ5nPjbYCnW5JfQ5tlD687DDHNHyyhRbUvfdSLyRIk0k8FAypP6esJ9rOB5hA7GAAp0omsJSSed4DAQhzcFK0uRDpgZPw8w1ThGjpIk76nDKajD/QrBzgpdjrr9a0lnU4GMHpTbVz+xwR1JsU45GQh28UDECRzb2YQYr/DgGpq7Y2eDtJlEZQr43kn0tnJ1Ls+U1Ck2csqL1oc3BV3FWZC0srTmL4aK9x3ngMEtEutiq+X/V0s1a1ji32eBLxLeLP42BnAO2HfW44Ocj3uId3uPwa7dGr3i1/YCP+ITPeHyzq/Flrw+44Lpz7mlgxZlDgJBliZVax4Hnw+G5ce5HDPxCSDaQgAu1mNQedR3AP8dS4mVGrEiRkGK41q0p2PQ28EC9LZhe0mnpdTmb284bfM0GrlMhr1dQTPk4lP0IjZIqtKJ+0+AnToJcjgo+3kjzkb1e5ztRa0ka9zcAAP//5r+52eICAAA=")
	data["docstring/repl-cls.md"] = decompress("H4sIAAAAAAAA/1JW0EjOKdZUSM5JTSwqVijJSFUoTi5KTc3jcgaJgAWSS4uKUvNKFJLz80pAdH4akjo9hZCMzGKFzGKFRIUg1wAf3fy8nEqFtNK85JLM/Dw9LkAAAAD//yVoAgJhAAAA")
	data["docstring/repl-doc.md"] = decompress("H4sIAAAAAAAA/1SOQQrCMBBF93OKD91Y0N5BsDsXIl5gSCZ0oJ0EJ4Xm9tKKC7fv8x+vwynmgJTfS4+oXmZujpjDuohVrpqNbj9cJ/mf9t9BvUjQpBIP0xmawNYgm3r1Aa9JHepgPMfH/ZJtbkirhd0xEHUdroZx46XMQoRvE3uz0NMnAAD//xQLeSyiAAAA")
	data["docstring/repl-help.md"] = decompress("H4sIAAAAAAAA/2SOS47bMBBE9zpFLWUh8AGyCxzvskiCHMAUVbIINJsMm/Tn9gPJM8YY09t6r7qGfqHk3YDH/QyWxd1Rl2BYE0SauTO7oZ+Sx5xK3OAPcEq+RWp1NSTd43hzMQu/47Thzu7qd6du6JsRau9/DovTM+FbKdQKdZGWneeLvxpzSpvtxZ4TcRC6groQ5gup3dD/b6E+gT8t1C3+e/z9q+v+LUQuKeaKkZKuCBO1hjnQNurLim+Yk0i6csJ43xBtcWRBmsFbLjQLSVfZVVyDCBZ3IUZSwYuT5ionBH1pX7fAHuoeP2BBz8JPffBOYdkpYpMashASlLbv3gIAAP//MChWfKMBAAA=")
	data["docstring/repl-quit.md"] = decompress("H4sIAAAAAAAA/1JW0CgszSzRVACRxQolGakKQa4BPlyeubmpKZmJJak5lQqpFTCp4ILSkpLUIrASPYWQjMxihcxihUQwXzc/L6dSIa00L7kkMz9PjwsQAAD//+kUzLVaAAAA")
	data["docstring/repl-use.md"] = decompress("H4sIAAAAAAAA/0TNQQrCQAxG4X1O8cNs7MLeQaQ7FyJeIITUGWiT0mSg3t6FgtvHg6/g1ENhMUAq20sDxqvGxqJ0/ZWsCun7rpYQt9Qj/9eIZ22BFmA8pvvt7La8MXeTbG4jUSm4GKaD121RInzF2X2gTwAAAP//xZr5jn4AAAA=")
	data["docstring/rest.md"] = decompress("H4sIAAAAAAAA/1SQvU7DMBSFdz/FkTqQLBmgQnRCDAzs7LXjHBOLW7vxD83jo7ooFZvlc7+j+90dusRckLn0SCw1hYwyE+03uvbOXCqDpfqcfYarwRYfAy5e5I+B2YZQZlPA1UqdeOtyPuUCCk8M99IzrXee00YO+Ci4xCoTRmJkoPPWG0GJsDPtN1xMMPgx4u8UavbhC7rLXF6bh8ZIFxNhjUjL2gIaMUFfvfSg1G6Ht4D31ZzOQqWAbqLDiofucMDzHk+P2L/0/TVop1h7pdD8eWP+6R+7DTkO6jcAAP//LTZxf1UBAAA=")
	data["docstring/str.md"] = decompress("H4sIAAAAAAAA/0yOzaqDQAyF93mKg26uLoT78wKX0jforhQZaqYOjDOSRO3jFxwpXeZLvpxT40tN4LNMbdvgntPKYroDhWU4qElIDzoJO2OFQ+LtgPCSJ9jIxxx84AGriwsrst83s+Q1DDyUlx1RXeM/4fx00xyZCKVBNXKMucL1Gz/4xd+tIbqMQcHlEFuIEcK2SPpIRF/Et1b1Hb0CAAD//3DnBHXVAAAA")
	data["docstring/take.md"] = decompress("H4sIAAAAAAAA/1SQvU4DMRCEez/FSClyJ1AkfhJoKSjokSgQipy7uZzFxQ7rNUl4emRffqCdnZ1vdyao1H4STUheEfl1VSMLEdoTnZOo4MANvUaELjsSfcNo3twwQKhJPCwG+3M4D7ORTnsKpiV4iiDouKNcwjoJmwJZce28d36d17KwlfDtWrYXGKom+MYqvVW29QwvozOGJA3/+KwQsQ+iFGhvfXFJHkdlO355fRRjGjRTz1fv8kcrQiUVWDszZjLBk8fz3m62A40BqpYd9phWN7jFHe7r+iQe8D7HAg94/Cha6XWBPQ61Ma+9i+AYM4KO1eVb/pe3PEVjjkW9nJnfAAAA//9Y1ktVpAEAAA==")
	data["docstring/to-assoc.md"] = decompress("H4sIAAAAAAAA/2SQzWrrQAyF9/MUB2dxY24aCt1510UfoXRRClFn5Fh0ftwZTWzn6YtdGgpdCaHzfRLaYa/pjkpJFoU//7ewKV44a1nbytFygSZQxBYSUrkwiuZqtWY2L+L9ylhSjqQMQmFF6n/xEv8YjniOXj4YOqw2io6yw+lbdEJfo1VJ8YBpEDtACjxdxS+wKYxV2R02MnOpftumg5Qbhmm96p0RSDkLebmyg4TATkjZL0djdjs8RjzNFEbPxgB7xz1mvHaRAqMpY1Xl3KCjM+P++PDW/oQW/Nt3E8t5UDR+LU27DW+vnLG05isAAP//+ZF5XlwBAAA=")
	data["docstring/to-list.md"] = decompress("H4sIAAAAAAAA/0yQTU7DQAyF93OKJ2XBjFoqAb0AC46AWFSVamYcxWJ+QsahSU+PEiBi6efvsy03sFruo1RF5c+dgy/5iwetSzly9lyhBYQFMW8S40J4Us6kDEJlRWn/0ZI3/oDXHOWDoR2jKuVAQ8DlZ8AF7Zi9Ssl7XDvxHaQi0k3iDF9SPyqH/WoOXMe4btFO6qbhulzzzkikPAhFuXGApMRBSDnOB2OaBs8ZLxOlPrIxgA3cYsLpAY94wvHs/rLZAIBN1MO2GafpDLvDhKNza+fO/ipuDbavTZid+Q4AAP//kawXNEYBAAA=")
	data["docstring/to-vector.md"] = decompress("H4sIAAAAAAAA/0yQwU7DMAyG73mKX9qBdIxJDF6AA4+AOCCkmcRVLRKnNO7W7elROwl6tP/vsy1v4K08nDhYGVD5575BKHriwepcjqyBK6yAcIPcu6Q0M4GMlYxBqGwo7YoXXRl7vGmSb4Z1jGqkkYaI423EEe2owaToDudOQgepSHSVdEEouR+N424xB65jWvZYJ/VPw3m+54uRyXgQSnLlCMmZo5Bxuuyd22zwonidKPeJnQN85BaTAwCfqYdvFR/TJ/x2iwmHplmiO/+IA57w3CyN1aP+E0yN+w0AAP//fvP0f0QBAAA=")
	data["docstring/vector.md"] = decompress("H4sIAAAAAAAA/1SOQYrrMBBE9zpFkWziwM8d/mJuMDDrHqmEBHLLqDuOffsBx7OYbb3i8a64rYzeB3If8/0+IQ6K0yBQvvCG4au2dpI/AK/SjWDjTHWDDMILwVXaU5zp0BqW0deamB74LNWQnxq9dkU1aEeqOXNQHV5ED8Gpb9U5pMF2ddnALXI5Xo7qiKL4JvzoSqgKQa7D/F9sYoYsVmrXRwjXK/4rPjaZl8YQgFtixoZLYWv9Mv0uOy5eOPhezogN+xR+AgAA//+vfXdiKQEAAA==")
	data["docstring/when.md"] = decompress("H4sIAAAAAAAA/1SQTU7DMBCF9z7FU7sgqUolKH8rJBYsOAHbOvGkGTEdV7ZDk9sjO6WUlW1974013xLVqSfFMZBD58NhtarRenWc2KsVmUDfVgabKMKiEd9+mY8OqacLcKXNrU0EjkhhSP2ESn1CZyXSGvmqLPW69PI3ETZcTdjgk0Uub5Bt+5IDK9IQdI1A+WTdzzNY7SXPXmEjOEUEioOkjTHLJd4U76M9HIWMASpHHUbcVHe4xxYPeMQTnvFS14UWC9UrKiHFWGNbGwCojoE1iWIx5t0a3i/OYM7V5tdGnA6NF+zyoFv1aZfzQyQH1pjIOvhupruiQUvrn7k/oVZdwcU3TtlNcy3cq0zg7rxu7mbFs3nzEwAA///KLwLw1AEAAA==")
	data["docstring/with-meta.md"] = decompress("H4sIAAAAAAAA/1yQsU4DMRBEe3/FKGlyUUhBSUdBQU+HKDb2Ol7pbn3crpNIiH9Hd4kIoh17tO/NGpuzeHkY2Am5TgMGGkdOHcidYmHD/JTICV5x/RNeM7wwbOQoWThdm2KINNKhZ9SMA4seQarVyTnt4EUMuWl0qYqz9D0m9jYpCMpniJqTxqXshRwn6htjxvt37Zfoxpj2eCt8j8UwTvUkiRNEl/ICWDPo5gfjz8YaeQdrsYAMpCCzGoVcTgzzqUVvE+9DWK/xrHi50DD2HAL+jBYAYJMV75cPbLZbXPDYdUv69ZRqBFZ0t168Um2Hnm0eRI+2+u7CTwAAAP//PWKwV4YBAAA=")
	data["docstring/with-ns.md"] = decompress("H4sIAAAAAAAA/3SQsW7yQBCEez/FyDSg/wcpKelSUKTPCyy+MZx03rNu18G8fcQp2EmR9ttvR6PZYHuLft2rQWUg+lyGfzvwU9IkTqvAIBeJag6plo3SsTktjsDoyP23HBV+JbqszrlygY3sYh8Z1oAD3vsqLgQh06DZIalQwh2co/l/RMctpoQz0RWKM+DMPhc+i0a91Kha4ICP63JhQKFNqfZ4KEnMq7dEFvpUlOHQNJsN3hSnWYYxsWmAbWCPGW2e3GLg777t7mE89xvu+3UcYP2N+scrsB2Loh1LVMfLES3mXc38gV+PD9p8BQAA//9UbL2/qQEAAA==")
}
