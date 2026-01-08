import Foundation

struct AuthTokens: Codable, Equatable {
    let accessToken: String
    let refreshToken: String
    let expiresAt: Date

    var isExpired: Bool {
        Date() >= expiresAt
    }
}
