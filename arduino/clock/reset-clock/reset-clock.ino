#include <Wire.h>
#include "RTClib.h"

RTC_DS3231 rtc;

void setup() {
  Serial.begin(9600);
  Wire.begin();

  if (!rtc.begin()) {
    Serial.println("No se encontró el módulo RTC");
    while (1);
  }

  // Detecta si el reloj perdió alimentación
  if (rtc.lostPower()) {
    Serial.println("⚠️ El RTC perdió energía, ajustando fecha/hora...");
    rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));  // Ajusta una vez
  }

  // Detectar si el reloj está detenido
  if (rtc.now().second() == 0) {
    delay(1000); // Espera 1 segundo
    if (rtc.now().second() == 0) {
      Serial.println("⏸️ Reloj detenido. Reiniciando...");
      rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));  // Reinicia el conteo
    }
  }

  Serial.println("✅ RTC activo");
}

void loop() {
  DateTime now = rtc.now();

  Serial.print("Hora actual: ");
  Serial.print(now.hour());
  Serial.print(":");
  Serial.print(now.minute());
  Serial.print(":");
  Serial.println(now.second());

  delay(1000);
}
