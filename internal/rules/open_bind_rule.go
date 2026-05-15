package rules

type OpenBindRule struct{}

var hostKeys = []string{"host", "address", "addr", "bind", "listen"}

func (r OpenBindRule) Check(config Config) []Issue {
	for _, key := range hostKeys {
		found, values := config.FindByKey(key)
		if !found {
			continue
		}
		for _, v := range values {
			s, ok := v.(string)
			if ok && s == "0.0.0.0" {
				return []Issue{{
					Severity:       MEDIUM,
					Description:    "Сервис слушает на всех интерфейсах (0.0.0.0)",
					Recommendation: "Укажите конкретный интерфейс или ограничьте доступ фаерволом",
				}}
			}
		}
	}
	return nil
}
