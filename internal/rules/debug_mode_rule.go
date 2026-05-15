package rules

import (
	"strings"
)

type DebugModeRule struct{}

func (rule DebugModeRule) Check(config Config) []Issue {
	var issues = make([]Issue, 0)

	ok, logInfos := config.FindByKey("log")

	if !ok {
		return nil
	}
	for _, logInfo := range logInfos {
		level, ok := logInfo.(string)

		if ok {
			ok, issue := tryGetDebugIssue(level)

			if ok {
				issues = append(issues, *issue)
			}
		} else {
			child, ok := logInfo.(map[string]any)

			if ok {
				childConf := Config(child)
				ok, levels := childConf.FindByKey("level")

				if ok {
					for _, level := range levels {
						levelStr, ok := level.(string)
						if ok {
							ok, issue := tryGetDebugIssue(levelStr)

							if ok {
								issues = append(issues, *issue)
							}
						}
					}
				}
			}
		}
	}

	return issues

}

func tryGetDebugIssue(level string) (bool, *Issue) {
	if strings.EqualFold(level, "debug") {
		return true, &Issue{

			Severity:       LOW,
			Description:    "Логирование в debug-режиме",
			Recommendation: "Поменяйте режим на более избирательный (info+)",
		}
	}
	return false, nil
}
