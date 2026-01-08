import SwiftUI

struct AuthFlowView: View {
    @Environment(AuthStore.self) private var authStore
    @State private var mode: AuthMode = .login

    @State private var displayName = ""
    @State private var email = ""
    @State private var password = ""
    @State private var showPassword = false

    var body: some View {
        ZStack {
            LinearGradient(
                colors: [PlatePilotTheme.accent.opacity(0.95), PlatePilotTheme.accentDeep],
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()

            AuthBackgroundView()

            VStack(spacing: 20) {
                AuthHeaderView()

                VStack(spacing: 16) {
                    Picker("", selection: $mode) {
                        ForEach(AuthMode.allCases) { mode in
                            Text(mode.title).tag(mode)
                        }
                    }
                    .pickerStyle(.segmented)

                    VStack(spacing: 14) {
                        if mode == .register {
                            AuthField(title: "Display Name") {
                                TextField("Jane Doe", text: $displayName)
                                    .textContentType(.name)
                            }
                        }

                        AuthField(title: "Email") {
                            TextField("you@example.com", text: $email)
                                .textInputAutocapitalization(.never)
                                .keyboardType(.emailAddress)
                                .textContentType(.emailAddress)
                        }

                        AuthField(title: "Password") {
                            Group {
                                if showPassword {
                                    TextField("Enter your password", text: $password)
                                        .textContentType(mode == .login ? .password : .newPassword)
                                } else {
                                    SecureField("Enter your password", text: $password)
                                        .textContentType(mode == .login ? .password : .newPassword)
                                }
                            }
                            .overlay(alignment: .trailing) {
                                Button {
                                    showPassword.toggle()
                                } label: {
                                    Image(systemName: showPassword ? "eye.slash" : "eye")
                                        .foregroundStyle(.secondary)
                                }
                                .padding(.trailing, 10)
                            }
                        }
                    }

                    if let errorMessage = authStore.errorMessage {
                        HStack(spacing: 8) {
                            Image(systemName: "exclamationmark.triangle")
                            Text(errorMessage)
                                .font(PlatePilotTheme.bodyFont(size: 13, weight: .medium))
                        }
                        .foregroundStyle(.white)
                        .padding(12)
                        .frame(maxWidth: .infinity, alignment: .leading)
                        .background(Color.white.opacity(0.15), in: RoundedRectangle(cornerRadius: 12, style: .continuous))
                    }

                    Button {
                        Task {
                            await handlePrimaryAction()
                        }
                    } label: {
                        HStack(spacing: 8) {
                            if authStore.isLoading {
                                ProgressView()
                                    .tint(.white)
                            }
                            Text(mode.primaryActionTitle)
                                .font(PlatePilotTheme.bodyFont(size: 16, weight: .semibold))
                        }
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 14)
                    }
                    .plateGlassButton(prominent: true)
                    .disabled(authStore.isLoading)

                    if mode == .login {
                        Text("Dev: seed@platepilot.local / platepilot")
                            .font(PlatePilotTheme.bodyFont(size: 12, weight: .medium))
                            .foregroundStyle(.white.opacity(0.8))
                    }
                }
                .padding(20)
                .background(Color.white.opacity(0.2), in: RoundedRectangle(cornerRadius: 24, style: .continuous))
                .padding(.horizontal, 20)
            }
            .padding(.vertical, 32)
        }
        .onChange(of: mode) { _, _ in
            authStore.clearError()
            password = ""
        }
    }

    private func handlePrimaryAction() async {
        authStore.clearError()
        switch mode {
        case .login:
            _ = await authStore.login(email: email, password: password)
        case .register:
            _ = await authStore.register(email: email, password: password, displayName: displayName)
        }
    }
}

private struct AuthHeaderView: View {
    var body: some View {
        VStack(spacing: 8) {
            ZStack {
                RoundedRectangle(cornerRadius: 18, style: .continuous)
                    .fill(Color.white.opacity(0.25))
                    .frame(width: 72, height: 72)

                Image(systemName: "fork.knife")
                    .font(.system(size: 28, weight: .semibold))
                    .foregroundStyle(.white)
            }

            Text("PlatePilot")
                .font(PlatePilotTheme.titleFont(size: 28))
                .foregroundStyle(.white)

            Text("Your personal meal companion")
                .font(PlatePilotTheme.bodyFont(size: 14, weight: .medium))
                .foregroundStyle(.white.opacity(0.85))
        }
    }
}

private struct AuthBackgroundView: View {
    var body: some View {
        GeometryReader { proxy in
            let size = proxy.size
            Circle()
                .fill(Color.white.opacity(0.12))
                .frame(width: 240, height: 240)
                .position(x: size.width * 0.85, y: size.height * 0.1)

            Circle()
                .fill(Color.white.opacity(0.08))
                .frame(width: 180, height: 180)
                .position(x: size.width * 0.15, y: size.height * 0.3)

            Circle()
                .fill(Color.white.opacity(0.1))
                .frame(width: 260, height: 260)
                .position(x: size.width * 0.2, y: size.height * 0.85)
        }
        .ignoresSafeArea()
    }
}

private struct AuthField<Content: View>: View {
    let title: String
    let content: () -> Content

    init(title: String, @ViewBuilder content: @escaping () -> Content) {
        self.title = title
        self.content = content
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            Text(title)
                .font(PlatePilotTheme.bodyFont(size: 12, weight: .semibold))
                .foregroundStyle(.white.opacity(0.85))

            content()
                .font(PlatePilotTheme.bodyFont(size: 15, weight: .medium))
                .padding(12)
                .background(Color.white.opacity(0.9), in: RoundedRectangle(cornerRadius: 12, style: .continuous))
                .overlay(
                    RoundedRectangle(cornerRadius: 12, style: .continuous)
                        .stroke(Color.white.opacity(0.4), lineWidth: 1)
                )
        }
    }
}

private enum AuthMode: String, CaseIterable, Identifiable {
    case login
    case register

    var id: String { rawValue }

    var title: String {
        switch self {
        case .login:
            return "Sign In"
        case .register:
            return "Create Account"
        }
    }

    var primaryActionTitle: String {
        switch self {
        case .login:
            return "Sign In"
        case .register:
            return "Create Account"
        }
    }
}
