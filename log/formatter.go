package log

import (
	"fmt"
	"strings"
	"time"
)

func formatArgs(useColors bool, args ...any) string {
	if len(args) == 0 {
		return ""
	}
	if len(args)%2 != 0 {
		panic("odd args length")
	}

	argsFormatted := make([]string, len(args)/2)
	argColor := "%s="
	if useColors {
		argColor = ARG_COLOR + argColor + RESET_COLOR
	}

	var newLine string

	for i := 0; i+1 < len(args); i += 2 {
		name, value := args[i], args[i+1]
		if _, ok := name.(string); !ok {
			panic("name field is not string")
		}

		valueString := fmt.Sprintf("%v", value)
		if strings.Contains(valueString, "\n") {
			lines := strings.Split(valueString, "\n")
			var b strings.Builder

			for _, s := range lines {
				if useColors {
					b.WriteString(fmt.Sprintf("  \u2502 %s\n", s))
				} else {
					b.WriteString(fmt.Sprintf("  | %s\n", s))
				}
			}

			argsFormatted[i/2] = fmt.Sprintf("\n  %s\n%s",
				fmt.Sprintf(argColor, name), b.String())
			if newLine == "" {
				newLine = "  "
			}
			continue
		}

		valueString = strings.ReplaceAll(valueString, "\"", "\\\"")
		if strings.Contains(valueString, " ") {
			argsFormatted[i/2] = fmt.Sprintf("%s%s\"%s\" ",
				newLine, fmt.Sprintf(argColor, name), valueString)
		} else {
			argsFormatted[i/2] = fmt.Sprintf("%s%s%s ",
				newLine, fmt.Sprintf(argColor, name), valueString)
		}
	}

	return strings.Join(argsFormatted, "")
}

func formatMsg(useColors bool, level Level, msg any, args ...any) string {
	lvlString := level.Name + strings.Repeat(" ", LEVEL_NAME_FILL-len(level.Name))
	timeString := time.Now().Format("02 January 2006, 15:04:05")

	var writeMsg string
	if useColors {
		writeMsg = fmt.Sprintf("%s %s %s %s %s",
			timeString, level.Color, lvlString, RESET_COLOR, msg)
	} else {
		writeMsg = fmt.Sprintf("%s [%s] %s", timeString, lvlString, msg)
	}

	if argsFormatted := formatArgs(useColors, args...); argsFormatted != "" {
		writeMsg += " " + argsFormatted
	}

	writeMsg += "\n"
	return writeMsg
}
