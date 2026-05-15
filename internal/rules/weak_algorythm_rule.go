package rules

import (
	"fmt"
	"strings"
)

type WeakAlgorithmRule struct{}

var weakAlgorithms = map[string]string{
	"md5":  "MD5 — криптографически сломан",
	"sha1": "SHA-1 — считается ненадёжным с 2017 года",
	"des":  "DES — устаревший, длина ключа 56 бит недостаточна",
	"rc4":  "RC4 — множество известных уязвимостей",
	"3des": "3DES — устаревший, уязвим к SWEET32",
}

var algorithmKeys = []string{"algorithm", "cipher", "digest", "hash"}

func (r WeakAlgorithmRule) Check(config Config) []Issue {
	for _, key := range algorithmKeys {
		found, values := config.FindByKey(key)
		if !found {
			continue
		}
		for _, v := range values {
			s, ok := v.(string)
			if !ok {
				continue
			}
			lower := strings.ToLower(s)
			if reason, weak := weakAlgorithms[lower]; weak {
				return []Issue{{
					Severity:       HIGH,
					Description:    fmt.Sprintf("Небезопасный алгоритм: %s. %s", s, reason),
					Recommendation: "Используйте SHA-256, AES-256-GCM или ChaCha20-Poly1305",
				}}
			}
		}
	}
	return nil
}
