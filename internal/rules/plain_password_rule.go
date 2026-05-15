package rules

import "strings"

type PlainPasswordRule struct{}

var passwordKeys = []string{"password", "passwd", "secret", "pass"}

func (r PlainPasswordRule) Check(config Config) []Issue {
	for _, key := range passwordKeys {
		found, values := config.FindByKey(key)
		if !found {
			continue
		}
		for _, v := range values {
			s, ok := v.(string)
			if !ok || s == "" {
				continue
			}
			if looksHashed(s) {
				continue
			}
			return []Issue{{
				Severity:       HIGH,
				Description:    "Пароль хранится в открытом виде",
				Recommendation: "Используйте переменные окружения или менеджер секретов",
			}}
		}
	}
	return nil
}

func looksHashed(s string) bool {
	if strings.HasPrefix(s, "$2") {
		return true
	}
	if len(s) >= 32 && isHex(s) {
		return true
	}
	return false
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
