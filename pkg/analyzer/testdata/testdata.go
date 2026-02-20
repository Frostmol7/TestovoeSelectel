package testdata

import "log/slog"

func _() {
    slog.Info("Starting server") // want "log message should start with lowercase letter"
    slog.Error("Failed to connect") // want "log message should start with lowercase letter"
    slog.Info("starting server")
    slog.Info("test ok")
    slog.Info("server started")
    slog.Info("started! @") // want "no special chars or emoji allowed"
    slog.Info("server started")
    slog.Info("user password is secret") // want "sensitive data detected"
    slog.Debug("api_key=value") // want "sensitive data detected"
    slog.Info("user authenticated")
}
