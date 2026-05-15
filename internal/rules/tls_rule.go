package rules

type TLSRule struct{}

func (r TLSRule) Check(config Config) []Issue {
	var issues []Issue

	tlsKeyfound, tlsValues := config.FindByKey("tls")

	if tlsKeyfound {
		tlsIsOff := false

		for _, tlsvalue := range tlsValues {

			tlsString, isString := tlsvalue.(string)

			if isString && tlsString == "enabled" {
				tlsIsOff = true
			}

			found, values := config.FindByKey("enabled")
			if found {
				for _, v := range values {
					if b, ok := v.(bool); ok && !b {
						tlsIsOff = true
						break
					}
				}
			}
		}

		if tlsIsOff {
			issues = append(issues, Issue{
				Severity:       HIGH,
				Description:    "TLS отключён",
				Recommendation: "Включите TLS для шифрования трафика",
			})
		}
	}

	// случай 2: verify/insecure отключена проверка сертификата
	for _, key := range []string{"verify", "tls_verify", "ssl_verify"} {
		found, values := config.FindByKey(key)
		if !found {
			continue
		}
		for _, v := range values {
			if b, ok := v.(bool); ok && !b {
				issues = append(issues, Issue{
					Severity:       HIGH,
					Description:    "Проверка TLS-сертификата отключена",
					Recommendation: "Никогда не отключайте проверку сертификата в продакшене",
				})
				break
			}
		}
	}

	found, values := config.FindByKey("insecure")
	if found {
		for _, v := range values {
			if b, ok := v.(bool); ok && b {
				issues = append(issues, Issue{
					Severity:       HIGH,
					Description:    "Установлен флаг insecure",
					Recommendation: "Уберите insecure: true и настройте корректные сертификаты",
				})
				break
			}
		}
	}

	return issues
}
