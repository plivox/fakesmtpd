package config

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(config *LogConfig) zerolog.Logger {
	var writers []io.Writer

	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to prase log level")
	}
	zerolog.SetGlobalLevel(level)

	if config.Console != nil {
		if config.Console.Json {
			writers = append(writers, os.Stdout)
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				NoColor:    !config.Console.Color,
				TimeFormat: time.RFC3339,
			})
		}
	}

	if config.File != nil {
		file, err := os.OpenFile(config.File.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		// Fix: defer file.Close()

		if config.File.Json {
			writers = append(writers, file)
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        file,
				NoColor:    true,
				TimeFormat: time.RFC3339,
			})
		}
	}

	// 	if !reflect.DeepEqual(zerolog.ConsoleWriter{}, writer) {
	// 		writer.FormatLevel = func(i interface{}) string {
	// 			return strings.ToUpper(fmt.Sprintf("| %s |", i))
	// 		}
	// 		// writer.FormatMessage = func(i interface{}) string {
	// 		// 	return fmt.Sprintf("%s | %s", guid.String(), i)
	// 		// }
	// 		writers = append(writers, writer)
	// 	}

	return zerolog.New(io.MultiWriter(writers...))
}
