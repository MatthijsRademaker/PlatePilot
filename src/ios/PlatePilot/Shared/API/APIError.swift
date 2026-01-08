import Foundation

enum APIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case decodingFailed
    case unauthorized
    case serverError(String)

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "Invalid server URL."
        case .invalidResponse:
            return "Unexpected server response."
        case .decodingFailed:
            return "Unable to read server response."
        case .unauthorized:
            return "Please sign in to continue."
        case .serverError(let message):
            return message
        }
    }
}
