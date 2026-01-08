import Foundation
import Observation

@MainActor
@Observable
final class AuthStore {
    private let apiClient: APIClient
    private let tokenStore: KeychainTokenStore

    private(set) var tokens: AuthTokens?
    var isLoading = false
    var errorMessage: String?

    var isAuthenticated: Bool {
        tokens != nil
    }

    var accessToken: String? {
        tokens?.accessToken
    }

    init(apiClient: APIClient = APIClient(), tokenStore: KeychainTokenStore = KeychainTokenStore()) {
        self.apiClient = apiClient
        self.tokenStore = tokenStore
        self.tokens = try? tokenStore.load()

        if let tokens, tokens.isExpired {
            Task {
                await refreshIfNeeded()
            }
        }
    }

    func login(email: String, password: String) async -> Bool {
        isLoading = true
        errorMessage = nil
        defer { isLoading = false }

        do {
            let response = try await apiClient.login(email: email, password: password)
            let tokens = mapTokens(from: response)
            try tokenStore.save(tokens)
            self.tokens = tokens
            return true
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Unable to sign in."
            return false
        }
    }

    func register(email: String, password: String, displayName: String) async -> Bool {
        isLoading = true
        errorMessage = nil
        defer { isLoading = false }

        do {
            let response = try await apiClient.register(email: email, password: password, displayName: displayName)
            let tokens = mapTokens(from: response)
            try tokenStore.save(tokens)
            self.tokens = tokens
            return true
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Unable to create account."
            return false
        }
    }

    func refreshIfNeeded() async {
        guard let tokens else { return }
        guard tokens.isExpired else { return }

        isLoading = true
        errorMessage = nil
        defer { isLoading = false }

        do {
            let response = try await apiClient.refresh(refreshToken: tokens.refreshToken)
            let newTokens = mapTokens(from: response)
            try tokenStore.save(newTokens)
            self.tokens = newTokens
        } catch {
            errorMessage = (error as? APIError)?.errorDescription ?? "Session expired. Please sign in again."
            self.tokens = nil
            try? tokenStore.clear()
        }
    }

    func logout() async {
        let refreshToken = tokens?.refreshToken
        tokens = nil
        try? tokenStore.clear()

        guard let refreshToken else { return }
        do {
            try await apiClient.logout(refreshToken: refreshToken)
        } catch {
            // Ignore logout errors; token already cleared locally.
        }
    }

    func clearError() {
        errorMessage = nil
    }

    private func mapTokens(from response: TokenResponseDTO) -> AuthTokens {
        let buffer: TimeInterval = 30
        let expiresIn = TimeInterval(response.expiresIn)
        let expiry = Date().addingTimeInterval(max(0, expiresIn - buffer))
        return AuthTokens(accessToken: response.accessToken, refreshToken: response.refreshToken, expiresAt: expiry)
    }
}
