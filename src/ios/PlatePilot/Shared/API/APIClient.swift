import Foundation

struct APIClient: @unchecked Sendable {
    let baseURL: URL
    let session: URLSession
    let tokenProvider: @Sendable () async -> String?

    init(
        baseURL: URL = APIConfig.defaultBaseURL,
        session: URLSession = .shared,
        tokenProvider: @escaping @Sendable () async -> String? = { nil }
    ) {
        self.baseURL = baseURL
        self.session = session
        self.tokenProvider = tokenProvider
    }

    func fetchRecipes(pageIndex: Int = 1, pageSize: Int = 20) async throws -> PaginatedRecipesDTO {
        try await request(
            path: "recipe/all",
            queryItems: [
                URLQueryItem(name: "pageIndex", value: String(pageIndex)),
                URLQueryItem(name: "pageSize", value: String(pageSize))
            ]
        )
    }

    func fetchRecipe(id: UUID) async throws -> RecipeDTO {
        try await request(path: "recipe/\(id.uuidString)")
    }

    func register(email: String, password: String, displayName: String) async throws -> TokenResponseDTO {
        let payload = RegisterRequestDTO(email: email, password: password, displayName: displayName)
        return try await request(path: "auth/register", method: "POST", body: payload)
    }

    func login(email: String, password: String) async throws -> TokenResponseDTO {
        let payload = LoginRequestDTO(email: email, password: password)
        return try await request(path: "auth/login", method: "POST", body: payload)
    }

    func refresh(refreshToken: String) async throws -> TokenResponseDTO {
        let payload = RefreshRequestDTO(refreshToken: refreshToken)
        return try await request(path: "auth/refresh", method: "POST", body: payload)
    }

    func logout(refreshToken: String) async throws {
        let payload = RefreshRequestDTO(refreshToken: refreshToken)
        let _: EmptyResponse = try await request(path: "auth/logout", method: "POST", body: payload)
    }

    private func request<T: Decodable>(
        path: String,
        method: String = "GET",
        queryItems: [URLQueryItem] = [],
        body: Encodable? = nil
    ) async throws -> T {
        guard var components = URLComponents(url: baseURL.appendingPathComponent(path), resolvingAgainstBaseURL: false) else {
            throw APIError.invalidURL
        }
        if !queryItems.isEmpty {
            components.queryItems = queryItems
        }
        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = method
        request.setValue("application/json", forHTTPHeaderField: "Accept")

        if let token = await tokenProvider() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        if let body {
            request.setValue("application/json", forHTTPHeaderField: "Content-Type")
            request.httpBody = try JSONEncoder().encode(AnyEncodable(body))
        }

        let (data, response) = try await session.data(for: request)
        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        if httpResponse.statusCode == 401 {
            throw APIError.unauthorized
        }

        guard (200...299).contains(httpResponse.statusCode) else {
            if let errorResponse = try? JSONDecoder().decode(APIErrorResponse.self, from: data),
               let message = errorResponse.error {
                throw APIError.serverError(message)
            }
            throw APIError.serverError("Server error (\(httpResponse.statusCode)).")
        }

        if data.isEmpty, T.self == EmptyResponse.self {
            return EmptyResponse() as! T
        }

        do {
            return try JSONDecoder().decode(T.self, from: data)
        } catch {
            throw APIError.decodingFailed
        }
    }
}

private struct AnyEncodable: Encodable {
    private let encoder: (Encoder) throws -> Void

    init(_ wrapped: Encodable) {
        self.encoder = wrapped.encode
    }

    func encode(to encoder: Encoder) throws {
        try self.encoder(encoder)
    }
}
