# PlatePilot iOS

This folder contains the native SwiftUI iOS app that mirrors the Vue frontend and connects to the Mobile BFF.

## Requirements

- Xcode with iOS 26 SDK (Liquid Glass APIs). The app target is iOS 26.
- `xcodegen` in your PATH
- Docker Desktop

## 1) Start the backend (Docker)

From the project root:

```bash
docker compose up
```

This starts Postgres, RabbitMQ, the services, runs migrations, and seeds sample recipes. Verify the BFF is up:

```bash
curl http://localhost:8080/health
```

You should see `OK`.

## 2) Generate the Xcode project

```bash
cd src/ios/PlatePilot
xcodegen generate
```

Open the project:

```bash
open PlatePilot.xcodeproj
```

## 3) Run the app in Simulator

- Select the `PlatePilot` scheme.
- Choose an iOS 26 simulator.
- Build and run.

The app defaults to `http://localhost:8080/v1`, matching the Mobile BFF in `docker-compose.yml` (`mobile-bff` exposes `8080:8080`). See `src/ios/PlatePilot/Shared/API/APIConfig.swift`.

## 4) Auth flow

- Use the Sign In / Create Account screens.
- Tokens are stored in Keychain; use the account sheet (avatar on Home) to sign out.
- If you already have a dev account, sign in with it. Otherwise, register a new one.

## 5) Running on a physical device

`localhost` is not reachable from a device. Update the base URL to your machine IP:

```swift
// src/ios/PlatePilot/Shared/API/APIConfig.swift
static let defaultBaseURL = URL(string: "http://192.168.1.42:8080/v1")!
```

Then rebuild and run.

Make sure your iPhone and Mac are on the same Wi-Fi network and that macOS firewall settings allow inbound connections to the Docker-exposed `8080` port.

## Troubleshooting

- If recipes do not load, confirm `docker compose up` is running and the seeder finished.
- If you see unauthorized errors, sign in again to refresh the access token.
- Regenerate the Xcode project after adding/removing Swift files:
  `xcodegen generate` in `src/ios/PlatePilot`.
