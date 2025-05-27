
#include <Servo.h>

int TRIGG_A = 8;
int ECHO_A = 9;

int TRIGG_B = 6;
int ECHO_B = 7;

int SERVO_DOOR = 3;

int ACTIVATION_DISTANCE = 10;

Servo servoDoor;

void setup() {
  pinMode(TRIGG_A, OUTPUT);
  pinMode(ECHO_A, INPUT);

  pinMode(TRIGG_B, OUTPUT);
  pinMode(ECHO_B, INPUT);

  servoDoor.attach(SERVO_DOOR);
  servoDoor.write(99);

  pinMode(13, OUTPUT);
  pinMode(12, OUTPUT);

  Serial.begin(9600);
  //Serial.setTimeout(1);

  Serial.println("READY");
}

unsigned long lastMeasurement = micros();
int probability = 0;

void loop() {
  if (Serial.available() > 0) {
    String cmd = Serial.readStringUntil('\n');

    if (cmd == "INFO") {
      Serial.println("# Probability.");
      Serial.println(probability);
    } else if (cmd == "OPEN_DOOR") {
      Serial.println("# Opening door...");
      servoDoor.write(0);
    } else if (cmd == "CLOSE_DOOR") {
      Serial.println("# Closing door...");
      servoDoor.write(99);
    }
  }

  int ultrasoundA = measureDistance(TRIGG_A, ECHO_A);
  int ultrasoundB = measureDistance(TRIGG_B, ECHO_B);

/*
  Serial.print("# Ultrasonido A: ");
  Serial.println(ultrasoundA);
  Serial.print("# Ultrasonido B: ");
  Serial.println(ultrasoundB);
*/

  if (ultrasoundA < ACTIVATION_DISTANCE && ultrasoundA != 0) {
    unsigned long time = micros();
    if (time - lastMeasurement < 500000) {
      return;
    }
    Serial.println("# Ultrasound A Activated!");
    digitalWrite(13, HIGH);
    digitalWrite(12, LOW);
    lastMeasurement = time;
  } else if (ultrasoundB < ACTIVATION_DISTANCE && ultrasoundB != 0) {
    if (lastMeasurement == 0) {
      return;
    }
    Serial.println("# Ultrasound B Activated!");

    float time = float(micros() - lastMeasurement) / 1000000.0;
    probability = (1 - ((time + random(10) * 0.01) / (1.5))) * 100;
    lastMeasurement = 0;

    digitalWrite(12, HIGH);
    digitalWrite(13, LOW);

    Serial.print("# Probability is: ");
    Serial.println(probability);
  }
}

long measureDistance(int triggerPin, int echoPin) {
  digitalWrite(triggerPin, LOW);
  delayMicroseconds(2);
  digitalWrite(triggerPin, HIGH);
  delayMicroseconds(10);
  digitalWrite(triggerPin, LOW);

  long duration = pulseIn(echoPin, HIGH, 30000);  // Timeout after 30ms.
  return duration * 0.034 / 2;
}