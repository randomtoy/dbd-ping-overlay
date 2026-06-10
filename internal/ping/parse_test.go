package ping

import "testing"

const englishPingOutput = `
Pinging 203.0.113.10 with 32 bytes of data:
Reply from 203.0.113.10: bytes=32 time=37ms TTL=55
Reply from 203.0.113.10: bytes=32 time=38ms TTL=55
Reply from 203.0.113.10: bytes=32 time=39ms TTL=55
Reply from 203.0.113.10: bytes=32 time=38ms TTL=55

Ping statistics for 203.0.113.10:
    Packets: Sent = 4, Received = 4, Lost = 0 (0% loss),
Approximate round trip times in milli-seconds:
    Minimum = 37ms, Maximum = 39ms, Average = 38ms
`

const russianPingOutputNoSpace = `
Обмен пакетами с 203.0.113.10 по 32 байт:
Ответ от 203.0.113.10: число байт=32 время=38мсек TTL=55
Ответ от 203.0.113.10: число байт=32 время=38мсек TTL=55
Ответ от 203.0.113.10: число байт=32 время=38мсек TTL=55
Ответ от 203.0.113.10: число байт=32 время=38мсек TTL=55

Статистика Ping для 203.0.113.10:
    Пакетов: отправлено = 4, получено = 4, потеряно = 0 (0% потерь),
Приблизительное время приема-передачи в мс:
    Минимальное = 38мсек, Максимальное = 38мсек, Среднее = 38мсек
`

const russianPingOutputWithSpace = `
Статистика Ping для 203.0.113.10:
    Пакетов: отправлено = 4, получено = 4, потеряно = 0 (0% потерь),
Приблизительное время приема-передачи в мс:
    Минимальное = 38 мсек, Максимальное = 38 мсек, Среднее = 38 мсек
`

const englishTimeoutOutput = `
Pinging 203.0.113.10 with 32 bytes of data:
Request timed out.
Request timed out.
Request timed out.
Request timed out.

Ping statistics for 203.0.113.10:
    Packets: Sent = 4, Received = 0, Lost = 4 (100% loss),
`

const russianTimeoutOutput = `
Обмен пакетами с 203.0.113.10 по 32 байт:
Превышен интервал ожидания для запроса.
Превышен интервал ожидания для запроса.
Превышен интервал ожидания для запроса.
Превышен интервал ожидания для запроса.

Статистика Ping для 203.0.113.10:
    Пакетов: отправлено = 4, получено = 0, потеряно = 4 (100% потерь),
`

func TestParsePingOutputEnglish(t *testing.T) {
	got := parsePingOutput(englishPingOutput)

	if !got.Available {
		t.Fatal("Available = false, want true")
	}
	if got.AverageMs != 38 {
		t.Errorf("AverageMs = %d, want 38", got.AverageMs)
	}
	if !got.PacketLossKnown {
		t.Fatal("PacketLossKnown = false, want true")
	}
	if got.PacketLossPercent != 0 {
		t.Errorf("PacketLossPercent = %d, want 0", got.PacketLossPercent)
	}
	if got.Message != "" {
		t.Errorf("Message = %q, want empty", got.Message)
	}
}

func TestParsePingOutputRussianNoSpace(t *testing.T) {
	got := parsePingOutput(russianPingOutputNoSpace)

	if !got.Available {
		t.Fatal("Available = false, want true")
	}
	if got.AverageMs != 38 {
		t.Errorf("AverageMs = %d, want 38", got.AverageMs)
	}
	if !got.PacketLossKnown || got.PacketLossPercent != 0 {
		t.Errorf("PacketLoss = (%v, %d), want (true, 0)", got.PacketLossKnown, got.PacketLossPercent)
	}
}

func TestParsePingOutputRussianWithSpace(t *testing.T) {
	got := parsePingOutput(russianPingOutputWithSpace)

	if !got.Available {
		t.Fatal("Available = false, want true")
	}
	if got.AverageMs != 38 {
		t.Errorf("AverageMs = %d, want 38", got.AverageMs)
	}
	if !got.PacketLossKnown || got.PacketLossPercent != 0 {
		t.Errorf("PacketLoss = (%v, %d), want (true, 0)", got.PacketLossKnown, got.PacketLossPercent)
	}
}

func TestParsePingOutputEnglishTimeout(t *testing.T) {
	got := parsePingOutput(englishTimeoutOutput)

	if got.Available {
		t.Fatal("Available = true, want false")
	}
	if got.Message != unavailableMessage {
		t.Errorf("Message = %q, want %q", got.Message, unavailableMessage)
	}
	if !got.PacketLossKnown || got.PacketLossPercent != 100 {
		t.Errorf("PacketLoss = (%v, %d), want (true, 100)", got.PacketLossKnown, got.PacketLossPercent)
	}
}

func TestParsePingOutputRussianTimeout(t *testing.T) {
	got := parsePingOutput(russianTimeoutOutput)

	if got.Available {
		t.Fatal("Available = true, want false")
	}
	if got.Message != unavailableMessage {
		t.Errorf("Message = %q, want %q", got.Message, unavailableMessage)
	}
	if !got.PacketLossKnown || got.PacketLossPercent != 100 {
		t.Errorf("PacketLoss = (%v, %d), want (true, 100)", got.PacketLossKnown, got.PacketLossPercent)
	}
}

func TestParsePingOutputEmpty(t *testing.T) {
	got := parsePingOutput("")

	if got.Available {
		t.Fatal("Available = true, want false")
	}
	if got.PacketLossKnown {
		t.Fatal("PacketLossKnown = true, want false")
	}
	if got.Message != unavailableMessage {
		t.Errorf("Message = %q, want %q", got.Message, unavailableMessage)
	}
}
