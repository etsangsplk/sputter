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

var data = make(map[string][]byte, 84)

func init() {
	data["core/00_builtins.lisp"] = decompress("H4sIAAAAAAAA/5RX23LjNgx991dwPdNp046m09e42X+BScjmhiJlEMpG/foOSTmJCGqszVuAc3A5ICn4dDqdVBwnZiSlA+GzOk/WsfXxcPjDYN8t/yrbK6XUswm6gxiR1dH2R6We44jaguv6QINimvBpzXPIFc8hHx/zTKjzmbAnHyGYNS9ZjhUK38CtUclSo2zs8LbuOVmOT4fD6aTOEK1WI6GxGhiFYLHz1smIrzj/DGRKjAEZDDBU3J+Wr13yrZJ/WOs6M1LVerWQNuYAS3LQFOqyDfbZLsLdHTtmd5sCo1KypOw4queSITEfB4uzZ3jvSswvzFqD5MD3Ebz5Z9u17enANYaVAUWtOA/n4Gq54jzUTcZ5qGW/oE/Ar6hiasyn5JF2FzS4UoqHAeMIWpy5fEJ8lKfGxx1j81FOLBGL8TG7Gyeu2cm2cW2zqEzWX4SoTEJUplqryPTtMSopyrRkw9uEXsoW8Sby5Wu+QvWWYtVfNtU4wlg/eclUwxwIWDLVMB3yWL7CkqkB+yFhP0RS9HWnDn2N8nwV54Cv8lBzjbqgqB9iDHqNyiZRmZVyWCnHG2oOtIIVU3qUxeDx1rhH6KXReoPvWD8PNna51ta7MI4teCpZWkuFy92dBiSr62K9zj2po/WaECJGBQl7RlLnWQWPtRIG7xSDOyl/LQKrIxiT0Em/0C8kcaa6D3iczkyg+THnzw/OMDm2o7P4mPT3B8nYN2t2MH77YIwUzKQxKr6iIhwgjZJUH0iBKuGMDLaO9nIPpq+oX2X2HA1vEzjLc13Kt5cdZOu36N+r3HwFljE+T0W6A+AmMdvvL3vimIDR/87qfmQ2w/27q6zPk7cZ55fKurf5NZy4TWNI97WX18zjZcMBvly+fvKabRDLrXYhToTy87c4fnFlcTCcDTT2n+LYs8bCOLq5tUJlR63yCMQ2bbQ1enGIhRFesbs3vb1TpfcvpWt8Ur+UX6TVweuJCL2eD41kl7BridBX8I2mk1n0TGGwEUXHxdxYAxZPKdfBf/PmPpArToj0BXlcdB+oQ9DXnYtSUgq4NdvikXuHS7/OGvjikcMd1eff6hcBjHJZSa9nC1w8NZ7hFdvBk0d8oyiMbXTyNE8lgb8sMyJI8wVvFKEOb0itMSFR+rSum4I856rRHGJzEy2HlxXbAevXobi65Ho6/B8AAP//UNzC3zMPAAA=")
	data["core/01_basics.lisp"] = decompress("H4sIAAAAAAAA/3TMsQ7CIBSF4f0+xbGDKUNfoCy+hzER4TZpYoFywcXIsxuxgw6dz/k/rbWGxJIzJ9iQeMTNyGyFqHc8LcamAMcTAc/RBTsYEc7oHE/di4CzNwvjYe6FLwRcey9DLBn9Zo5eFGo71fZS6lfm9WO0AUfYsET5MrMMvG4J6qkt/+lht/UhYwdQ9A4AAP//DfilB/EAAAA=")
	data["core/02_branching.lisp"] = decompress("H4sIAAAAAAAA/6SU7WrbMBSG//sqzjIIMizQ9mfDoPcRDFXs40SgSEGSM8qYr33oSHJk2dkKLTQh7/nQ8+rj7Pf7Pdjr4BwaaLXBVzgartqzUKeqYh32F94aDUq7CuD3a6fbHbcWHWyUdps/FcDhxmVTAbwz0cN44xJ6Li2CMwPWdd5kUBKtLfuIntqwg0PrGrj/KSHrKQDujCpGw0qkKiFh9KG6SAWUFps8dfRKSp5x/TqjKqm8tso1p+q1uTQF0+hFyssSt5RqZ0Ss0zC+kV6H/AXWbmXjv4BG2+UD/0PziTlegdZq1ZVYXotYTQnUSj744wjfQdvGX9ZfnkTJogZPdTX5uqvPmcr8gjC+MUN7SZ8xr65LYL7k5QtcurGf55Xo4MBV932G3URCb4eCjGeYEyDFSspv/8Ccr//OlHap9d113kybspU2Xzof8qvNY7s+5v+ZNkvDC68P+VadUst1o6LfSXTEfBSqE+oUh4W/Uj4EY9Lp9WdPM88PE4PsMA9k3I6bk5382d0NW6cNxKIaNqk8Bi6DdXBEsMP1KgV2m1T7E16ASVRT6YPaVivHhYIXQIkXVM5uwoUPe28/LsBS3VP94/4WIMyTFHuejiWe2uhL6WU3UxG9t6B79zQeF6dEM6jc3TRZVvY3my1J28JRdx9rpxEGjI/Ssn8DAAD//8rzYP+MBgAA")
	data["core/03_functions.lisp"] = decompress("H4sIAAAAAAAA/3ySz47bIBDG7zzFpxxSOKTaJDdbqvIekaVlYZxaNRAxuJeqfvYK421YxQoHCw3z++bP57ZtW/B9SokiTIjUoJ+8SUPwLIS01DttYoBmppgOOt5YAPLaoRw/jGoJmFFPTB3kXfvBoHHErG+EXUWiJDHcxAkfhLseItmdKhL7z/dOZOl3aYK3otSZ5Sf7piBrzfkiI3FC+a5ZapHMp6GRaW32qbmH6lHlU0/cewH8aWwwh1wuYdf73V8BXPfoQ3RLl3IcOOHbusDG6V90MGPgKdJSvx5Bam8hf0CO5IuCWoYZ+PCbTAoRconiTSn1uL+gj1v0saKPr+jTFn2q6NMr+rxFnyv6/NWBa6f+b4QfGxu1+7B6Ff1qgKVnC3KsmOC1I1ROvGcQ8xKWvV9v82VL+XsWSDreKMFR+hks9sg/UxGSOT6v73NJUJgvOUMp8S8AAP//CMSxMzMDAAA=")
	data["core/04_exceptions.lisp"] = decompress("H4sIAAAAAAAA/5xW7W7bOBD876eYc4ErCZyB5n46OCDvYRgoI61SXWhKWdJpdUX97Ad+SaKsJEb9Iwi4y93Zmd2l7u/v72H7s3PEqDqmPehHRb1rO2M3G1FTc1IVdyDmjjfA4U9UWp0t2eMG+CpO6pl2wQihrO0qXB6Sg5RyHqFXpq1WIrBqLUHEGNeXDVq7q5TWG+Dnvuf2VTmC4zP98sHscErhfDShTA3R2p3uKqVhh5PcIP7CaWtd8p7O6cX7QTQtT9Zleld92z22pm7N0zqOpuNTgeCVKtcx/PmU6x/8DaHJxOMSW0AsvAF3a/nX80a86BVbqgsAnjN8Dleval7WBGFc9vLZR0fTueBt6QVi/6i76jnlkguMTWuU1sPvoUyXr3CW6bPXCOANv1jzR17vFUM/er61kj9869/C90qhU1LHw65n6snU64ljRjzT8L3jet7zVWf+TYBwyHZ/arOi4okcygBSHhfJvfm9msO4irPRZO1IY57WUOnPxKmQSBr4/3LJQv4KXkKTw6EpB87KvzJZABiC6Q1bDzHCBctjMvl669FtbHA06KU/mPGLPsNrZHEjIw13ljdGY3nHd0pOcnUn0jG7sSdtCbOfCFsR+xNZq54I25PSfgVQHUSJQ5pSb4suncyPrMz768H69b2zQ9pQ1hK7neInu8llzLXENrIT48JPTa9apnrrC4kbX6v/Bn9ppmjutbdlfVUc3MW4aHBXOAQykcSft4AsdLaTAmGjf04P2F6T26D4HXzOWfGFrVFejTgpY4y6CzDkcabzxHUusmjPnOANfaLb7fqsEPz9G5n1mbudfaBnqhf8f5ETKyW1mdxwKRc4ty4bEIta5DodH6+YUA8x+7Vl/NO8JeZtAvpVaHV6rBUOF2Jen/7Lg1B9rweEAq7VmwGVy9k8eEgIsZfwb3jhbBKwT3ORoKZlndUpOzsiGD975p9MjoerD6Ybxvdqc4Dp5dwyWSgHTco6dCbzEKc6kJ6eiPmGve6kuNSwfEEnh5Ca7PIhnhzGF2H5oo8SezBM1SdPVdW9Eo/DWne4PITE8ojLpG3OKsueJ+YY5BO+lBYmmy134xBcHsRM6Qx0nLS2yQHDN6sPIeNfL+P/AQAA//+nk61XTAsAAA==")
	data["core/05_concurrency.lisp"] = decompress("H4sIAAAAAAAA/3RPTWrzMBDd6xSPGL5vBM0F5E3vYQJ15YlranvSkbQIpT17saxiJ6HaaHhv5v3UdV0jXFKMrPCi7OBl9kmVZ381hjo+T61XQS8G+HSd+GMbAkccejl8GaD5h1fpricD0DiEiP9Fz03tOx/9KCEpo1kWAPIyh7uVXrKCtfbGkGfWNvKDbcHvzV9o5IjGv7VzBVo++5Qt81ticAVyeUDe2vM8DRHk8reSJW8vWIWVQxpjBeoE38858Gm7B60OdgeVC7ti5AJ/FO3bpucUkz70XNE/Wl5UpiEXKtM+7kZuUUuIX8pa8xMAAP//D2AMr/sBAAA=")
	data["core/06_lazy.lisp"] = decompress("H4sIAAAAAAAA/4STza7aMBCF93mKo7tonUWkdgub+x4REq4zuVhN7GCbAq2aZ688zg9YpDc75hwfj78Z9vv9Hn64hEAOyjraoZO/7/B0vpBR5ItCNNT2UjmLYCvpvVUF8GfXWBV/UcDbXH/7WwD1l3jYHwrgKOQwdHewCKGsUTJgfI96WZZZdKd9eJEcy1vBUfs09xepYN2L5CRsZSf1/+lOmg/Kk7nIqaI+4PE7il7+pIoN+AZtWnwvSzb28nbYNI69vK1ObbC6n5xj1F6Z4QMNhy3zGNXsZXEJKk/n/HFzfab2wzb3SE3wKL5Om7Tje1Rn/cUR6mgAk/SZZY7jnKUFgxDF60l3kW49OGqgbNelm9bWAHE9kak6Cqj9YlkECD4qWu18gGcok8zNLELxNCeI9X6kBEcpgL9HUC2vVu3pXNFtcB4rk6P4IENOBkqtttZVJNUJ4+oW1OsA0ViM7xODlcIgXdBBW5NPYRGmNVP2YkJ6PsQiYi7P4rQUqRCH/kDsGevCT3seD5/O4U1/zkQrv+YJZ94S3y0aZ4e1jZntvwAAAP//r/v8qpIEAAA=")
	data["core/07_io.lisp"] = decompress("H4sIAAAAAAAA/9RS3Y6bPBS85ynm209d2VFp071MtO2+B0KqA4fEqrGJbYhQ1Tx7ZZsQ6K/au3KB4jAzZ8549vv9Hq7rvSeLyljaQb41WcZqajQ6m7eiyy/Sn3ItVQZ83nVWDsITvO3pSwYUTa8rODqXGcBa0WUAwBqNYhCqBJMNmHSBj0EoHl5gkRSOnEe8ozPn97EoHtEY27ooqsijcHQG+8YQnLf/JSAv09xpWkAHzfgnwN7cltxtnK9N7zfYXaz0BNZI63wywCfvxuYkqhMKUtSCWZoA5e/lNq4TFW0QmJwvd9LrpUTXqTGsmuxnv1bVdFFS02YpKLX/g5z+lZik9uonUYWF/zYtIS3V+UCVN/bDj4sc6hrH6TqGk8CptOwZ7FVIWKfzE8d2tt2KyhrEsE1HGsVB6lrqo8MjDqYek6pzZH0u7NGl+NaWcONwPNyVZqG2dx4HgsAnGvNBqJ6QmA8xjMroOsk+J5d3ue3tdj+y2uD6EhxNN/j+O/DTHRzbdGWzhS3H4vSOv56Q01Mp4+h/sCN5rFm7+AlaqrkZAPN2XAiwRXrXl1Sl9J698XJlfuI1UgulxpUXdjmRjleYmpOscbDpR3qyrwEAAP//cTSIgPwEAAA=")
	data["core/08_os.lisp"] = decompress("H4sIAAAAAAAA/1SNQQqDMBBF93OKj0KZQKXtttn0HiJUdATBJDIZtz17MVKwfzOLeY/nvffI62YmiiGpPJEyEY8yhX7QBJuDENBeMCUNuSPgzYsY2my9Wg3wsKlKtGZH3ZXwm0reFqvBY8LnVfTzW+JY7/ff7wrAq87Rlgi+gZsDPXoOj3uZQxVy5ehUco6+AQAA//+xXZhQ0AAAAA==")
	data["core/09_predicates.lisp"] = decompress("H4sIAAAAAAAA/7xWzW7zNhC8+ymmW/QLdRDQ9pig6IMYBspIS4f4JFImKX9xi/rZC5L6S00FLVD0YkrameHOaMXk5eXlBX4YQ2CHxjp+xuC41Y0M7A8H0bIy0L7mKxscr7Ib+QTxC8QPSDf4ucKPVbUCbdvu4X6qqllx2aPu5O+32vPlAPzxPDh9lYER3Mh/HoCjGk2Dxnbd6QCIDRYQWkFoH+8TokpPI4oDjgoicYXSzocMqE4TJJPVcgeIxhoPBfHYGLKO40Wm2hCfhJKd5/Xhk4jNV9XstZeNs2hZ1av0YP0ndo3sGT0H2cogs+/kaDRNnWrC33oIH1yG0q80W/st7oj7Ap0j+abDWx0l1wCU2eBwvD/loL7g/hStnrbhdNKHYjT3nE0Kb1bI/Gob031283kshs//TSz0Hf1v0cQejPo+sY7vJwhjA0QO5r3aztznUSaNnSD/bZIfM/vyMTXpPbtQS3f203c0fd2iY4NVHzRfox99wCtjkNpxS9USft9CBFtL722zUpfA7Rzyw/hPk3NPDd77tiojDZ8fkIvxFRfPHaM7kNEdLVk/t7aJrXEAZQBVBWIaGB/cLjHWikS+gDxf9ol8KRJjztSx2SXGWomoTcvv3IKmC/o7cX5eIueXRGmh0q65UqL2chjitnmlR+pUKJrVPoDiL+2YjaUS88pNsA6U18KuU6E4DoOzvfYMmi4eopqfF/1ykHn8Sz2/yYzY2djG16DiO1KUiaDmjZuvHhJm7F/ZQVmHwXod9JWhjdJGh1uxFcPnLFcnvX05w2f5D+SkARk5j94HufAmA+T051p7xINs3qLsNR89lJa5t8A+eGi1VZLIkOJYDkN3A6VlV6ORJp4/EaR35swP3GjZ1cq6HrS9+6gas5KY6kj1orevfPtmXQuaLj73N4OKnd36V9uB8kpL7GWlCVX8lGwjO1BaaPP+ikIGo7mMstNKc7uqFmTT4U/xt+BxmjDtkQCltuJ/fGTbtuhs5UdEdfgrAAD///GoDVhwCgAA")
	data["core/10_threading.lisp"] = decompress("H4sIAAAAAAAA/+SV3WrbThDF7/0UBwf+rCCCf3tpg8l7GJOMpZG1ZLUr9sOtKdWzl/1wrMRuA+1NSn1jeXZ8Ruc3w+x6vV7DjcF7tmiM5RV8b5laqQ+LhWi50yVQe1Mr6fwC+LYarTySZ3gb+PsC2HqyB/a7BSCCVuwchHQpH/moWgCAuIqU5yoXG6ixBvUmFWlNU5Nz7LEsr9BJ6/wyFhTbI6nAO6Sv6hLBf+iMHdyu1GOPrYJ47QEiKeXMqrpPufnTnc/Uq7CFsJyiuxJ9EvUGYuow5brTg60wPeS0Ily99fUzY4o+lK9i7MEWc+8Zc2bgS9OwVPKZ8VRvnu6xDx6WfbDa4VFL9QjZgfQpKYGjPHl2IAcXmv7DQEjaR1J3BcHuJV986VnX2vg04lqq+Jp31UxPZB4JYpa4PRu3IN6g+I9jvAxjSniXI7n5KGZTDsFJfQBhb4JuoWngKIDROOml0aTQKf4q91JJf5rjS7k3GCaJ3wD5hs2UdK7gRBOYFMrxteXzch7omevG6LZuFAXHV+sZKzdyI0nVaVReNrY7DWlddxrb/Ne5i9FyC6F9j3yG/6v7PGvz4KeZG1FukVUUPJLaQcgOUxaa3GmICDFFjSo+XjcuuZi1Lv4uzVEnnBvpe2vCoQchLlDTZShx8H3P0iLDJnsIA+v3l2p2cu5g3MrW12QPrhiLw8lH1pGLPmdXFZblEUNwHnvGSNJyu8wDnC6Hl4tBDDRCvO0U6k0FMZL1ySQ+X9RvgvljMvGS+UvA/JLMjwAAAP//qMN6ILQIAAA=")
	data["docstring/and.md"] = decompress("H4sIAAAAAAAA/3SPMW/CMBCFd/+KJxgClEZqu3VjYO7YMTolF2zJsaPzhdB/XzmmCh3IFJ/eve+7LXYUOvRRhsNhj5El/yYQko2ir62TdnLqwgUVha4y5yv5iZQT1DJKuJc4wHOv0AhxF6s1TgkpxgBKiIHB61oEoerJJ/6pkMd8xOy8h7BOEqCWtMxrfKllmV1iOC2hUWLL3OWae2d2yy6Bb7oI1cZstzgFnG80jJ6NQbly94I3vONjb1C+xeLvsQlR7w7cbfbGfD9INUu0OSLwleUR3fzfa46YLQtT+nzOVZlW7DNkvimpFMbaX5vfAAAA//8ZW3AItAEAAA==")
	data["docstring/apply.md"] = decompress("H4sIAAAAAAAA/1TPTUrAMBAF4P2c4kEXJggFFdy76A3cl7GZ2kCaxPzUentJgwWXM3zvDTNAcYzuB2v1C7J8abTZSganz7qLLxklgC9QbPA0HewqF8komyCmcFgjpmWr+EXA3twdtmQ0LlfJP//XBxVS9wt/ONHgHrvPj0TDgDeP6eQ9OiEClJEVJx7UE57xonVb9T8ecWqi981mSPf4ts4hSanJY36dR/oNAAD//8CDF8f1AAAA")
	data["docstring/assoc.md"] = decompress("H4sIAAAAAAAA/2zQwWrjQAwG4Ps8xY9z2YRNWNjbshR6KH2BQs/KWBOLjDVGktP47UudQFvoUeLj/5E2+EXuLeP/mRdcqM78sNttkY0p2EFQfsNKhEIuDA+bc8zG6VVqvcMf3JGcezRFDIzJ2kV67nHmZb+2YCIx/41mMI7ZFKTgcYrlW4oUaENpNjrIPoMOeBnEUWbNIU0h/uF6KYWNNRAD3Zq/hlUJNqrwRYOu4GvmaaUBCWRSHBmxHtRDFIQi5rHPldxRyAdpekhps8Gj4ulK41Q5Jdx/mADgn9LI6HyaI9i6245ODPw5/L1NlfSE7rl12/QeAAD//+eMAcV/AQAA")
	data["docstring/chan.md"] = decompress("H4sIAAAAAAAA/2yTQW/bMAyF7/4VD+3FDtr8gNyKoafdht2GAWZkOhYmSylFJc1+/UA5SetiucQQn8j3PkmPaN1EsYMTJuUMQol+8MJOfYoUYOXIoXm5fcGbaiAlZJXitAhDJ1IrlMwDNOHAkYWUQQj094LMb4WjY6QRJwqF8xY/J4ZwLkGXlhPl6XmmI1yK2Wf18WByiuh59tpjLLG6egKhdyFlXq3FAXQftMUPVvF8si4UwYFnjopR0gyd+MPRTBds9iG5P5snnMnXuWOSqor8rothS7VnmBPlAUnumhuXRVB9DVu8mrAO/9hPd+3ZhwAKOaFORonqQ+22L+PIYkTGUPLEA8jgXEGlEV6NbIVzYlll2TbN4yO+XWd850tums0u89sG9lvFXtzezmloNjuLVoWGq8aUO14bbNvHJDP61spLrK5vNrua2bbSxwZNC4oVor6ta11fjb5EvL7TfAzcNEAbWPHLTdcb+bsxK+0h1X+g3cJNqCbxMKb00P2vsCf5Wlgmds3ST9PziZ0mQWtg4KbOaldmGcJH4Ww3xXyHdH4OfOKAEgeWo4/Rx0O+0ehv+D7fw/Pk3WTn5yPmlBWOMueqPwqPLMIDznRZPZPVI7m+DV7YgPbpxHCphMGOTPgsdjhxV6l9CnTrtiT9Sqr9RKjrmn8BAAD//w/0PTP6AwAA")
	data["docstring/concat.md"] = decompress("H4sIAAAAAAAA/0yOzUrEQBCE7/0UBTmYIAr+3kV8A28iy5AtzUBvzzrdcV2fXjKB4LW+r6q7Qz8WG1PA+XU5QNNv1jPWjJaCvpCZNtLlubIlafHOG8BpKs6lFLRAdsREVPqsgfLxby3bZ2NUHmjhC2UaJxxr+c577rfNa5Guw5Ph5ScdjkoRoI9ypdlj+/ntBre4e8dFf48HPA7DIPI6Zccpq6Iy5mrtXmvt+qZjdXfyFwAA///8jynh/AAAAA==")
	data["docstring/cond.md"] = decompress("H4sIAAAAAAAA/7SSMW/bMBCFd/6KB3uIFQSGMmRIEKTo0AKd2z0+U2eTCH0UyFMl/fvipMQtCnSsJlF47+m+e9xi57N0eO4Ld9DA8nJ7C06VPzXouZxyuVSYJGrMQgnHQuJDlLP7mguYfMDhD/cBPtFQ+c6OsO/RkzLGmBKODP5JaSDl7g4kHeIJURErtAwaZuwkK06ULMBeJaZmSbqx9BvYPCa/xiwphXUoYplZA5cxVl5MwpO+z2OmvmTPtXK3d267xWfBl4kufWLngF3HJ0x4fGyWkyE7ANg9Y8JD2wCbkSoS1woNJHhoN6vgBRPu27ZZBefCpFxWzX37Lvr9LKIoOLKOzLJpnPsm0BArPBn269+K1+vyPjj3+J7iOWiaccmFQVw1sEZPKc3oE1ONcsaYh9SZTzNsBSR4smrxxvOYS/f0P0jXP/yD9EcwCKrZiElX7DGXt2r9HNkvVVl1hyXn8DHqtfBqMOtl2btfAQAA//+f2055vgIAAA==")
	data["docstring/conj.md"] = decompress("H4sIAAAAAAAA/2SQvW7jMBCEez7FAG7OuIOAy3/rIm8QIEVgWCtyZTGmdhVyZcV5+kCR7SYt55vhzqzwx6u8o/AHWs393zUohAJO3LNYgSloVkcWz27zW5vtGoWaxFeuwksXCxru6Bg1Y4opIcS25YzAA0uIsocKrOM5wGc2hp0GrrBBisUWy5B/4H8gHNnbJYmG5XXqYmKQgEpRH8nikdHTgQtEsR8pkxhzATU6GjQHzlH25+vaUbxFlSVT1DBpPmCK1iHR1+napqCMvgMVqHDBkDWMngOaE+qehhqaUbcxGee6cm61wkbw/En9kNg5nBd++48b3OJui3s84BFPa+de558z25iXKc4ldxf2gm53lfsOAAD//0Od9nWqAQAA")
	data["docstring/cons.md"] = decompress("H4sIAAAAAAAA/3ySvY7cMAyEez/FAFvEFg4L3OWnT5EuZYCUZ57EtQnYlEPK+/P2gezsFkFypYYc8iOpA9qY1XHKNsP5V4eY5zdRdpCCJ55ZCy5SRlCNr6yRm5/bW5EtsXF6BJ7gaxxBDsIkXpANZ44l2xPKyDD2dSqQGle+/JWzt6l5G81ivLAmTih5U7PJIErTEff+q75DoCD3HIWKnBlkRreNwrgCaMawkpEWrlwDWRIdsGSXIlmPTfNjZCjNjL5uqN+pz+xFBkY+4TKy4rv4ApmXfVFUrY5QDaHYGgunx6BeJ408TY63GxYSqw0JfSTr0YYYshbWmnfa5g0UUjJ2x0JWqhos8CBe2Lr7TfqY/uNOgaPt9/uH/9g0hwO+Kr5dqeI3DdAmPuGKD+1HfMJnfOm6u3j7801ecN3F7fWMW1f3JA7eq+Ai0wTjspritX3GC+61Xo/N7wAAAP//UYiuRW0CAAA=")
	data["docstring/def.md"] = decompress("H4sIAAAAAAAA/1ySzYrbQBCE73qKAh/WXjaCbJ5gDzkEcgzkEALbkkqegflRpntsK08fZmxYyKUR/VP9VWsOOC5ckSQSay7xhMmnRSE9pZvMBJOVffjpQ+hFCC4SKmH50QafYI6YaylM9jH6gqvzs4NXPFVlecK0Y+EqNdiIt4egT2eFFMLHWE2mQEhaUKg1WNOWBJaSCyb6dEYRr1zg11YQM8bN2oooS4cq/NQ5zTGO+OG8YqKTi8+ltS1+Xdk5zUlCzGr47nXTF4i2oR3X5vXMxCIh7FjFB6gPTBb2BqR1dphFqeMwHA54S/h6k7gFDgPuF70NAHCMsvUP4Lgm/Np/Y8fzM15Pj6zyz+cWXlv4cjoNQ+flXe0OMheKEYIgf3dE2XB1LARldmBgbF7y2n+BuUJiK/niFy5NtDLN1G481ylwQd1yQmkVtRHf7L7FHBMedxP72GW5C//3HPB+ex+HfwEAAP//rkOu6z4CAAA=")
	data["docstring/defmacro.md"] = decompress("H4sIAAAAAAAA/1SQwW7yMBCE73mKEUh/DH+FyrVShTj0GXpAHBZ7o1hy7NS7KfTtK5tAIYdoMt9kdrVLGMfdQDYnRBoYAyvtcCh6vT6iS3n4v8LJRycgZCbHGTXffPoQKgFdHZx+riWaoD3DTjlz1OrJSJZfcO697eEF7SSc2/KH446moBvs5xYv4MtI0bErvDTNc31EykVoAgUtomdImrJl2OS4FTjSsueYWTgqqU8RJ+5SZnit3d8UJlJ2m6ZZLrGP+LjQMAZuGjxcw6boGgBY+AIHjip1XFtIWy+zqIHDP9hAk7Ac67c59xxhhL92N7CqADC+g3nHFiZwvMMbBcxs4fXBC14UrYyTKuc3393JU/5Pb1dPiRTlujNMZtH5fZ9dn+Y3AAD//92LbYMIAgAA")
	data["docstring/defn.md"] = decompress("H4sIAAAAAAAA/2yQz2obMRDG73qKj/iQXWKbuMdAKT20T1DoIQQyqx1lh0qjjf6QOqXvXiSnsQ/RRfrmm3/6bTDM7BRKgRG40BfcU3rKD3AxhZsRk+icQT0hr2QZrqotEtX8FO+7D3oPYjqeepWIsjBsTYm1nMu3eFnELpCM65o5XbeKmR1VX/bGbDb4qvj2m8Lq2Ri8redkMgDw526OFleJS02aQQr2HNqA6Po8J1NUslaQ+bmyWr762yvv5aHfg4069xcwfIbgdsTtpT6MOFzqT2d9xz5zi99gcDJh2HV/vFCHsR1jfiySwadvYOXUaGYktjXlhulFygKNKCQelrxHXIsEeaVGcQvSudmBQ3yL7fE9JhBCTO8dSUvPzIXsr51Lwjr74wcQ8MTKqff5v9UWmbkzm6OtjeHJjg6Pnl6Pu8zPj3vzLwAA//8qpIYoIQIAAA==")
	data["docstring/do.md"] = decompress("H4sIAAAAAAAA/zzOParDMBDE8V6n+IMbSTweSZkmkCJnSC2Sdbwgy0YfQccPdsDV7vCbYgbsa2Fc8uy9Qz4htlClMLdYdY2yUzEPjfFQJDynHdBEbTn9kWW7mt7USRg1haOvSyIUtBaylBbrvzHDwC1x72FeoxgDVkfslc7FbYsMgF1zorvf7z2d82lLzplvAAAA//+liQqetwAAAA==")
	data["docstring/drop.md"] = decompress("H4sIAAAAAAAA/1SQT0vEMBTE7/kUA3vYBqWg+0evHjx4FzyILDGZ0kKa1ORVWz+9ZLvq7vXN8Jt5s0LlUhxg4xgEmR9XGuWQIS3RdCkL6NkzSEZsimNksMzqpfMeiTKmAANvvuc/EdIaASfrR8dz0voYs74gFnVI8bNzdP94VDYGa4TBCJ2u8bQ4cxyT5ZnPJCK3MQlTiQ1HVypyFrrlr2uYAPaDzPBdFnyV6u88taerlVqt8BDwOJl+8FQKqBwbTFhXN7jFBlutf48zXnfY4w73b8ut7LfBhFkr9dx2GVwwS85polLrcqRDtcUOe32o1U8AAAD//6EiwBaGAQAA")
	data["docstring/eq.md"] = decompress("H4sIAAAAAAAA/3SQsW4jMQxEe33FnF3Yxhn7D1dcESBAypQGrR1ZQmTKFrlx8vfBeh0gTRoSGIAzj7PGllek1s/38XcHp7mhJAiMjpbwLnWiQTpRRqqXKBXe4JlIpZuH11IrOn3qikOSajxADNaazrs4qLFN6uwGWeI8i6MYtPkvtgNePLPfihG3nwneJx6GENZr/FP8/5DzpTIEANiOTMhYZdbaVrtF4/VbQN6F8FzeiHMzx6VzLFGctofnYkiTRi9NEUUh1RqOhPIkzhHHz/ngQh2Lnu6Ymz8bxCxdorMPeNLFJYpx/8BcyI98wHOcq51/llrvRSzF8jpJHcJXAAAA//+7FebRjwEAAA==")
	data["docstring/eval.md"] = decompress("H4sIAAAAAAAA/1JW0EgtS8xRSMsvytVUSK0oSMxLKVZIzEtRAAmXJpakFiskgmW5AAEAAP//Q1NDtisAAAA=")
	data["docstring/filter.md"] = decompress("H4sIAAAAAAAA/1yQsWozQQyE+32KARe+4/8x2E6ZJoQU6dMZc2zOupxA1l52tbYvTx92TUycUhrNfGIWaAYWo4gha49En/9aiP9imXEVUllm0p6Se47kjRJ8OZlvAs5jSIQ+qJEaOMFGQqSUxRAG+GmSmfWjrqcYTnygQwUaB4WFKpDQkdRScdwd3vgrvA5/ostU4rn3NYsTlhazjfMSjQZDN3hJ1P1HHZSla4tJq/PkJRPOLIJ3AmsvuQBZf1HK3z8frJxbLPCkeLn44yTkHG79NYNid9mjecQF27bFbo0NtnjYt869jZyunEiW4zX/vsOuWWPTdu47AAD//8wwT2mSAQAA")
	data["docstring/first.md"] = decompress("H4sIAAAAAAAA/3yQMU+EQBSE+/0Vk1whJIZCL0YqY2Fhb3/AMisvPnZhd/Hw35sDo7nGdibzvTdzQOEkpozEuURkXqJPyAOxy1SO9BnBbWLivNBbmrdBEtzibZbgcRbVn/A/2YlWnLD/pdwiRHhRyDUdksBxyl8VXjPOYdEeHdHR04mVVpED7ED7ARciWny2Kn9YLEn8O5oicX7aijXo6EIkbKu6eduHzeV+E5lyUxlzOODZ42Vtx0lpDFD0dFhxU9Q1Ho64v8PxsSwvxt5vLY3BNgT30NUOp7o+VeY7AAD//0ew0thgAQAA")
	data["docstring/fn.md"] = decompress("H4sIAAAAAAAA/2yQwa7aMBBF9/6KK1g0oQKVlKLuEIv+QaUuoghNnAlYsseRx6bw908k7N7b2kdHc88a1SgQCnxC4EwntJSu2mGMKXyvYRNTZgUJSKI8QyyKsYjNLor557x/I18TyDfKCPREz5hIlQdQikUGOAFhdEnz1npSRSARTjtj1mucBX8eFCbPxgDVwCOGWHrPBgAqzxltKB5NNz9gXrFaEMWdfGFdoX10qDYbPBCKr+u6nmU5bu9sc0yLK9D0duNbtUeDnzjgF44z/vfmFLxcgv+vtYlzSa9djMWCS9vggCN+Y/8D+6a77Iw5fyqhmFIcimUQrI9aEi9xbJwc6yzsnQxOrrp83OjO6JkFiUdOLJYHjCmGmdWS5o5OrlAbJ96ZjwAAAP//Jcru6MwBAAA=")
	data["docstring/future.md"] = decompress("H4sIAAAAAAAA/1yRUZLbIBBE/zlFl/0jqxLfIR+5QGoPIAytFRU8uGDQrm6fAinOOnzOMP1mus8Y5qo1E3PK93G8gKuN1SoLLAoVae6tAls2cUtOkmqJm/lFrVnar0MgCHTZddpQq4vTkOSKt4V45LQGT3+ofYQYceMT59t8Iz5stkrokml9E+InXW0632DFw8oGZ2Ms0LQDD8wuOd5icr9HVNEQnwsVLHYlbqTApfsjUhm3f/CrMeczfgh+ftrWNQYYPOd221+HDAAMmr6vdJoyhncK2669AQy8B8Up058ur6WUrbzz/+rGGNPH6dJe581VL8a8LaGA+xrNNF8dv7jcTqfHNFedoItVuCQrs5Z+bGapUUsPQF4iQzcGQTTBYr9hT2bq5KlLH7n0r18s3Nn5GXnzjVfzJwAA//872/RsPgIAAA==")
	data["docstring/generate.md"] = decompress("H4sIAAAAAAAA/1SSTXLbMAyF9zrFm3hjpYnu0EUu0Ml00+nEMAlZnFKAyh8n6uk7JCXHWQoQ3vvwwAOOFxYOlBijhvlbj/07ghD5b2YxDIqrmCmoaI5+7V6u5HP9J02MuLBxo2NbJSKc1NGFqmyaApOFjuAPNjk5lQE/OOUgXyzSRAnvznu41ID0ygEka5ktPsW0WlICBQbPLiW2A362xl0NOTq5gODVkPcrotGlAGYxhWCXLMA4HctUk+9PA163xqZYqdmCCq7R56A5OeEnzExSXBoQihGSVoS2yOPZq/nziCzJ+epnNASOi4otg+x5ZklwEYGj+itbnNfqIjHPHHbMPaSh6w4HfBe8fNC8eO464Gh5hFGvIX4eswOAttZDYPvQ3xc0kFz4a21l7/X9oe/7qpn0+comadiU+657nVwEN9+23tmJrXie/q2fhwz1tm2V0r09sD2gUhSaGacmvkV+51l+bBYlwFt+VI4UHJ09l8zGXC67ZWWfQA1HNoL2Zpri268aw23328K/34bufwAAAP//XdpC+AUDAAA=")
	data["docstring/gensym.md"] = decompress("H4sIAAAAAAAA/0SQsa7bMAxFd33FhT00BpJ8QDMUHTp071YEKRPTFgFZSkTKjf/+Qc9O3kiQ55C8LXYjR10mqOUfHW6ZyVhBKFEehaHLdE1hj6I8lACJmOiWk7rfA6hCEkeI4p7TLD33e5gnezX+Swi4cqV7WMKjUJBhgXneFoQFI0fOZNxvy47440UxlHgzSfGl1o3qOYelyq/saZaUMaQMT+oPRhLeGq3H6hKNngcvow8yeqvc+sDRubbFz4hfT5rugZ0DdoENf/WdSTNTbrqzA2pP1PBN78WM8/dtsvEcQmrO0K6rhtOphsGPIjMFjgZLDvi3imfK7RdRq865z195vWHNK7OVHHHZJi9H9xEAAP//lS3GgKYBAAA=")
	data["docstring/get.md"] = decompress("H4sIAAAAAAAA/0SQvW7jMBCEez7FQC7OBg6u7ho3QQKkSBukN1biyCJik8ruyj8I8u6BGBkpuTM7/HZWWB/o6Iue8M4bInuZjv6wgdI18UyD4CzHiWhvsyW80ifNBh+4CJfkQ8oQGD8m5o7wQRzJIGalS+KM1VR3bGSX+sQ4p23x0tdp/bzQkIuD12R+j607S/Df+ZXB5AO1KgvwQpJs5p40M6IoJIOqRetckjFuQ1it8JjxfJXTeGQIwDqyh5aW6vjcZTkRzVNpG+zkQPz7/7WZTXNNi2knMSrN0Fy05EOzCeFtSAb+ZC4Ihv2i79Gyk8lYie+dpDPrxfmPoyvZZa4w/4aPWkaq37bhOwAA//+YoO0YpQEAAA==")
	data["docstring/go.md"] = decompress("H4sIAAAAAAAA/zyNTY7rIBgE95yin72BLHKHLN4N5gJfTNugwXwWPxl8+5EnUralqu4ZdlOsWvbbzUHqmZdQNGuv6QRfkro0VgieSZdv8xWIo+grevq/quInpoQnP7JHzBBUHlKkES0Uioeu4ODSW9R8xyOfKKw9tZg3XCUvowW+jz6rPtZFiqe/GzPPeGT8H7IficYA1nPFgD2K7rHSuYttCjswBaak0xvVVmCHw4RT+7/Jmd8AAAD//0iWKE72AAAA")
	data["docstring/has-meta.md"] = decompress("H4sIAAAAAAAA/5SRsY7bMBBEe33F2ClsOYmB5AdSB0jpLgiElbgyiVBLgbu0nb8PKPnOh4Oba7Uzo3nDT9hPbPQDY8rT5xbGaoqrZ/OcYZ4x53QJjt2iUHi6MKrFkVHzcwTFeD/xhWIhY1gCiSQjY4f6kfVLzRKYD4qxyGAhCa4hRmS2kgWd5cLdESfPGENWgyT5+oi5a7VEQ5Cl2WvOGhHkjG6kqNwdmwbYOx4FLpU+Mn7f/mB/OOCG721bjyv1em2b5lSL8Y2mOfKTXuh5oKKPfyoGEvT8AD02za/wlzElNcyZXRjIVvC3zNVGUVP1Cp8Xtv5fNcwsrjJUtN1mh8FTpsE411GCYmIShXmyZfM0PnueqSzLLZJ0fWkXkqyTbFbsrecY0xbfPkCuloOcFSQOUqae87LBzt6t8D8AAP//HDug1VICAAA=")
	data["docstring/if.md"] = decompress("H4sIAAAAAAAA/3SQzU7DMBCE736KUXtoLEVIUP5OIA4ceIs49bpZabOpvE6bvj1yKygXbrY833hm1mg44ZApogykIDF69zhQTlMeDcbjQQh9DrobWPfuK1Uh6BhkDoXiheVdKAQ2lDyX4YxGp4IUxKhFPSqLby/cpv6yQTWv+ptN0IhMZc5KscVUBsonNrpCNdUVasEJQc8tTiyCnv6xuHNuvcaH4nMJtYFzQBMpYcGmuccDtnjEE57xglfvL6+c0LyhEVIsHlvvAGC11Jg971d/rjYGkZV3P2PYeewnQTerkFlXJbNRBKsVChFTQsepa68TV0KmPe+C3Hb9rcN6pFxqge8AAAD//8CeuC6bAQAA")
	data["docstring/is-assoc.md"] = decompress("H4sIAAAAAAAA/4yRsc7UMBCEez/F/Lni58QRCehoEAUFEuV1CEUbZxNbOHbk3eQ4oXt3ZOc4AaKgtLwzO/vNAS9IJNn3GFOeXx6hLCq4OFbHGeoYS06bH3ioEwLKjCrxpH5jMZ9GUAj3X94orKQMTaC4D56KTYQ6LxjXaNWniIsPAZl1zRGd5pW7FmfHGH0WRUzxVdX+mpM1KHysgR4eu9zHCd1IQbhrjTkc8CHi43eal8DG4HHfj3eRZkbT+xCa2+O58DRdmxu+vMYbvP16NOZccvJu8GfMfQl6trQK1yzqfN7JwAsIG1tNuTXms//GmJMolsyDt6Qsp78Y2EIoSELPiDyR8oD+WgQLx6HcVVY8Pz3DOspklXOB5AUzUxSoI63s0/ivpua1ktRi/3tjbcXydOfSOA4hNf8FYO/J/AwAAP//ysZCyzUCAAA=")
	data["docstring/is-indexed.md"] = decompress("H4sIAAAAAAAA/3yST28TMRDF7/4Ur80hjSiR+E9PiAMHJI69IRRNvG+zFs548cym5dsj75aCSul95vn9fuMVLpJ2vGX3AX2px2cbOM0NNwN9YIUPxFjLKXXs5gmDVOJuCcYfEzXSwucekvPdCE+SJ3HCCwQnyenPKHwQRRTFnpAYacYO+59L5mV7UeFDMvSTRk9FcZNyRqVPVbHzOnG3xfVA9KmaQ4s+f9jn94pN2ZF0xriPW5KSHrDrJRt32xBWK3xUfLqV45gZAv7ysr54gZd4hdcbfH2Dt3iH9982IVy3ilwWHmsYwpf0nTgWc4yVXYritMsHaM2DZCtNhvIgvrgYK0dq1zq26uuzNeIgVaKzNvZkOFLUmkyfxZf+sVsdp1mQt/jlDP9cbjvTnt3jng/MuZzj6upJyMUc9owyGWFekx7+8zm24VcAAAD//zBMU69oAgAA")
	data["docstring/is-len.md"] = decompress("H4sIAAAAAAAA/3ySP4/UMBDFe3+Kd7fF3oojEv+5ClFQIFFeh1A060xiC2ccPOM9+PbIyQHSstC/Gb/fb7zDTWJ5hzGX+ckBxmqKh8AWuMACYyn5FAce1oSCCsPnKsYDlL9VFs/qPo6glB4jfKJUyRiWQThRin+isEACT4LCSy6GaIrEMlm4bc8JLETFWMVbzIKHmBIKWy2C3krlvsN9YIyxqEGyPD0v82tEazJEWRl+r9s2RZnQj5SU+8653Q7vBR++07wkdg6PQvY3z/AcL/DygM+v8Bpv8PbLwbn7Vo+38KV2zn2KXxlzVsNSeIiejPX2DKsJoKQZR4bwRA3g+KMNLCxD69dq76/28IEKeePSuKNiZhJtFm01nsdLR5rrKsfa+s3/XyfrVtKrFfU6cEr5Gnd3/wXcjOHInqoy1EqU6R8/onM/AwAA//9ltyBNWQIAAA==")
	data["docstring/is-list.md"] = decompress("H4sIAAAAAAAA/3yRT2sbMRTE7/oUY/tgm7oL/V+fSg89FHr0zZTlWfvWEtFKi95bO/n2Qdo4kBByFTNP85tZYRO86C/0KQ8ftlAWFVwdq+MMdYwxp4vvuKsKAWVGcYj524NCeHrmC4WJlKEJVAW74o5Q5wX9FK36FHH1ISCzTjmi1Txx2+DgGL3PoogpfizWm0ymoPCxxng+Mbt9PKPtKQi3jTGrFX5H/LmnYQxsDG5Q680nfMYXfN3i+A3f8QM//2+NOZRMPKtfRpov4sSWJuH6sbBNcaaHFxAubDXlxph//o4xJFGMmTtvSVl2r4AtRVCQhBMj8pmUO5weimHk2BWK8sd6sYZ1lMkq59KIFwxMUaCOtNac+rfWGKZam5bzdZWm0i9m/KXjENISx/3+fex5CfMYAAD//6e3brsNAgAA")
	data["docstring/is-mapped.md"] = decompress("H4sIAAAAAAAA/4yRv87UMBDEez/FfLni48QRCehoEAUFEuV1CEWbZBJbOLZlb+44oXt3lD8cAlFQ2tqdnfnNAS8mSYn9ewwxTy+PUBYtuFqqZYZaIuV4cT37daJAMrHtmE8DxPv9nxfxsyihEbJPQG+Jp0UlQK0rGObQqYsBV+c9MnXOAY3mmU2NsyUGl4sixPBqV9gHy+wVLqyGHiLbvgsjmkF8YVMbczjgQ8DH7zIlT2PwO+CPd0Emomqd99X98Uwcx1t1x5fXeIO3X4/GnBen3BT+NLpdQctO5sLVjFqXNzRwBYILO425Nuaz+0ZMsShSZu86UZbTXxQ6CRBfIloicBRlj/a2LCSGfgm2nHh+ekZnJUunzAsmVzBRQoFa0bWCOPyrqmleWeoiv1GoVyJPv5BUlt7H6r+ybyX9DAAA///YqAyuMQIAAA==")
	data["docstring/is-nil.md"] = decompress("H4sIAAAAAAAA/3yQP4/TQBDF+/0U7y5FLuKwxJ8mFaKgQKJMh5A1scfZEeNZa3ecwLdHuw5IIHTt7syb3+/t8GSiHzClPL86wLl4wS2yR87wyFhyusrIY5sooMww0fB5AqneH/lKupIzPNXP57po8CgF02qDSzLcRBWZfc2G3vPKfYdTZEySi8OSvTbR31NlVYdYA/iTsC2LXdBPpIX7LoTdDh8Nn37QvCiHgLvN/ukN3uId3h8qzyGEU2XhbexvlC0KZx5oLbxdbEhVDVJAUCnehfBFvjPmVBxL5lEGci7P/1gOZCAtCWeG8YWcR5x/1oWFbazs9cD+YY8hUqbBOdcapGBmsgKP5K3YNP2v/XltXXmLF+2a8UNTfoysmh7x9Xj89qLw1n34FQAA//8j90f0+QEAAA==")
	data["docstring/is-promise.md"] = decompress("H4sIAAAAAAAA/4xRvW7yQBDs7ykGKMD6+CxBHiBKkSJSSrooshZ7zZ1yvrNu15C8fXQ2WApKkXZ3Z3Z+Vtj0KXZO+BFtTN2/AsqigotltZygltGneHYNN+OFgNI4yiAxLy3I++uGz+QHUoZG0O1mmzkC1DpBO4RaXQy4OO+RWIcUUGkauCpxsIzWJVGEGP5f0bdLGbzChVHPzDIRuHBC1ZIXrkpjVis8BTx/Utd7NgbYNNyi381Gi2Ie7u+GcxT9Li/fdtjj4b0w5pDF88T5U/v0F0euaRAe5al1acoKTkA4c60xlca8ug9GF0XRJ25cTcqyvQumpgDyEnFkBD6RcoPjVwb0HJpsNb9YL9aoLSWqlVNOzgk6piBQSzo2EtvfuuuGMV7N9LcOyzGlxWx+adn7uPyT/ak68x0AAP//p+aZ8UcCAAA=")
	data["docstring/is-seq.md"] = decompress("H4sIAAAAAAAA/4yQzU4CQRCE7/MUBRyAiCT+y8l48GDikZsxm2a3l5k4O7NM94K+vZnlJ9Fw8NzVle+rESbCmyfUMTUXUyiLCnaW1XKCWkab4tZVXPUJASWG8KbjULKY1xrk/eHEW/IdKUMjCFvyrjpFZ7krQK0T1F0o1cWAnfMeibVLAYWmjos5lpZRuySKEMPl8f0Ylc4rXOjBTjX7BhfWKGrywsXcmNEIzwEvX9S0no3BwXI8ucI1bnA7xfsd7vGAx4+pMcuMxfvwOSpj3twno4miaBNXriRlmf3RKSmAvESsGIHXpFxh9Z0fWg5V5svY48EYpaVEpXLKvk7QMAWBWtJ+zlifW77p+lE01/8eV+a94aBXHFr2Pg6xWPxH7CcAAP//82gtTf8BAAA=")
	data["docstring/is-str.md"] = decompress("H4sIAAAAAAAA/3yRQU/rMBCE7/4V0/bQRu+9SA+49IQ4cEDi2BtC0TbZ1BaOHXk3Lfx7ZKdFRUKcPTOeb3aFjWi6Rx/T8KeCsqjgZFktJ6hljCkeXcddUQgoMUSTCwcxTz3I+/MDH8lPpAyNF8HfHBCg1gn6KbTqYsDJeY/EOqWARtPETY2dZfQuiSLE8G92X4QyeYULpctXyOzPqqYnL9zUxqxWeAh4fKdh9GwMzmDrzX/c4BZ3FZaWvY/LyphdrsSz9HujOQ57bmkSnn8tzTIlnIDgnWhtzLN7YwxRFGPizrWkXIivYVsKIC8Re0bgAyl32H9kw8ihy/3zB+vFGq2lRK1yyms4wcAUBGpJy8ax/+kYw1Qm0xx/3rwu5Itr9Aov2+3rr9TzHcxnAAAA//+lOibUDQIAAA==")
	data["docstring/is-vector.md"] = decompress("H4sIAAAAAAAA/5SRvW7jMBCEez7F2C5s43wC7v9cHa644oCU7oxAWFMrkwhFCtyVnbx9QMnOH9Kk5sxwvtkFVie2mvIftCl3n9ZQFhWcHavjDHWMPqeTb7gZFQLKjMkj5n8LCuHywCcKAylDE+gi2ZSECHVe0A7Rqk8RZx8CMuuQI2rNA9cVdo7R+iyKmOLnyXwVyhAUPo5lnkImv49H1C0F4boyZrHA34h/99T1gY3BM9xy9QVf8Q3f19j/wE/8wu/btTG70osn/etaUyYObGkQnr4e6xVWeAEheNHKmBt/x+iSKPrMjbekLJs3xJYiKEjCgRH5SMoNDg/F0HNsCkT5YDlbwjrKZJVzmcQLOqYoUEc6Lp3a907SDeNuWuIvp6lG/NmVf+44hDTHfrv9MLewTbF5AT6FVuYxAAD////9lvg9AgAA")
	data["docstring/last.md"] = decompress("H4sIAAAAAAAA/3SPsU7DMBRFd3/FlTqQSKgDRKidEAMDO3vjONfkiRc7iR0a/h7VIFAH1nd1jt7ZoVKbMhLnGgvzuoSEPBDlSuXIkBF9uSXOK4OjeR0kwa/BZYkBZ1H9Yf9HJzrxwv5Xcou4IIhCruWQBI5T/tzjJeMcV+3RER0DvTixihzhBrp3+LjA4sOq/GmxJglvaKvE+bFktejo40I4q1q2y4Pt3pjdDk8Bz5sdJ6UxQNXTY8NNdTziocH9HZpDXV+G0rTVxqC085u5Sj81h9PefAUAAP//Q8dqPVEBAAA=")
	data["docstring/lazy-seq.md"] = decompress("H4sIAAAAAAAA/0yPsU7EMBBEe3/FSGnsQ4eC6Ogo+AO6iGLtrBULx87ZG5S7r0f26QTTzr6Z2QE60u16rnyBz2U9nQy2kufdcQWh8mXn5BiykCBU8A/FnYRnRLqFeFVqGPCe8HHQukV+UwrQM/sEH2xPnb4UAOjIgskH2wpXOxMmgr17//U3R7uc2gbdIAv9BII1TQ+qOyNejDG91+dyZnILpgNa6JvxOo79qgU+OL2VkCQmHA37XNpX9/XoToUsDB9Klc77YHMi5wLSvlou9Vn9BgAA//8aeNH1NgEAAA==")
	data["docstring/len.md"] = decompress("H4sIAAAAAAAA/1SPwWrzMBCE73qKAR/+GILhb/MCPfTQe6HHIKvjSCDvOtKKNH36YjtQelx2vm93OhwyBZXXHoXWilRYJDLlYhE6bVPltVEC3UfK+ZHbFtLmkWWNMXOmWEUSeARtYn7Mv+iAt78upIrsv1O+I+i8NOPnEb7eJcSioq0eoQVqkeWWKpEk+GVT6oSRSS77lRWzmCqmJsGSCm7bk36FvIClaBmc6zq8CF6//LxkOoe9+L/DfzzhGae+d+591XBPPCx71fPpPLifAAAA///Twvx/KwEAAA==")
	data["docstring/let.md"] = decompress("H4sIAAAAAAAA/1yQP2vDMBDFd32KBxkck2Lo/w6l0KFDl06FDiEQRb7UAkXn+s6x3U9fLDuhrRah073fvXsLLAMp1o/RHgh7bg5Pq026Vzl2PpaCwM4GHG1oScyHDwGuIasEi0jd/C2Oa7pICh8/oRVBanJ+76mctVCGVlanXuwGjDMLvCq6kaoVRdDYO8L/EkZDgs5r5eNviI0lGtK2iUnQkLRBwfv0Clb0BPQcC2MWCzxHvPT2UAcyBvPyPbLlJa5wjZvcYD4D1re4wz0eNptUXDqOzip6DHluzHvlBTShpgXOsQQvCmdDoBJZnyWXFkdyys25PmT/3Qf7PWAaQjFZnjZhIQh9tRQdSYE3TvFYTSLtOOUoKBmRFdSP07lV8SWdotgG0m1KsTA/AQAA//+1iXTr8wEAAA==")
	data["docstring/list.md"] = decompress("H4sIAAAAAAAA/0zMO2rFMBCF4V6rONjN9SVkDymyhtQiOkaCkWRG49fug1+Q9p8zX4+XpGYYq+b3e8Cv0hsbPApXHCf3k0Tu/i9jjbURFGYWa/BKWCS4eJm9MZxkw6R1SYHhA1WhtFnLtcuT7ZeURpR6zw/mefl0ru/xVfC9+TwJnQNegSM2dJEitRuesqOzSOVVTnXDPri/AAAA///Sf+9X3wAAAA==")
	data["docstring/map.md"] = decompress("H4sIAAAAAAAA/1SPwWrDMBBE7/qKAV+sFAJ1Q+m1lP5BbyGEJVnXgvVK1cqN3a8vUcDg67xl9k2DdqSEftILjH+ePIT+giwYKdk9mVgvbO4jMxU20J0vK8BtiMZg4ZG1GCgzysDIbJMUxB6UkixBv2uccvwNV77WfyVERYkVrAWx3x6uBnvnmgbvis+ZxiTsHB7qba84zie0ux1mdN7j+IwOLzicvHNfQzDcgggylylrbd9OOLcdDnjFmz/v3X8AAAD//4F0fMMQAQAA")
	data["docstring/meta.md"] = decompress("H4sIAAAAAAAA/2yPsU7EMBBE+/2KUdIkp4iCMh0FBT0dQqe9eEMsOetTvOaCEP+OfE4HnTVvxs9u0a1ijDlua49NLG+aYIugxI4rQi3Qy3xH5QyfwKrR2MQNJda/M1vYavvmQ8BFDoM4cJmDU4qTZ/OfgmRbnixv8kDUtnhSPO+8XoMQAZ2TWeFivgTB2/6O7nTCjse+L/D+hQp7otfFpyqstuPNIcSb14//nSPRGfge7esqaOask/mozYBReZXj7gGjixOa5udMvwEAAP//4WKBkTkBAAA=")
	data["docstring/not.md"] = decompress("H4sIAAAAAAAA/2yOQUpDQQyG9znFT9+iFqR3cOENBJfDtM1zAumkJJlXe3vxKVih2yTfl2/CU7fEbH7eQe1DjlX1BukLewayMdJHNukcAZvXycVtkROfVozeRRXOObyjzFWDC+TBISSwXWW37TMsG/tVgnG959MHlz3RNOGl4/Wzni/KRPjJ3DRWtc2O6K1J/CO/3x3MlOtfxYGPdQSvu6XqYJRfQ7mP2dNXAAAA//9/HD/YBgEAAA==")
	data["docstring/ns-put.md"] = decompress("H4sIAAAAAAAA/0yPMW7rMBBEe55iYDVWYd/Bxb/BB1JaK3ElLUAuCXKlJLcPRDuOWw7fzrwOZ62XvBm0Qiky5lRij1HUV1B7qpkmBquVb/chIbQQhJ3CxrD0/AZR2MqomSeZhf0ffMXtiYkuFVQYEuNmNAYGqUfhugU7LpCCS0kFI4suKCSVPWQ+AjLjmA1SEcm36sKXtsZWjlfnug43xb8vijmwc3jZnbViSjrL0oN1xymX5LfJJOmpd+7/KhX8wPD5cmw6Vo4h93fiflQfYfMeWPfh1354tAxv8u4nAAD//+R4x/5lAQAA")
	data["docstring/ns.md"] = decompress("H4sIAAAAAAAA/0zPQU7EMAwF0H1O8aVuWmmGKyAWLDgBC8Qibf5ApNSJbA9Tbo8IQ9Xtt/2+PGAUg8SVjxOUflUxxB5YiwvDay7lPjjmWKo4N8ccjQlVkP3PecDLBf7Jw26qNEh1xKKM6RvcsvkJ2XH75WdiUUZnwsxLVWJmlo97LVMnpXYR2dC0fuXEdOo9y1WV4oe+f3S/D2EY8CR43uLaCkMAxkLHm/T3ra48C2/nnZjeAwCMTQUyTeEnAAD//wJiTz4nAQAA")
	data["docstring/nth.md"] = decompress("H4sIAAAAAAAA/3SSMY/bMAyFd/2Kh8vQBLgGaJvOxQ0dgm5t9zvaomMCMuWKku/SX19IdjL1RvE9kY+ftMNe84ghpgmint/geaAS8rcDEuckvLCBsFAojO66mtxPziWpIY+8SXmkjJ4UHWOIRT0oN9lm7mUQ9lv/OECywfhPYe35iPPQfKsqhlhyNdVaVxvZ7XS78lhPCpY8cmrKlnmLIlajl6TsERNIwSnF1Ookxv6IH8wzRDGJ+jV6G19XTYy/nOLHjqrTud0OT4rvbzTNgZ0D9p4HBHzYf8JnfMHpcKjFSjHghIfXFPXycHDu9ygGXu9tgQzPm/6Mjnsqdtv7hH1d5Gse18Khpu34Gls+Bqu/Y7gDDWJ5TXhu6T1+bYjwNM9BesoS1bmzWmZqDYqJXlqbF83jC4aifTU9bgD8nbLBSj+CrI0xkHos3OeYrL0zBYsYaWFQupSJtXrqVPbwkrjP4Yoc66zp+D63/zB7lRBwFQ7r6kZT5Wf1gWn9cnPiRWK54z26fwEAAP//yi3PN8gCAAA=")
	data["docstring/or.md"] = decompress("H4sIAAAAAAAA/3SPsW7rMAxFd33FRTI4ycsz0Bbo0C1D5o4dBcGhIwGyaFB0nP59IdtB3KGaBOJensMtdixoWbrDYY+epHwzHLJn0f9NkGYIGtIVFUtlzjcXB6eUoZ4wZ1vhDpFahTIkXL3WOGVk5gSXwYlAzxrDoVIZ1H9XKGM6YgwxQkgHSVDvdJ7X+FRPMoZMCDqHeuGG6FLWLDuLWnFJdNdJqDZmu8Up4Xx3XR/JGEw37v7hBa942xtMr3Ux0/LfJNZFgC6bvTFfKyP7bo9IdCNZM+3Ut2CB/V23R4yehFz+eLDXrBTig/oXsdyTVSbMM2Rr8xMAAP//wi8tD64BAAA=")
	data["docstring/partial.md"] = decompress("H4sIAAAAAAAA/0yQy2oDMQxF9/6KC9lkSBsaKIUsu2g/oBS6TDSx4hF1ZONHpvn7Ms5za52je60Z5pFSEfLYV92Bklt0cKycqHAGteciQdFTZougIFwVitGfzBeXmnRilUd8XvlxCJkhKmc0uXpgLRmUGDHxcx+qWpSA0sCYwlEs2yV+BlaUgcp9l2SIHsMv2yeQnm7ww9pRvEeWQ/Qn9DxVY52IFvCYeFd63ofE2JH3oq5hIYkTJX+LXhozm+Fd8fFHh+jZGGBueY/oa1693M+3wCveum4aX0brzpjvQTL4rJ4rpnYtbFbrDShPn69uwHbefKy7LUa+lGK7NP8BAAD//+QMK6ahAQAA")
	data["docstring/partition.md"] = decompress("H4sIAAAAAAAA/0yPzU7DMBCE736KkXKxJaia0tJyQhx4AyQOCDVpsmlXctbGXvPz9ihBLT3Zmm/0aaeCjW1SVg6CLhRRZKX4iEwfDheU0U5JIenIvLL3/+iKgEUDjimUmBEGNLOwAXkaSTTfgKVL85/liMMP9ESQMh4oTf1zDz0NLNSDBc10TgMb0kXHwznlDAmKmMIn99S7hTFVhSfB83c7Rk/GAFbDrees10NXuMNbPT9rbHCPLXZ4QL18d86YlxNn0J8BX9PYRFqSYG9tjZWDXWPjYLfYOdh66dx+YX4DAAD//y+QZ9FJAQAA")
	data["docstring/promise.md"] = decompress("H4sIAAAAAAAA/1SQMW7DMAxFd53iI1mSIMgdMnTo2hswEh0TpUlBluz29oXgxkVHiR/vP/KIUy4+ycxn5OKpRZ5BhmaFZ9eFExbSxuGDays2g/CbhxjqyBi8TPABhKFZrOJ2w/sAylmFE1apo7famVSebWKrV9RR5j2OVVRxeajHz8sVK0kVe3YuaCtHdTwYiVUWLl1JCIRIqn3ULV5SO7SOVCEWtaVto1f7Dfc9Hcngpt//6W6RbyEcj7gb3r5oysohAKfEA/Lfwc79L+Mwsqofttc5/AQAAP//etNsWFIBAAA=")
	data["docstring/quote.md"] = decompress("H4sIAAAAAAAA/0zPQU7DMBCF4b1P8aQu2kgoUoELsGDBgh0H6DSZkJHssfFM2ub2KC4gtrbs9387HL6W7Iwp19Shsi9VDT4zrPAgk/DY7iCKkZyQ8sjhnUlFP+EzOaKYG0hH2JrOORquEiM0O84MvlBcyHns8TGLIdFQM8TA08SDy4Xjet+jxCBDqVxYx+17UvCtVDaTrLiKz9sRlWxec5kZh33Xh7Db4UXxeqNUIoeAX9PhiEc84bnrQmjbretubJNRnCvFJkAln7luJIXXtfEyqJSfPl3SmSuOD1ukTBDHlSuDMC06uGTtgTffaLZq1jXlxe7R2/N/kNP+r+zUh+8AAAD//+9CmPqDAQAA")
	data["docstring/range.md"] = decompress("H4sIAAAAAAAA/1yQQW6rQBBE93OKkrwBiY/AXxzA+vqLnCCLKII2NPYoMz1kprFwTh95EnvhZVWp66lrhyKSnBjeCjxtsDKWGCOTcgIhh+bfQzv6uiLx58oyMvRMiiVyYtEEPTNk9UeOCXMMPncWVka3JnvhEhoyouDtblU4Xm/IyJ5FaxycA8XT6nMhRUZY1AYhV2PwVgZMPNPqbrSAvukrDJ62J9vK3FcgmTAk5eUpbfsaLzOCuCuC8IMHm7DEcLETTxVs1pofn0Dph1Mbs9vhIPi/kV8cGwMAhYY/Fx41RBRKH4zuvmrbwMqMrixLY16tc4isa5S81e9J/9Y2aDvsG+w7/G3e+9p8BwAA//+GOQX8lwEAAA==")
	data["docstring/read.md"] = decompress("H4sIAAAAAAAA/xTKuw2AMBAD0J4pLNHALixxil2kSaL7CMZH6V7xTlwuIyL9xlbA7cXTY6FNCn3kRA19ywZF0NL2rpbliuMPAAD//w7hDINBAAAA")
	data["docstring/reduce.md"] = decompress("H4sIAAAAAAAA/1ySTY/TMBCG7/4Vr9QDrbZELFA+jiBx4L6cEFqZZNKM1rXDzDht/z2yk7SIU+TRTJ5nPjbYCnW5JfQ5tlD687DDHNHyyhRbUvfdSLyRIk0k8FAypP6esJ9rOB5hA7GAAp0omsJSSed4DAQhzcFK0uRDpgZPw8w1ThGjpIk76nDKajD/QrBzgpdjrr9a0lnU4GMHpTbVz+xwR1JsU45GQh28UDECRzb2YQYr/DgGpq7Y2eDtJlEZQr43kn0tnJ1Ls+U1Ck2csqL1oc3BV3FWZC0srTmL4aK9x3ngMEtEutiq+X/V0s1a1ji32eBLxLeLP42BnAO2HfW44Ocj3uId3uPwa7dGr3i1/YCP+ITPeHyzq/Flrw+44Lpz7mlgxZlDgJBliZVax4Hnw+G5ce5HDPxCSDaQgAu1mNQedR3AP8dS4mVGrEiRkGK41q0p2PQ28EC9LZhe0mnpdTmb284bfM0GrlMhr1dQTPk4lP0IjZIqtKJ+0+AnToJcjgo+3kjzkb1e5ztRa0ka9zcAAP//5r+52eICAAA=")
	data["docstring/repl-cls.md"] = decompress("H4sIAAAAAAAA/1JW0EjOKdZUSM5JTSwqVijJSFUoTi5KTc3jcgaJgAWSS4uKUvNKFJLz80pAdH4akjo9hZCMzGKFzGKFRIUg1wAf3fy8nEqFtNK85JLM/Dw9LkAAAAD//yVoAgJhAAAA")
	data["docstring/repl-doc.md"] = decompress("H4sIAAAAAAAA/1SOQQrCMBBF93OKD91Y0N5BsDsXIl5gSCZ0oJ0EJ4Xm9tKKC7fv8x+vwynmgJTfS4+oXmZujpjDuohVrpqNbj9cJ/mf9t9BvUjQpBIP0xmawNYgm3r1Aa9JHepgPMfH/ZJtbkirhd0xEHUdroZx46XMQoRvE3uz0NMnAAD//xQLeSyiAAAA")
	data["docstring/repl-help.md"] = decompress("H4sIAAAAAAAA/2SOS27jMBBE9zpFLWVh4APMbuDxbhaTIAcwJZUkAs0mwyb9uX0gOTFipLf1XnV17UJJuw73++stibuhLN6wJgg0czObrh3jgCnmsMFf4BiHGqjFFR91j+PVhST8jdOGz3F3arq2GqH2+eSwOJ2JoeZMLVAXaMkNfJJXY4p3exB77MNB6DLKQtiQSW269r368gBeqi9b/Hr8/69p3hYi5RhSQU+JF/iRWvzkaRv1Y8UvTFEkXjiiv22I1tAzI07gNWWa+air7AouXgSLOxM9qeDZSXWFI7w+ta9bYHd1jz8wr7PwWx8Gp7DkFKFK8UkI8UrbNx8BAAD//3pr7pSgAQAA")
	data["docstring/repl-quit.md"] = decompress("H4sIAAAAAAAA/1JW0CgszSzRVACRxQolGakKQa4BPlyeubmpKZmJJak5lQqpFTCp4ILSkpLUIrASPYWQjMxihcxihUQwXzc/L6dSIa00L7kkMz9PjwsQAAD//+kUzLVaAAAA")
	data["docstring/repl-use.md"] = decompress("H4sIAAAAAAAA/0TNQQrCQAxG4X1O8cNs7MLeQaQ7FyJeIITUGWiT0mSg3t6FgtvHg6/g1ENhMUAq20sDxqvGxqJ0/ZWsCun7rpYQt9Qj/9eIZ22BFmA8pvvt7La8MXeTbG4jUSm4GKaD121RInzF2X2gTwAAAP//xZr5jn4AAAA=")
	data["docstring/rest.md"] = decompress("H4sIAAAAAAAA/1SQvU7DMBSFdz/FkTqQLBmgQnRCDAzs7LXjHBOLW7vxD83jo7ooFZvlc7+j+90dusRckLn0SCw1hYwyE+03uvbOXCqDpfqcfYarwRYfAy5e5I+B2YZQZlPA1UqdeOtyPuUCCk8M99IzrXee00YO+Ci4xCoTRmJkoPPWG0GJsDPtN1xMMPgx4u8UavbhC7rLXF6bh8ZIFxNhjUjL2gIaMUFfvfSg1G6Ht4D31ZzOQqWAbqLDiofucMDzHk+P2L/0/TVop1h7pdD8eWP+6R+7DTkO6jcAAP//LTZxf1UBAAA=")
	data["docstring/seq.md"] = decompress("H4sIAAAAAAAA/2yPwWqEMBiE73mKIXtoAtV36KGH3gs9iEhWRwzExE3+de3bF5UtPfQyEPLNfPwXmMIbxpRnCyfCeZECSehTXJnl+NnfDoW3O2NP9eVDeLJ/UZmIJafVDxz+6cGP8AJf4ovAhUw3fNf4GI/egfcuxiS48jnJ4RVd9KHDY3deiUy558ihVupywVvE++bmJVApwDwmxipQ0JTzLD0xhKRtqwDASKpW9pIyzOwWmDGi2VqYIhkbdKWtRbHWKvU5+QKe06f7FKNr9FRpaO4RfiNVuu1q9RMAAP//+bJX0E4BAAA=")
	data["docstring/str.md"] = decompress("H4sIAAAAAAAA/5yPS2rzMBSF51rFSTL5Y34CfWwglG6g6ayURLaPa4Esmasrp9l9ieyGQmcd6ug87rfBv6SCLspQVVs0MUwUTUVI0AiLpOLCh3kSWmWCReB5EdFJHKA9l7frHFtM1mcmxK78jBIn17KdK3fGbDbYBzx/2mH0NAbzBeue3sc13u5wjwc8vm+Nee1dAmcjzs57CDVL+LGIY8ndUsd54IW2peBQPMmYvVdKsOom+st/xEA0NiAnlq5TUlmd0OXQqIvhyj1KbHPDG/+CRknF0FtFNdhLhZqwtec1UxNC26K+lNbDmFUpRaPsvlFXf2f9FTzuzFcAAAD//9gkDpPCAQAA")
	data["docstring/sym.md"] = decompress("H4sIAAAAAAAA/3yOQU7DQAxF9z7FV7OgWZA7AOIG7CuTcZmRHE/qcQq9PZoAEquu7K/3bb0Bx3Zb0MJHzNWu4tHAPRf7QLGoYOtT3CSh3Zb3qvTy14wsWL1eS+rw7tGE5xoZl421nIsksCVonVl/Cw3sAp5nWUPSRDQMeDK8fvGyqhABxyRnZFGtj116Nz/s+TCOncvlH37Y15HoLZcG+XmDz6IKl9jccArf5DTRdwAAAP//mE7l2gUBAAA=")
	data["docstring/take.md"] = decompress("H4sIAAAAAAAA/1SQvU4DMRCEez/FSClyJ1AkfhJoKSjokSgQipy7uZzFxQ7rNUl4emRffqCdnZ1vdyao1H4STUheEfl1VSMLEdoTnZOo4MANvUaELjsSfcNo3twwQKhJPCwG+3M4D7ORTnsKpiV4iiDouKNcwjoJmwJZce28d36d17KwlfDtWrYXGKom+MYqvVW29QwvozOGJA3/+KwQsQ+iFGhvfXFJHkdlO355fRRjGjRTz1fv8kcrQiUVWDszZjLBk8fz3m62A40BqpYd9phWN7jFHe7r+iQe8D7HAg94/Cha6XWBPQ61Ma+9i+AYM4KO1eVb/pe3PEVjjkW9nJnfAAAA//9Y1ktVpAEAAA==")
	data["docstring/thread-first.md"] = decompress("H4sIAAAAAAAA/1yRz86bMBDE7zzFSDkE8lGSRs21Ug499B6px7CBDbZk1she8qdPX9nQSv1utjUzu7/xBuWX73iQmxl3H8a421VQE5j6uD6rCX4eDDpyLoIi1LANuNsQFRSGeWTR4rJ6ttm0/edSw4jzNDnL/TKhwVne+Qg1pLAR4hXkUsAbJKAk70itFzytc7gxOi8PDso9rKiHF25wMbzumFVqWJLUSlyUFGE1ftoUJP0ygfs6X2KKa4pis8FZ8ONF4+S4KJCrOaD8wNdDhXKHY4Vyj1NVFcWvNJBfU/KrR1vuk6D8wCGLjxVOVZvTA+scBNdv1wY/BT70HJLlxqoc8LBxJmd/M56GdBsxeCsDvNSJEaMdjKIz3kdOroUNhLZuEx9hctSx8S7F3n3IhS8fyP3STvMfS/2Xpl556kz0JwAA//8UupbuCQIAAA==")
	data["docstring/thread-last.md"] = decompress("H4sIAAAAAAAA/1yRMW/cMAyFd/2KB9xw1tV13QBZA2To0D1Ax5ixeZYAmTIk+pL01xeSDhk6SSK+98hHndB9f3rCjcLBuMa05cvFQl1iWvK9rC7FY3WYKYQMylDHPiFQVlBaj41Fzctdcq6a85dIHSMf+x48L63BgGf5rFeoI4XPkKigUAw+QQIq+Ezqo+Ddh4A3xhzlxkl5gReNiMIDXhzfR6yUOpaCesmNpAyvGVcvFL4mBcnSOvDS10cudoMxpxOeBb8+aNsDG4O2mhHdN/wcLboLHiy6H3i01pg/pSN/7MVAI6ZSr0jDMVprp2qfWI8keB2Hh8fXAb8FMS2ciuqNVTnh5vNBwf9lvDvSc8YavayI0pec2PzqFLOLMXNRtXwgTP1UMhL2QDO7GIrtNaa69PaHvLQNDf/nQd8S1bPM3ltr/gUAAP//Adj/ig4CAAA=")
	data["docstring/to-assoc.md"] = decompress("H4sIAAAAAAAA/2SQzWrrQAyF9/MUB2dxY24aCt1510UfoXRRClFn5Fh0ftwZTWzn6YtdGgpdCaHzfRLaYa/pjkpJFoU//7ewKV44a1nbytFygSZQxBYSUrkwiuZqtWY2L+L9ylhSjqQMQmFF6n/xEv8YjniOXj4YOqw2io6yw+lbdEJfo1VJ8YBpEDtACjxdxS+wKYxV2R02MnOpftumg5Qbhmm96p0RSDkLebmyg4TATkjZL0djdjs8RjzNFEbPxgB7xz1mvHaRAqMpY1Xl3KCjM+P++PDW/oQW/Nt3E8t5UDR+LU27DW+vnLG05isAAP//+ZF5XlwBAAA=")
	data["docstring/to-list.md"] = decompress("H4sIAAAAAAAA/0yQTU7DQAyF93OKJ2XBjFoqAb0AC46AWFSVamYcxWJ+QsahSU+PEiBi6efvsy03sFruo1RF5c+dgy/5iwetSzly9lyhBYQFMW8S40J4Us6kDEJlRWn/0ZI3/oDXHOWDoR2jKuVAQ8DlZ8AF7Zi9Ssl7XDvxHaQi0k3iDF9SPyqH/WoOXMe4btFO6qbhulzzzkikPAhFuXGApMRBSDnOB2OaBs8ZLxOlPrIxgA3cYsLpAY94wvHs/rLZAIBN1MO2GafpDLvDhKNza+fO/ipuDbavTZid+Q4AAP//kawXNEYBAAA=")
	data["docstring/to-vector.md"] = decompress("H4sIAAAAAAAA/0yQwU7DMAyG73mKX9qBdIxJDF6AA4+AOCCkmcRVLRKnNO7W7elROwl6tP/vsy1v4K08nDhYGVD5575BKHriwepcjqyBK6yAcIPcu6Q0M4GMlYxBqGwo7YoXXRl7vGmSb4Z1jGqkkYaI423EEe2owaToDudOQgepSHSVdEEouR+N424xB65jWvZYJ/VPw3m+54uRyXgQSnLlCMmZo5Bxuuyd22zwonidKPeJnQN85BaTAwCfqYdvFR/TJ/x2iwmHplmiO/+IA57w3CyN1aP+E0yN+w0AAP//fvP0f0QBAAA=")
	data["docstring/vector.md"] = decompress("H4sIAAAAAAAA/1SOQYrrMBBE9zpFkWziwM8d/mJuMDDrHqmEBHLLqDuOffsBx7OYbb3i8a64rYzeB3If8/0+IQ6K0yBQvvCG4au2dpI/AK/SjWDjTHWDDMILwVXaU5zp0BqW0deamB74LNWQnxq9dkU1aEeqOXNQHV5ED8Gpb9U5pMF2ddnALXI5Xo7qiKL4JvzoSqgKQa7D/F9sYoYsVmrXRwjXK/4rPjaZl8YQgFtixoZLYWv9Mv0uOy5eOPhezogN+xR+AgAA//+vfXdiKQEAAA==")
	data["docstring/when.md"] = decompress("H4sIAAAAAAAA/1SQTU7DMBCF9z7FU7sgqUolKH8rJBYsOAHbOvGkGTEdV7ZDk9sjO6WUlW1974013xLVqSfFMZBD58NhtarRenWc2KsVmUDfVgabKMKiEd9+mY8OqacLcKXNrU0EjkhhSP2ESn1CZyXSGvmqLPW69PI3ETZcTdjgk0Uub5Bt+5IDK9IQdI1A+WTdzzNY7SXPXmEjOEUEioOkjTHLJd4U76M9HIWMASpHHUbcVHe4xxYPeMQTnvFS14UWC9UrKiHFWGNbGwCojoE1iWIx5t0a3i/OYM7V5tdGnA6NF+zyoFv1aZfzQyQH1pjIOvhupruiQUvrn7k/oVZdwcU3TtlNcy3cq0zg7rxu7mbFs3nzEwAA///KLwLw1AEAAA==")
	data["docstring/with-meta.md"] = decompress("H4sIAAAAAAAA/1yQsU4zMRCEez/FKGly+fOnoKSjoKCnQxQbex2vdLc+btdJJMS7o7tEBNGOPdrvmzU2Z/Hyf2An5DoNGGgcOf3rQO4UCxvmt0RO8Irrp/CS4YVhI0fJwulaFUOkkQ49o2YcWPQIUq1OzmkHL2LITaNLVZyl7zGxt0lBUD5D1Jw0LmUv5DhR3xgz359rP0Q3xrTHa+F7LIZxqidJnCC6lBfAmm96MP5orJFtB2uxgAxkVqOQy4lhPrXobWLbh7Be40nxfKFh7DkE/JosAMAmK94u79hst7jgoeuW9PMx1Qis6K68SKXaDj3bvIYebfXVhe8AAAD//z7LZmuEAQAA")
	data["docstring/with-ns.md"] = decompress("H4sIAAAAAAAA/3SQsW7yQBCEez/FyDSg/wcpKelSUKTPCyy+MZx03rNu18G8fcQp2EmR9ttvR6PZYHuLft2rQWUg+lyGfzvwU9IkTqvAIBeJag6plo3SsTktjsDoyP23HBV+JbqszrlygY3sYh8Z1oAD3vsqLgQh06DZIalQwh2co/l/RMctpoQz0RWKM+DMPhc+i0a91Kha4ICP63JhQKFNqfZ4KEnMq7dEFvpUlOHQNJsN3hSnWYYxsWmAbWCPGW2e3GLg777t7mE89xvu+3UcYP2N+scrsB2Loh1LVMfLES3mXc38gV+PD9p8BQAA//9UbL2/qQEAAA==")
}
