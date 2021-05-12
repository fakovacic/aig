#include "DHT.h"
#include <WiFi.h>
#include <WebServer.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>

#define SERVER_IP "" // set server IP
#define CTRL "johnny"
#define DHTPIN 4
#define DHTTYPE DHT22

DHT dht(DHTPIN, DHTTYPE);
WebServer server(80);

StaticJsonDocument<200> output;
IPAddress ip;

float humidity;
float temperature;
float heatIndex;
int soil;

bool water;
bool lights;
bool vent;

uint8_t soilPin = A0;
uint8_t waterPin = 12;
uint8_t ventPin = 13;
uint8_t lightsPin = 14;

void setup() {
  Serial.begin(9600);

  dht.begin();

  pinMode(waterPin, OUTPUT);
  pinMode(ventPin, OUTPUT);
  pinMode(lightsPin, OUTPUT);

  WiFi.begin("", ""); // set wifi creds

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  ip = WiFi.localIP();
  Serial.println("Connected! IP address: ");
  Serial.print(ip);
  registerController();

  server.on("/water/on", waterOn);
  server.on("/water/off", waterOff);
  server.on("/lights/on", lightsOn);
  server.on("/lights/off", lightsOn);
  server.on("/vent/on", ventOn);
  server.on("/vent/off", ventOn);
  server.on("/stats", stats);
  server.onNotFound(notFound);
  server.begin();
}

void loop() {
  server.handleClient();
  setPower();
}

void setPower() {
  if (water == true) {
    digitalWrite(waterPin, HIGH);
  } else {
    digitalWrite(waterPin, LOW);
  }

  if (lights == true) {
    digitalWrite(lightsPin, HIGH);
  } else {
    digitalWrite(lightsPin, LOW);
  }

  if (vent == true) {
    digitalWrite(ventPin, HIGH);
  } else {
    digitalWrite(ventPin, LOW);
  }
  return;
}

void waterOn() {
  water = true;
  Serial.println("Water: ON");
  stats();
}

void waterOff() {
  water = false;
  Serial.println("Water: OFF");
  stats();
}

void lightsOn() {
  lights = true;
  Serial.println("Lights: ON");
  stats();
}

void lightsOff() {
  lights = false;
  Serial.println("Lights: OFF");
  stats();
}

void ventOn() {
  vent = true;
  Serial.println("Vent: ON");
  stats();
}

void ventOff() {
  vent = false;
  Serial.println("Vent: OFF");
  stats();
}

void notFound() {
  server.send(404, "application/json", "");
}

void stats() {
  readStats();
  printStats();
}

void readStats() {
  readDHT22();
  readSoil();
  return;
}

void readSoil() {
  soil = analogRead(soilPin);
  if (soil == 0) {
    Serial.println(F("Failed to read from Soil sensor!"));
    return;
  }
  soil = map(soil, 550, 10, 0, 100);
  return;
}

void readDHT22() {
  float h = dht.readHumidity();
  float t = dht.readTemperature();

  if (isnan(h) || isnan(t)) {
    Serial.println(F("Failed to read from DHT sensor!"));
    return;
  }

  humidity = h;
  temperature = t;
  heatIndex = dht.computeHeatIndex(t, h, false);
  return;
}

void printStats() {
  char url[50];
  char body[200];
  output["id"] = CTRL;
  output["ip"] = ip.toString();
  output["humidity"] = humidity;
  output["temperature"] = temperature;
  output["heatIndex"] = heatIndex;
  output["soil"] = soil;
  output["water"] = water;
  output["lights"] = lights;
  output["vent"] = vent;

  serializeJson(output, body);

  server.send(200, "application/json", body);
  return;
}

void registerController() {
  char url[50];
  char body[200];
  output["id"] = CTRL;
  output["ip"] = ip.toString();
  serializeJson(output, body);

  snprintf(url, 50, "http://%s/register", SERVER_IP);
  sendHttpRequest(url, body);
  return;
}

void sendHttpRequest(char* url, char* body) {
  if ((WiFi.status() == WL_CONNECTED)) {
    Serial.println(url);
    Serial.println(body);
    WiFiClient client;
    HTTPClient http;
    http.begin(client, url);
    http.addHeader("Content-Type", "application/json");

    int httpCode = http.POST(body);
    if (httpCode < 0) {
      Serial.println(sprintf("register failed, error: %s\n", http.errorToString(httpCode).c_str()));
    }
    http.end();
  }
  return;
}
