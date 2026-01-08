import Foundation
import Security

struct KeychainTokenStore {
    private let service = "com.platepilot.app"
    private let account = "authTokens"

    func load() throws -> AuthTokens? {
        var query = baseQuery()
        query[kSecReturnData as String] = true
        query[kSecMatchLimit as String] = kSecMatchLimitOne

        var item: CFTypeRef?
        let status = SecItemCopyMatching(query as CFDictionary, &item)
        if status == errSecItemNotFound {
            return nil
        }
        guard status == errSecSuccess else {
            throw KeychainError.readFailed(status: status)
        }
        guard let data = item as? Data else {
            return nil
        }
        return try JSONDecoder().decode(AuthTokens.self, from: data)
    }

    func save(_ tokens: AuthTokens) throws {
        let data = try JSONEncoder().encode(tokens)
        var query = baseQuery()
        let update: [String: Any] = [kSecValueData as String: data]

        let status = SecItemUpdate(query as CFDictionary, update as CFDictionary)
        if status == errSecItemNotFound {
            query[kSecValueData as String] = data
            let addStatus = SecItemAdd(query as CFDictionary, nil)
            if addStatus != errSecSuccess {
                throw KeychainError.writeFailed(status: addStatus)
            }
        } else if status != errSecSuccess {
            throw KeychainError.writeFailed(status: status)
        }
    }

    func clear() throws {
        let status = SecItemDelete(baseQuery() as CFDictionary)
        if status != errSecSuccess && status != errSecItemNotFound {
            throw KeychainError.deleteFailed(status: status)
        }
    }

    private func baseQuery() -> [String: Any] {
        [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrService as String: service,
            kSecAttrAccount as String: account
        ]
    }
}

enum KeychainError: LocalizedError {
    case readFailed(status: OSStatus)
    case writeFailed(status: OSStatus)
    case deleteFailed(status: OSStatus)

    var errorDescription: String? {
        switch self {
        case .readFailed:
            return "Unable to read credentials."
        case .writeFailed:
            return "Unable to save credentials."
        case .deleteFailed:
            return "Unable to delete credentials."
        }
    }
}
