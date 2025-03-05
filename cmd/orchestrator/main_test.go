package main

import (
	"os"
	"testing"
)

// TestCheckEnvironmentVariable проверяет функцию checkEnvironmentVariable.
func TestCheckEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name        string
		envVarName  string
		envVarValue string
		setEnv      bool
		wantErr     bool
	}{
		// Успешные случаи
		{"Valid positive number", "TIME_TEST_MS", "1000", true, false},
		{"Valid large number", "TIME_TEST_MS", "5000", true, false},

		// Ошибки
		{"Variable not set", "TIME_TEST_MS", "", false, true},
		{"Non-numeric value", "TIME_TEST_MS", "abc", true, true},
		{"Negative value", "TIME_TEST_MS", "-100", true, true},
		{"Zero value", "TIME_TEST_MS", "0", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменную среды, если нужно.
			if tt.setEnv {
				os.Setenv(tt.envVarName, tt.envVarValue)
			} else {
				os.Unsetenv(tt.envVarName)
			}

			// Вызываем функцию.
			err := checkEnvironmentVariable(tt.envVarName)

			// Проверяем результат.
			if (err != nil) != tt.wantErr {
				t.Errorf("checkEnvironmentVariable() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Очищаем переменную среды после теста.
			os.Unsetenv(tt.envVarName)
		})
	}
}

// TestGetPortFromEnv проверяет функцию getPortFromEnv.
func TestGetPortFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envPort  string
		setEnv   bool
		wantPort string
		wantErr  bool
	}{
		// Успешные случаи
		{"Default port (not set)", "", false, defaultPort, false},
		{"Valid port", "8080", true, "8080", false},
		{"Valid port (different)", "9090", true, "9090", false},
		{"Negative port", "-80", true, "-80", false}, // Теперь это допустимо
		{"Zero port", "0", true, "0", false},         // Теперь это допустимо

		// Ошибки
		{"Non-numeric port", "abc", true, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменную среды, если нужно.
			if tt.setEnv {
				os.Setenv("PORT", tt.envPort)
			} else {
				os.Unsetenv("PORT")
			}

			// Вызываем функцию.
			gotPort, err := getPortFromEnv()

			// Проверяем результат.
			if (err != nil) != tt.wantErr {
				t.Errorf("getPortFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && gotPort != tt.wantPort {
				t.Errorf("getPortFromEnv() = %v, want %v", gotPort, tt.wantPort)
			}

			// Очищаем переменную среды после теста.
			os.Unsetenv("PORT")
		})
	}
}
