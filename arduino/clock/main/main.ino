#include <Wire.h>
#include "RTClib.h"
#include <LiquidCrystal.h>

RTC_DS3231 rtc;

// Pines de conexión del LCD
const int rs = 12, en = 11, d4 = 5, d5 = 4, d6 = 3, d7 = 2;
LiquidCrystal lcd(rs, en, d4, d5, d6, d7);

void setup() {
  Serial.begin(9600);
  Wire.begin();

  // Inicializar el LCD
  lcd.begin(16, 2);
  lcd.print("Iniciando...");

  // Inicializar el RTC
  if (!rtc.begin()) {
    Serial.println("No se encontró el módulo RTC");
    lcd.clear();
    lcd.print("RTC NO ENCONTRADO");
    while (1);
  }

  // Descomenta UNA SOLA VEZ para ajustar la hora con la del compilador:
  rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));

  // O bien ajústala manualmente así (año, mes, día, hora, minuto, segundo):
  // rtc.adjust(DateTime(2025, 5, 27, 14, 30, 00));

  lcd.clear();
}

void loop() {
  DateTime now = rtc.now();

  // Formatear la hora como HH:MM:SS con ceros a la izquierda
  char buffer[9];
  snprintf(buffer, sizeof(buffer), "%02d:%02d:%02d", now.hour(), now.minute(), now.second());

  // Mostrar en la primera línea
  lcd.setCursor(0, 0);
  lcd.print("Fecha ");
  // Opcional: mostrar fecha
  char fecha[11];
  snprintf(fecha, sizeof(fecha), "%02d/%02d/%04d", now.day(), now.month(), now.year());
  lcd.print(fecha);

  // Mostrar hora en la segunda línea
  lcd.setCursor(0, 1);
  lcd.print("Hora: ");
  lcd.print(buffer);

  delay(1000);
}
